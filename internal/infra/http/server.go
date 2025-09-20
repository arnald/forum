package http

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/infra/http/category"
	"github.com/arnald/forum/internal/infra/http/comment"
	"github.com/arnald/forum/internal/infra/http/health"
	"github.com/arnald/forum/internal/infra/http/post"
	userLogin "github.com/arnald/forum/internal/infra/http/user/login"
	userRegister "github.com/arnald/forum/internal/infra/http/user/register"
	"github.com/arnald/forum/internal/infra/http/vote"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/infra/session"
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
	sessionManager user.SessionManager
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
	httpServer.AddHTTPRoutes()
	return httpServer
}

func (server *Server) AddHTTPRoutes() {
	// Health endpoint
	server.router.HandleFunc(apiContext+"/health", health.NewHandler(server.logger).HealthCheck)

	// User authentication endpoints
	server.router.HandleFunc(
		apiContext+"/login/username",
		userLogin.NewHandler(server.config, server.appServices, server.sessionManager, server.logger).UserLoginUsername,
	)
	server.router.HandleFunc(
		apiContext+"/login/email",
		userLogin.NewHandler(server.config, server.appServices, server.sessionManager, server.logger).UserLoginEmail,
	)
	server.router.HandleFunc(
		apiContext+"/register",
		userRegister.NewHandler(server.config, server.appServices, server.sessionManager, server.logger).UserRegister,
	)

	// Post endpoints
	createPostHandler := post.NewCreatePostHandler(server.appServices.PostServices.Queries.CreatePost)
	getPostsHandler := post.NewGetPostsHandler(
		server.appServices.PostServices.Queries.GetAllPosts,
		server.appServices.PostServices.Queries.GetPostByID,
		server.appServices.PostServices.Queries.GetPostsByCategory,
	)
	updatePostHandler := post.NewUpdatePostHandler(server.appServices.PostServices.Queries.UpdatePost)

	server.router.HandleFunc(apiContext+"/posts", createPostHandler.CreatePost)
	server.router.HandleFunc(apiContext+"/posts/all", getPostsHandler.GetAllPosts)
	server.router.HandleFunc(apiContext+"/posts/get", getPostsHandler.GetPostByID)
	server.router.HandleFunc(apiContext+"/posts/update", updatePostHandler.UpdatePost)

	// Comment endpoints
	createCommentHandler := comment.NewCreateCommentHandler(server.appServices.CommentServices.Queries.CreateComment)
	getCommentsHandler := comment.NewGetCommentsHandler(
		server.appServices.CommentServices.Queries.GetCommentsByPost,
		server.appServices.CommentServices.Queries.GetCommentTree,
	)
	updateCommentHandler := comment.NewUpdateCommentHandler(server.appServices.CommentServices.Queries.UpdateComment)

	server.router.HandleFunc(apiContext+"/comments", createCommentHandler.CreateComment)
	server.router.HandleFunc(apiContext+"/comments/post", getCommentsHandler.GetCommentsByPost)
	server.router.HandleFunc(apiContext+"/comments/tree", getCommentsHandler.GetCommentTree)
	server.router.HandleFunc(apiContext+"/comments/update", updateCommentHandler.UpdateComment)

	// Vote endpoints
	voteHandler := vote.NewVoteHandler(
		server.appServices.VoteServices.Queries.CastVote,
		server.appServices.VoteServices.Queries.GetVoteStatus,
		server.appServices.VoteServices.Queries.GetUserVotes,
	)

	server.router.HandleFunc(apiContext+"/votes/cast", voteHandler.CastVote)
	server.router.HandleFunc(apiContext+"/votes/status", voteHandler.GetVoteStatus)
	server.router.HandleFunc(apiContext+"/votes/user", voteHandler.GetUserVotes)

	// Category endpoints
	categoryHandler := category.NewCategoryHandler(
		server.appServices.CategoryServices.Queries.CreateCategory,
		server.appServices.CategoryServices.Queries.GetAllCategories,
		server.appServices.CategoryServices.Queries.GetCategoryByID,
		server.appServices.CategoryServices.Queries.UpdateCategory,
		server.appServices.CategoryServices.Queries.DeleteCategory,
		server.appServices.CategoryServices.Queries.GetCategoryWithPosts,
	)

	server.router.HandleFunc(apiContext+"/categories", categoryHandler.CreateCategory)
	server.router.HandleFunc(apiContext+"/categories/all", categoryHandler.GetAllCategories)
	server.router.HandleFunc(apiContext+"/categories/get", categoryHandler.GetCategoryByID)
	server.router.HandleFunc(apiContext+"/categories/update", categoryHandler.UpdateCategory)
	server.router.HandleFunc(apiContext+"/categories/delete", categoryHandler.DeleteCategory)
	server.router.HandleFunc(apiContext+"/categories/posts", categoryHandler.GetCategoryWithPosts)
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
	server.sessionManager = session.NewSessionManager(server.db, server.config.SessionManager)
}
