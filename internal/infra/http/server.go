package http

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/arnald/forum/internal/infra/http/health"
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
	router *http.ServeMux
}

func NewServer() *Server {
	httpServer := &Server{
		router: http.NewServeMux(),
	}
	httpServer.AddHTTPRoutes()
	return httpServer
}

func (server *Server) AddHTTPRoutes() {
	// server.router.HandleFunc(apiContext+"/users", user.NewHandler(server.appServices.UserServices).GetAllUsers)
	server.router.HandleFunc(apiContext+"/health", health.NewHandler().HealthCheck)
}

func (server *Server) ListenAndServe(port string) {
	log.Printf("Server started port %s", port)

	srv := &http.Server{
		Addr:         port,
		Handler:      server.router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed: %v", err)
	}
}
