package http

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/session"
	createcategory "github.com/arnald/forum/internal/infra/http/category/createCategory"
	deletecategory "github.com/arnald/forum/internal/infra/http/category/deleteCategory"
	getallcategories "github.com/arnald/forum/internal/infra/http/category/getAllCategories"
	getcategorybyid "github.com/arnald/forum/internal/infra/http/category/getCategoryByID"
	updatecategory "github.com/arnald/forum/internal/infra/http/category/updateCategory"
	createcomment "github.com/arnald/forum/internal/infra/http/comment/createComment"
	deletecomment "github.com/arnald/forum/internal/infra/http/comment/deleteComment"
	getcomment "github.com/arnald/forum/internal/infra/http/comment/getComment"
	getcommentsbytopic "github.com/arnald/forum/internal/infra/http/comment/getCommentsByTopic"
	updatecomment "github.com/arnald/forum/internal/infra/http/comment/updateComment"
	"github.com/arnald/forum/internal/infra/http/health"
	getnotifications "github.com/arnald/forum/internal/infra/http/notification/getNotifications"
	getunreadcount "github.com/arnald/forum/internal/infra/http/notification/getUnreadCount"
	markallasread "github.com/arnald/forum/internal/infra/http/notification/markAllAsRead"
	markasread "github.com/arnald/forum/internal/infra/http/notification/markAsRead"
	streamnotification "github.com/arnald/forum/internal/infra/http/notification/streamNotification"
	oauthlogin "github.com/arnald/forum/internal/infra/http/oauth"
	createtopic "github.com/arnald/forum/internal/infra/http/topic/createTopic"
	deletetopic "github.com/arnald/forum/internal/infra/http/topic/deleteTopic"
	getalltopics "github.com/arnald/forum/internal/infra/http/topic/getAllTopics"
	gettopic "github.com/arnald/forum/internal/infra/http/topic/getTopic"
	updatetopic "github.com/arnald/forum/internal/infra/http/topic/updateTopic"
	getme "github.com/arnald/forum/internal/infra/http/user/getMe"
	userLogin "github.com/arnald/forum/internal/infra/http/user/login"
	"github.com/arnald/forum/internal/infra/http/user/logout"
	userRegister "github.com/arnald/forum/internal/infra/http/user/register"
	castvote "github.com/arnald/forum/internal/infra/http/vote/castVote"
	deletevote "github.com/arnald/forum/internal/infra/http/vote/deleteVote"
	getCounts "github.com/arnald/forum/internal/infra/http/vote/getVoteCounts"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/infra/storage/notifications"
	"github.com/arnald/forum/internal/infra/storage/sessionstore"
	oauth "github.com/arnald/forum/internal/pkg/oAuth"
	"github.com/arnald/forum/internal/pkg/oAuth/githubclient"
	"github.com/arnald/forum/internal/pkg/oAuth/googleclient"
)

const (
	apiContext               = "/api/v1"
	readTimeout              = 5 * time.Second
	writeTimeout             = 10 * time.Second
	idleTimeout              = 15 * time.Second
	stateManagerDefaultLimit = 10
)

type Server struct {
	appServices    app.Services
	config         *config.ServerConfig
	router         *http.ServeMux
	sessionManager session.Manager
	oauth          *OAuth
	notifications  *notifications.NotificationService
	middleware     *middleware.Middleware
	db             *sql.DB
	logger         logger.Logger
}

type OAuth struct {
	stateManager   *oauth.StateManager
	githubProvider *githubclient.GitHubProvider
	googleProvider *googleclient.GoogleProvider
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
	httpServer.initNotifications()
	httpServer.initOAuthServices()
	httpServer.initMiddleware(httpServer.sessionManager)
	httpServer.AddHTTPRoutes()
	return httpServer
}

func middlewareChain(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}

