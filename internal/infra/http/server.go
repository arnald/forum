package http

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/session"
	createcategory "github.com/arnald/forum/internal/infra/http/category/createCategory"
	deletecategory "github.com/arnald/forum/internal/infra/http/category/deleteCategory"
	getcategorybyid "github.com/arnald/forum/internal/infra/http/category/getCategoryByID"
	updatecategory "github.com/arnald/forum/internal/infra/http/category/updateCategory"
	"github.com/arnald/forum/internal/infra/http/health"
	createtopic "github.com/arnald/forum/internal/infra/http/topic/createTopic"
	deletetopic "github.com/arnald/forum/internal/infra/http/topic/deleteTopic"
	getalltopics "github.com/arnald/forum/internal/infra/http/topic/getAllTopics"
	gettopic "github.com/arnald/forum/internal/infra/http/topic/getTopic"
	updatetopic "github.com/arnald/forum/internal/infra/http/topic/updateTopic"
	userLogin "github.com/arnald/forum/internal/infra/http/user/login"
	userRegister "github.com/arnald/forum/internal/infra/http/user/register"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/infra/storage/sessionstore"
)

const (
	apiContext   = "/api/v1"
	readTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 15 * time.Second
)

type Server struct {
	appServices    app.Services
	config         *config.ServerConfig
	router         *http.ServeMux
	sessionManager session.Manager
	middleware     *middleware.Middleware
	db             *sql.DB
	logger         logger.Logger
}

func NewServer(cfg *config.ServerConfig, db *sql.DB, logger logger.Logger, appServices app.Services) *Server {
	httpServer := &Server{
		router:      http.NewServeMux(),
		appServices: appServices,
		config:      cfg,
		db:          db,
		logger:      logger,
	}
	httpServer.initSessionManager()
	httpServer.initMiddleware(httpServer.sessionManager)
	httpServer.AddHTTPRoutes()
	return httpServer
}

// FOR MIDDLEWARE CHAINING.
func middlewareChain(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}

func (server *Server) AddHTTPRoutes() {
	server.router.HandleFunc(apiContext+"/health",
		middlewareChain(
			health.NewHandler(server.logger).HealthCheck,
			server.middleware.Authorization.Required,
		))

	// User routes
	server.router.HandleFunc(apiContext+"/login/email",
		userLogin.NewHandler(server.config, server.appServices, server.sessionManager, server.logger).UserLoginEmail,
	)
	server.router.HandleFunc(apiContext+"/login/username",
		userLogin.NewHandler(server.config, server.appServices, server.sessionManager, server.logger).UserLoginUsername,
	)
	server.router.HandleFunc(apiContext+"/register",
		userRegister.NewHandler(server.config, server.appServices, server.sessionManager, server.logger).UserRegister,
	)

	// Topic routes
	server.router.HandleFunc(apiContext+"/topics/create",
		middlewareChain(
			createtopic.NewHandler(server.appServices, server.config, server.logger).CreateTopic,
			server.middleware.Authorization.Required,
		),
	)
	server.router.HandleFunc(apiContext+"/topics/update",
		middlewareChain(
			updatetopic.NewHandler(server.appServices, server.config, server.logger).UpdateTopic,
			server.middleware.Authorization.Required,
		),
	)
	server.router.HandleFunc(apiContext+"/topics/delete",
		middlewareChain(
			deletetopic.NewHandler(server.appServices, server.config, server.logger).DeleteTopic,
			server.middleware.Authorization.Required,
		),
	)
	server.router.HandleFunc(apiContext+"/topics/get",
		gettopic.NewHandler(server.appServices, server.config, server.logger).GetTopic,
	)
	server.router.HandleFunc(apiContext+"/topics/all",
		getalltopics.NewHandler(server.appServices, server.config, server.logger).GetAllTopics,
	)

	// Category routes
	server.router.HandleFunc(apiContext+"/categories/create",
		middlewareChain(
			createcategory.NewHandler(server.appServices, server.config, server.logger).CreateCategory,
			server.middleware.Authorization.Required,
		),
	)
	server.router.HandleFunc(apiContext+"/categories/delete",
		middlewareChain(
			deletecategory.NewHandler(server.appServices, server.config, server.logger).DeleteCategory,
			server.middleware.Authorization.Required,
		),
	)
	server.router.HandleFunc(apiContext+"/categories/update",
		middlewareChain(
			updatecategory.NewHandler(server.appServices, server.config, server.logger).UpdateCategory,
			server.middleware.Authorization.Required,
		),
	)
	server.router.HandleFunc(apiContext+"/categories/get",
		middlewareChain(
			getcategorybyid.NewHandler(server.appServices, server.config, server.logger).GetCategoryByID,
			server.middleware.Authorization.Optional,
		),
	)
}

func (server *Server) ListenAndServe() {
	corsWrappedRouter := middleware.NewCorsMiddleware(server.router)

	srv := &http.Server{
		Addr:         server.config.Host + ":" + server.config.Port,
		Handler:      corsWrappedRouter,
		ReadTimeout:  server.config.ReadTimeout,
		WriteTimeout: server.config.WriteTimeout,
		IdleTimeout:  server.config.IdleTimeout,
	}
	server.logger.PrintInfo("Starting server", map[string]string{
		"host":        server.config.Host,
		"port":        server.config.Port,
		"environment": server.config.Environment,
	})
	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		server.logger.PrintFatal(err, nil)
	}
}

func (server *Server) initSessionManager() {
	server.sessionManager = sessionstore.NewSessionManager(server.db, server.config.SessionManager)
}

func (server *Server) initMiddleware(sessionManager session.Manager) {
	server.middleware = middleware.NewMiddleware(sessionManager)
}
