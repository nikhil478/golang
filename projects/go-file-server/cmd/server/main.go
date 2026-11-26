package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/nikhil478/golang/go-file-server/internal/handler"
	"github.com/nikhil478/golang/go-file-server/internal/metrics"
)

const maxFileSize int64 = 10 << 30 // 10 GiB payload

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		log.Fatalf("%s must be a positive integer; got %q", key, value)
	}
	return n
}

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg := handler.Config{
		MaxFileSize:    maxFileSize,
		ChunkSize:      envInt("CHUNK_SIZE", 1024*1024),
		MaxConcurrency: envInt("MAX_CONCURRENCY", 32),
	}

	m := &metrics.Metrics{}
	fileHandler := handler.New(cfg, m)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", fileHandler.Health)
	mux.HandleFunc("GET /readyz", fileHandler.Ready)
	mux.Handle("GET /metrics", m)
	mux.HandleFunc("GET /files/{size}", fileHandler.Download)
	mux.HandleFunc("HEAD /files/{size}", fileHandler.Download)

	server := &http.Server{
		Addr:              host + ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
		// Do not set WriteTimeout: large downloads may legitimately take a
		// long time. Client cancellation and graceful shutdown are supported.
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("ZIP file server listening on http://%s:%s", host, port)
		log.Printf("max payload: %d bytes; chunk: %d bytes; concurrency: %d",
			cfg.MaxFileSize, cfg.ChunkSize, cfg.MaxConcurrency)

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
	case err := <-serverErr:
		if err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
		return
	}

	fileHandler.Shutdown()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		_ = server.Close()
	}
	log.Println("server stopped")
}
