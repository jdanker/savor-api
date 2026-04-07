package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type healthResponse struct {
    Status string `json:"status"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
    log.Println("GET /health")
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(healthResponse{Status: "ok"})
}