func (server *Server) AddHTTPRoutes() {
	server.router.HandleFunc(apiContext+"/health",
		middlewareChain(
			health.NewHandler(server.logger, server.notifications).HealthCheck,
			server.middleware.Authorization.Optional,
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
	server.router.HandleFunc(apiContext+"/logout",
		middlewareChain(
			logout.NewHandler(server.sessionManager, server.logger).Logout,
			server.middleware.Authorization.Required,
		))
	// New handler for retrieving current user data from backend
	server.router.HandleFunc(apiContext+"/me",
		middlewareChain(
			getme.NewHandler(server.logger).GetMe,
			server.middleware.Authorization.Required,
		))
	// OAuth routes
	server.router.HandleFunc(apiContext+"/auth/github/login",
		oauthlogin.NewOAuthHandler(
			server.oauth.githubProvider,
			server.config,
			&server.appServices.UserServices.Queries.UserLoginGithub,
			server.oauth.stateManager,
			server.sessionManager,
			server.logger,
		).Login,
	)
	server.router.HandleFunc(apiContext+"/auth/github/callback",
		oauthlogin.NewOAuthHandler(
			server.oauth.githubProvider,
			server.config,
			&server.appServices.UserServices.Queries.UserLoginGithub,
			server.oauth.stateManager,
			server.sessionManager,
			server.logger,
		).Callback,
	)
	server.router.HandleFunc(apiContext+"/auth/google/login",
		oauthlogin.NewOAuthHandler(
			server.oauth.googleProvider,
			server.config,
			&server.appServices.UserServices.Queries.UserLoginGithub,
			server.oauth.stateManager,
			server.sessionManager,
			server.logger,
		).Login,
	)
	server.router.HandleFunc(apiContext+"/auth/google/callback",
		oauthlogin.NewOAuthHandler(
			server.oauth.googleProvider,
			server.config,
			&server.appServices.UserServices.Queries.UserLoginGithub,
			server.oauth.stateManager,
			server.sessionManager,
			server.logger,
		).Callback,
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
	server.router.HandleFunc(apiContext+"/topic",
		middlewareChain(
			gettopic.NewHandler(server.appServices, server.config, server.logger).GetTopic,
			server.middleware.Authorization.Optional,
		),
	)
	server.router.HandleFunc(apiContext+"/topics/all",
		middlewareChain(
			getalltopics.NewHandler(server.appServices, server.config, server.logger).GetAllTopics,
			server.middleware.Authorization.Optional,
		),
	)

	// Comment routes
	server.router.HandleFunc(apiContext+"/comments/create",
		middlewareChain(
			createcomment.NewHandler(server.appServices, server.config, server.logger, server.notifications).CreateComment,
			server.middleware.Authorization.Required,
		),
	)
	server.router.HandleFunc(apiContext+"/comments/update",
		middlewareChain(
			updatecomment.NewHandler(server.appServices, server.config, server.logger).UpdateComment,
			server.middleware.Authorization.Required,
		),
	)
	server.router.HandleFunc(apiContext+"/comments/delete",
		middlewareChain(
			deletecomment.NewHandler(server.appServices, server.config, server.logger).DeleteComment,
			server.middleware.Authorization.Required,
		),
	)
	server.router.HandleFunc(apiContext+"/comments/get",
		getcomment.NewHandler(server.appServices, server.config, server.logger).GetComment,
	)
	server.router.HandleFunc(apiContext+"/comments/topic",
		getcommentsbytopic.NewHandler(server.appServices, server.config, server.logger).GetCommentsByTopic,
	)

	// Category routes
	server.router.HandleFunc(apiContext+"/category/create",
		middlewareChain(
			createcategory.NewHandler(server.appServices, server.config, server.logger).CreateCategory,
			server.middleware.Authorization.Required,
		),
	)
	server.router.HandleFunc(apiContext+"/category/delete",
		middlewareChain(
			deletecategory.NewHandler(server.appServices, server.config, server.logger).DeleteCategory,
			server.middleware.Authorization.Required,
		),
	)
	server.router.HandleFunc(apiContext+"/category/update",
		middlewareChain(
			updatecategory.NewHandler(server.appServices, server.config, server.logger).UpdateCategory,
			server.middleware.Authorization.Required,
		),
	)
	server.router.HandleFunc(apiContext+"/category",
		middlewareChain(
			getcategorybyid.NewHandler(server.appServices, server.config, server.logger).GetCategoryByID,
			server.middleware.Authorization.Optional,
		),
	)
	server.router.HandleFunc(apiContext+"/categories/all",
		getallcategories.NewHandler(server.appServices, server.config, server.logger).GetAllCategories,
	)

	// Vote routes
	server.router.HandleFunc(apiContext+"/vote/cast",
		middlewareChain(
			castvote.NewHandler(server.appServices, server.config, server.logger, server.notifications).CastVote,
			server.middleware.Authorization.Required,
		),
	)

	server.router.HandleFunc(apiContext+"/vote/delete",
		middlewareChain(
			deletevote.NewHandler(server.appServices, server.config, server.logger).DeleteVote,
			server.middleware.Authorization.Required,
		),
	)

	server.router.HandleFunc(apiContext+"/vote/counts",
		middlewareChain(
			getCounts.NewHandler(server.appServices, server.config, server.logger).GetCounts,
			server.middleware.Authorization.Optional,
		),
	)

	// Notifications routes

	server.router.HandleFunc(apiContext+"/notifications/stream", // get
		middlewareChain(
			streamnotification.NewHandler(server.notifications).StreamNotifications,
			server.middleware.Authorization.Required,
		),
	)

	server.router.HandleFunc(apiContext+"/otifications/unread-count", // get
		middlewareChain(
			getunreadcount.NewHandler(server.notifications).GetUnread,
			server.middleware.Authorization.Required,
		),
	)

	server.router.HandleFunc(apiContext+"/notifications", // get
		middlewareChain(
			getnotifications.NewHandler(server.notifications).GetNotifications,
			server.middleware.Authorization.Required,
		),
	)

	server.router.HandleFunc(apiContext+"/notifications/mark-read", // post
		middlewareChain(
			markasread.NewHandler(server.notifications).MarkAsRead,
			server.middleware.Authorization.Required,
		),
	)

	server.router.HandleFunc(apiContext+"/notifications/mark-all-read", // post
		middlewareChain(
			markallasread.NewHandler(server.notifications).MarkAllAsRead,
			server.middleware.Authorization.Required,
		),
	)
}

func (server *Server) ListenAndServe() {
	wrappedRouter := middleware.NewCorsMiddleware(server.router)

	if server.config.RateLimit.Enabled {
		wrappedRouter = middleware.NewRateLimiterMiddleware(
			wrappedRouter,
			server.config.RateLimit.RequestsLimit,
			server.config.RateLimit.WindowSeconds,
			server.config.RateLimit.Cleanup,
		)
		server.logger.PrintInfo("Rate Limit wrapped", nil)
		log.Printf("  2. Rate Limit middleware (limit: %d req/%ds cleanup: %s)",
			server.config.RateLimit.RequestsLimit,
			server.config.RateLimit.WindowSeconds,
			server.config.RateLimit.Cleanup.String())
	}

	srv := &http.Server{
		Addr:         server.config.Host + ":" + server.config.Port,
		Handler:      wrappedRouter,
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

func (server *Server) initNotifications() {
	server.notifications = notifications.NewNotificationService(server.db)
}

func (server *Server) initMiddleware(sessionManager session.Manager) {
	server.middleware = middleware.NewMiddleware(sessionManager)
}

func (server *Server) initOAuthServices() {
	server.oauth = &OAuth{
		stateManager: oauth.NewStateManager(stateManagerDefaultLimit * time.Minute),
		githubProvider: githubclient.NewProvider(
			server.config.OAuth.GitHub.ClientID,
			server.config.OAuth.GitHub.ClientSecret,
			server.config.OAuth.GitHub.RedirectURL,
			server.config.OAuth.GitHub.Scopes,
		),
		googleProvider: googleclient.NewProvider(
			server.config.OAuth.Google.ClientID,
			server.config.OAuth.Google.ClientSecret,
			server.config.OAuth.Google.RedirectURL,
			server.config.OAuth.Google.TokenURL,
			server.config.OAuth.Google.Scopes,
		),
	}
}
