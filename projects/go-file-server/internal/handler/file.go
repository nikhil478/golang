package handler

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nikhil478/golang/go-file-server/internal/metrics"
)

type Config struct {
	MaxFileSize    int64
	ChunkSize      int
	MaxConcurrency int
}

type FileHandler struct {
	cfg     Config
	metrics *metrics.Metrics
	slots   chan struct{}
	buffers sync.Pool
	ready   atomic.Bool
}

func New(cfg Config, m *metrics.Metrics) *FileHandler {
	if cfg.MaxFileSize <= 0 {
		cfg.MaxFileSize = 10 << 30
	}
	if cfg.ChunkSize <= 0 {
		cfg.ChunkSize = 1 << 20
	}
	if cfg.MaxConcurrency <= 0 {
		cfg.MaxConcurrency = 32
	}
	if m == nil {
		m = &metrics.Metrics{}
	}

	h := &FileHandler{
		cfg:     cfg,
		metrics: m,
		slots:   make(chan struct{}, cfg.MaxConcurrency),
	}
	h.buffers.New = func() any {
		return make([]byte, cfg.ChunkSize)
	}
	h.ready.Store(true)
	return h
}

// patternReader generates reproducible, non-zero bytes without holding the
// complete payload in memory.
type patternReader struct {
	offset int64
}

func (r *patternReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = byte((r.offset + int64(i)) % 251)
	}
	r.offset += int64(len(p))
	return len(p), nil
}

// countingWriter counts actual bytes written to the HTTP response, including
// ZIP headers and the ZIP central directory.
type countingWriter struct {
	w http.ResponseWriter
	n int64
}

func (cw *countingWriter) Write(p []byte) (int, error) {
	n, err := cw.w.Write(p)
	cw.n += int64(n)
	return n, err
}

func (h *FileHandler) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "OK\n")
}

func (h *FileHandler) Ready(w http.ResponseWriter, _ *http.Request) {
	if !h.ready.Load() {
		http.Error(w, "not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "READY\n")
}

func (h *FileHandler) Download(w http.ResponseWriter, r *http.Request) {
	h.metrics.Requests.Add(1)

	size, err := parseSize(r.PathValue("size"))
	if err != nil || size > h.cfg.MaxFileSize {
		http.Error(w, "invalid or excessive file size", http.StatusBadRequest)
		return
	}

	// HEAD reports the archive media type and filename, but intentionally
	// omits Content-Length because the ZIP is generated as a live stream.
	if r.Method == http.MethodHead {
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="test-%s.zip"`, r.PathValue("size")))
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusOK)
		return
	}

	select {
	case h.slots <- struct{}{}:
		defer func() { <-h.slots }()
	case <-r.Context().Done():
		return
	default:
		h.metrics.Rejected.Add(1)
		w.Header().Set("Retry-After", "2")
		http.Error(w, "server capacity reached", http.StatusServiceUnavailable)
		return
	}

	h.metrics.Active.Add(1)
	defer h.metrics.Active.Add(-1)

	filename := "test-" + r.PathValue("size") + ".zip"
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	started := time.Now()
	cw := &countingWriter{w: w}
	zw := zip.NewWriter(cw)

	// Store avoids CPU-intensive compression, making this useful for
	// measuring downloader and network throughput.
	entryHeader := &zip.FileHeader{
		Name:   "payload.bin",
		Method: zip.Store,
	}
	entryHeader.SetMode(0644)

	entry, err := zw.CreateHeader(entryHeader)
	if err != nil {
		h.metrics.Failed.Add(1)
		return
	}

	buf := h.buffers.Get().([]byte)
	defer h.buffers.Put(buf)

	payloadWritten, copyErr := io.CopyBuffer(
		entry,
		io.LimitReader(&patternReader{}, size),
		buf,
	)

	// Close writes the ZIP central directory and is required for a valid ZIP.
	closeErr := zw.Close()

	h.metrics.PayloadBytes.Add(payloadWritten)
	h.metrics.ResponseBytes.Add(cw.n)

	if copyErr != nil || closeErr != nil || payloadWritten != size {
		h.metrics.Failed.Add(1)
		fmt.Printf("download failed file=%s payload=%d/%d copy_err=%v close_err=%v\n",
			filename, payloadWritten, size, copyErr, closeErr)
		return
	}

	h.metrics.Completed.Add(1)
	fmt.Printf("download complete file=%s payload=%d response=%d duration=%s\n",
		filename, payloadWritten, cw.n, time.Since(started))
}

func (h *FileHandler) Shutdown() {
	h.ready.Store(false)
}

// Keep strconv referenced here as a compile-time reminder that sizes are
// represented as exact byte counts in HTTP metadata where available.
var _ = strconv.IntSize
