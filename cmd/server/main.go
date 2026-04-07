package main

import (
	"log"
	"net/http"
	"os"

    "github.com/jdanker/savor-api/internal/handlers"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/health", handlers.HealthCheck)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    addr := ":" + port
    log.Printf("starting server on %s", addr)

    server := &http.Server{Addr: addr, Handler: mux}
    if err := server.ListenAndServe(); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}

