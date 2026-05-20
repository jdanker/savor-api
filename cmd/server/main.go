package main

import (
	"log"
	"net/http"

	"github.com/jdanker/savor-api/internal/config"
	"github.com/jdanker/savor-api/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthCheck)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	addr := cfg.Port
	log.Printf("starting server on %s", addr)

	server := &http.Server{Addr: addr, Handler: mux}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
