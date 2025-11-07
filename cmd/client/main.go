package main

import (
	"log"

	"github.com/arnald/forum/cmd/client/config"
	"github.com/arnald/forum/cmd/client/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadClientConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Create client server
	clientServer, err := server.NewClientServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create client server: %v", err)
	}

	// Setup routes
	clientServer.SetupRoutes()

	// Start server
	if err := clientServer.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
