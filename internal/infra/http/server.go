package http

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra/http/health"
	"github.com/arnald/forum/internal/infra/storage/sqlite"
)

const (
	apiContext   = "/api/v1"
	readTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 15 * time.Second
)

type Server struct {
	// ctx context.Context
	// appServices app.Services
	config *config.ServerConfig
	router *http.ServeMux
}

func NewServer() *Server {
	httpServer := &Server{
		router: http.NewServeMux(),
	}
	httpServer.loadConfiguration()
	httpServer.loadDatabase()
	httpServer.AddHTTPRoutes()
	return httpServer
}

func (server *Server) AddHTTPRoutes() {
	// server.router.HandleFunc(apiContext+"/users", user.NewHandler(server.appServices.UserServices).GetAllUsers)
	server.router.HandleFunc(apiContext+"/health", health.NewHandler().HealthCheck)
}

func (server *Server) ListenAndServe() {
	srv := &http.Server{
		Addr:         server.config.Host + ":" + server.config.Port,
		Handler:      server.router,
		ReadTimeout:  server.config.ReadTimeout,
		WriteTimeout: server.config.WriteTimeout,
		IdleTimeout:  server.config.IdleTimeout,
	}

	log.Printf("Server started port: %s (%s environment)", server.config.Port, server.config.Environment)
	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed: %v", err)
	}
}

func (server *Server) loadConfiguration() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	server.config = cfg
}

func (server *Server) loadDatabase() {
	_, err := sqlite.InitializeDB(*server.config)
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
}
