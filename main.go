package main

import (
	"log"

	"github.com/ben/ubiquiti-monitor/internal/api"
	"github.com/ben/ubiquiti-monitor/internal/config"
	"github.com/ben/ubiquiti-monitor/internal/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.Init(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize and start the API server
	server := api.NewServer(cfg.Server, db)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
