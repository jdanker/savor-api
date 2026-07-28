package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jdanker/savor-api/internal/config"
	"github.com/jdanker/savor-api/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthCheck)

	addr := cfg.Port
	log.Printf("starting server on %s", addr)

	server := &http.Server{
		Addr: addr, 
		Handler: mux, 
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 60 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}


}
