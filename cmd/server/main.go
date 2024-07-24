package main

import (
	"cm_open_api/internal/config"
	"cm_open_api/internal/handlers"
	"cm_open_api/internal/metrics"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	go metrics.SetupPrometheus()

	r := handlers.SetupRouter(cfg)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
