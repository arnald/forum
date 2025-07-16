package main

import (
	"log"
	"net/http"

	"github.com/arnald/forum/cmd/client/config"
	"github.com/arnald/forum/cmd/client/handler"
)

func main() {
	cfg, err := config.LoadClientConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	router := setupRoutes()
	log.Printf("Client started port: %s (%s environment)", cfg.Port, cfg.Environment)
	err = http.ListenAndServe(":"+cfg.Port, router)
	if err != nil {
		log.Fatal(err)
	}
}

func setupRoutes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", handler.HomePage)
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../frontend/"))))

	return router
}
