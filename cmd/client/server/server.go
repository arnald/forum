package server

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"

	"github.com/arnald/forum/cmd/client/config"
	"github.com/arnald/forum/cmd/client/middleware"
	"github.com/arnald/forum/internal/pkg/path"
)

// ClientServer represents the frontend client server.
type ClientServer struct {
	Config     *config.Client
	Router     *http.ServeMux
	HTTPClient *http.Client
}

// NewClientServer creates and initializes a new ClientServer.
func NewClientServer(cfg *config.Client) (*ClientServer, error) {
	// Create a cookie jar to persist cookies between requests
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	// Create HTTP client with cookie jar
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: cfg.HTTPTimeouts.Read,
	}

	return &ClientServer{
		Config:     cfg,
		Router:     http.NewServeMux(),
		HTTPClient: httpClient,
	}, nil
}

// SetupRoutes configures all HTTP routes and binds them to handler methods.
func (cs *ClientServer) SetupRoutes() {
	resolver := path.NewResolver()

	// Static file serving
	cs.Router.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir(resolver.GetPath("frontend/static/")))),
	)

	// Create auth middleware
	authMiddleware := middleware.AuthMiddleware(cs.HTTPClient)

	// Public Routes (with optional auth - shows user if logged in).
	// Homepage
	cs.Router.HandleFunc("/", applyMiddleware(cs.HomePage, authMiddleware))

	// Categories page
	cs.Router.HandleFunc("/categories", applyMiddleware(cs.CategoriesPage, authMiddleware))

	// Topics page
	cs.Router.HandleFunc("/topics", applyMiddleware(cs.TopicsPage, authMiddleware))

	// Topic detail page
	cs.Router.HandleFunc("/topic/", applyMiddleware(cs.TopicPage, authMiddleware))

	// Register page
	cs.Router.HandleFunc("/register",
		applyMiddleware(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				cs.RegisterPage(w, r)
			case http.MethodPost:
				cs.RegisterPost(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}, authMiddleware))

	// Login page
	cs.Router.HandleFunc("/login",
		applyMiddleware(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				cs.LoginPage(w, r)
			case http.MethodPost:
				cs.LoginPost(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}, authMiddleware))

	// OAuth Register
	cs.Router.HandleFunc("/auth/github/login", applyMiddleware(cs.GitHubRegister, authMiddleware))
	cs.Router.HandleFunc("/auth/google/login", applyMiddleware(cs.GoogleRegister, authMiddleware))
	cs.Router.HandleFunc("/auth/callback", applyMiddleware(cs.Callback, authMiddleware))

	// Protected Routes (require authentication).
	// Logout route - clears cookies
	cs.Router.HandleFunc("/logout", applyMiddleware(cs.Logout, middleware.RequireAuth, authMiddleware))
}

// ListenAndServe starts the HTTP server.
func (cs *ClientServer) ListenAndServe() error {
	server := &http.Server{
		Addr:              ":" + cs.Config.Port,
		Handler:           cs.Router,
		ReadHeaderTimeout: cs.Config.HTTPTimeouts.ReadHeader,
		ReadTimeout:       cs.Config.HTTPTimeouts.Read,
		WriteTimeout:      cs.Config.HTTPTimeouts.Write,
		IdleTimeout:       cs.Config.HTTPTimeouts.Idle,
	}

	log.Printf("Client started on port: %s (%s environment)", cs.Config.Port, cs.Config.Environment)
	return server.ListenAndServe()
}

func applyMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

// Standardized way to make requests to the backend server, used in handlers.
func (cs *ClientServer) newRequest(ctx context.Context, method string, url string, req any) (*http.Response, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, backendError("Failed to marshal request: " + err.Error())
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, backendError("Failed to create request:" + err.Error())
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := cs.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, backendError("Registration request failed: " + err.Error())
	}

	return resp, nil
}
