package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/arnald/forum/cmd/client/config"
	"github.com/arnald/forum/cmd/client/handler"
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

	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	staticPath := filepath.Join(basePath, "frontend", "static")
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	// Homepage
	router.HandleFunc("/", handler.HomePage)
	// Register page
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.RegisterPage(w, r)
		case http.MethodPost:
			handler.RegisterPost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})
	// Login page
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.LoginPage(w, r)
		case http.MethodPost:
			handler.LoginPost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Create post page (GET/POST)
	router.HandleFunc("/post/create", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.CreatePostPage(w, r)
		case http.MethodPost:
			handler.CreatePostPost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/category/", handler.CategoryPage)

	return router
}
