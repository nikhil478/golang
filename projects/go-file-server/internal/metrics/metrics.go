package metrics

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// Metrics contains process-local counters for the test server.
type Metrics struct {
	Requests       atomic.Int64
	Active         atomic.Int64
	Completed      atomic.Int64
	Failed         atomic.Int64
	Rejected       atomic.Int64
	PayloadBytes   atomic.Int64
	ResponseBytes  atomic.Int64
}

func (m *Metrics) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	fmt.Fprintln(w, "# TYPE zip_file_requests_total counter")
	fmt.Fprintf(w, "zip_file_requests_total %d\n", m.Requests.Load())

	fmt.Fprintln(w, "# TYPE zip_file_active_downloads gauge")
	fmt.Fprintf(w, "zip_file_active_downloads %d\n", m.Active.Load())

	fmt.Fprintln(w, "# TYPE zip_file_downloads_completed_total counter")
	fmt.Fprintf(w, "zip_file_downloads_completed_total %d\n", m.Completed.Load())

	fmt.Fprintln(w, "# TYPE zip_file_downloads_failed_total counter")
	fmt.Fprintf(w, "zip_file_downloads_failed_total %d\n", m.Failed.Load())

	fmt.Fprintln(w, "# TYPE zip_file_downloads_rejected_total counter")
	fmt.Fprintf(w, "zip_file_downloads_rejected_total %d\n", m.Rejected.Load())

	fmt.Fprintln(w, "# TYPE zip_file_payload_bytes_total counter")
	fmt.Fprintf(w, "zip_file_payload_bytes_total %d\n", m.PayloadBytes.Load())

	fmt.Fprintln(w, "# TYPE zip_file_http_response_bytes_total counter")
	fmt.Fprintf(w, "zip_file_http_response_bytes_total %d\n", m.ResponseBytes.Load())
}