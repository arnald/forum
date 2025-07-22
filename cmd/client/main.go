package main

import (
	"log"
	"net/http"

	"github.com/arnald/forum/cmd/client/config"
	"github.com/arnald/forum/cmd/client/handler"
	"github.com/arnald/forum/internal/pkg/path"
)

func main() {
	cfg, err := config.LoadClientConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	router := setupRoutes()
	client := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: cfg.HTTPTimeouts.ReadHeader,
		ReadTimeout:       cfg.HTTPTimeouts.Read,
		WriteTimeout:      cfg.HTTPTimeouts.Write,
		IdleTimeout:       cfg.HTTPTimeouts.Idle,
	}

	log.Printf("Client started port: %s (%s environment)", cfg.Port, cfg.Environment)
	err = client.ListenAndServe()
	if err != nil {
		log.Fatal("Client error: ", err)
	}
}

func setupRoutes() *http.ServeMux {
	router := http.NewServeMux()

	resolver := path.NewResolver()
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(resolver.GetPath("frontend/static/")))))
	router.HandleFunc("/", handler.HomePage)

	return router
}
