package server

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"

	"github.com/arnald/forum/cmd/client/config"
	"github.com/arnald/forum/cmd/client/helpers"
	"github.com/arnald/forum/cmd/client/middleware"
	"github.com/arnald/forum/internal/pkg/path"
)

// ClientServer represents the frontend client server.
type ClientServer struct {
	Config      *config.Client
	Router      *http.ServeMux
	HTTPClient  *http.Client
	SseClient   *http.Client
	BackendURLs *BackendURLs
}

// getSecureTLSConfig returns a TLS configuration with explicit cipher suites.
func getSecureTLSConfig(isDevelopment bool) *tls.Config {
	return &tls.Config{
		// Only skip verification in development
		InsecureSkipVerify: isDevelopment, //nolint:gosec

		// Minimum TLS version - TLS 1.2 is widely supported and secure
		MinVersion: tls.VersionTLS12,

		// Prefer server cipher suites for better security
		PreferServerCipherSuites: true,

		// Explicit cipher suites - prioritize modern AEAD ciphers
		CipherSuites: []uint16{
			// TLS 1.3 cipher suites (automatically used when TLS 1.3 is negotiated)
			// Note: TLS 1.3 suites are not configurable via CipherSuites in Go

			// TLS 1.2 AEAD cipher suites (recommended)
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,

			// Fallback cipher suites (still secure but less preferred)
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		},

		// Curve preferences for ECDHE
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
	}
}

// NewClientServer creates and initializes a new ClientServer.
func NewClientServer(cfg *config.Client) (*ClientServer, error) {
	// Create a cookie jar to persist cookies between requests
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	// Get TLS configuration with explicit cipher suites
	tlsConfig := getSecureTLSConfig(cfg.Environment == "development")

	// Create custom transport that skips TLS verification for self-signed certs in development
	transport := &http.Transport{
		TLSClientConfig:   tlsConfig,
		ForceAttemptHTTP2: false, // Disable HTTP/2 to avoid protocol issues
	}

	// Log TLS configuration
	if cfg.Environment == "development" {
		log.Printf("TLS Config: MinVersion=TLS1.2, InsecureSkipVerify=true (development mode)")
	} else {
		log.Printf("TLS Config: MinVersion=TLS1.2, InsecureSkipVerify=false, CipherSuites=%d configured", len(tlsConfig.CipherSuites))
	}

	// Create HTTP client with cookie jar and custom transport
	httpClient := &http.Client{
		Jar:       jar,
		Timeout:   cfg.HTTPTimeouts.Read,
		Transport: transport,
	}

	sseClient := &http.Client{
		Timeout:   0,
		Transport: transport,
	}

	// Create backend URLs instance
	backendURLs := NewBackendURLs(cfg.BackendURL)

	return &ClientServer{
		Config:      cfg,
		Router:      http.NewServeMux(),
		HTTPClient:  httpClient,
		SseClient:   sseClient,
		BackendURLs: backendURLs,
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
	authMiddleware := middleware.AuthMiddleware(cs.HTTPClient, cs.BackendURLs.MeURL())

	// Public Routes (with optional auth - shows user if logged in).
	// Homepage
	cs.Router.HandleFunc("/", applyMiddleware(cs.HomePage, authMiddleware))

	// Categories page
	cs.Router.HandleFunc("/categories", applyMiddleware(cs.CategoriesPage, authMiddleware))

	// Topics page
	cs.Router.HandleFunc("/topics", applyMiddleware(cs.TopicsPage, authMiddleware))

	// Topic detail page
	cs.Router.HandleFunc("/topic/", applyMiddleware(cs.TopicPage, authMiddleware))

	// Topic CRUD routes
	cs.Router.HandleFunc("/topics/create", applyMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cs.CreateTopicPage(w, r)
		case http.MethodPost:
			cs.CreateTopicPost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}, middleware.RequireAuth, authMiddleware))
	cs.Router.HandleFunc("/topics/edit", applyMiddleware(cs.UpdateTopicPost, middleware.RequireAuth, authMiddleware))
	cs.Router.HandleFunc("/topics/delete", applyMiddleware(cs.DeleteTopicPost, middleware.RequireAuth, authMiddleware))

	// Comment CRUD routes
	cs.Router.HandleFunc("/comments/create", applyMiddleware(cs.CreateCommentPost, middleware.RequireAuth, authMiddleware))
	cs.Router.HandleFunc("/comments/edit", applyMiddleware(cs.UpdateCommentPost, middleware.RequireAuth, authMiddleware))
	cs.Router.HandleFunc("/comments/delete", applyMiddleware(cs.DeleteCommentPost, middleware.RequireAuth, authMiddleware))

	// Vote API routes (these are API endpoints, not pages)
	cs.Router.HandleFunc("/api/vote/cast", applyMiddleware(cs.CastVote, middleware.RequireAuth, authMiddleware))
	cs.Router.HandleFunc("/api/vote/counts", applyMiddleware(cs.GetVoteCounts, authMiddleware))
	cs.Router.HandleFunc("/api/vote/delete", applyMiddleware(cs.DeleteVote, middleware.RequireAuth, authMiddleware))

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
	// Activity page
	cs.Router.HandleFunc("/activity", applyMiddleware(cs.ActivityPage, middleware.RequireAuth, authMiddleware))
	// Notification routes
	cs.Router.HandleFunc("/api/notifications/stream", applyMiddleware(cs.StreamNotifications, middleware.RequireAuth, authMiddleware))
	cs.Router.HandleFunc("/api/notifications", applyMiddleware(cs.GetNotifications, middleware.RequireAuth, authMiddleware))
	cs.Router.HandleFunc("/api/notifications/unread-count", applyMiddleware((cs.GetUnreadCount), middleware.RequireAuth, authMiddleware))
	cs.Router.HandleFunc("/api/notifications/mark-read", applyMiddleware(cs.MarkNotificationAsRead, middleware.RequireAuth, authMiddleware))
	cs.Router.HandleFunc("/api/notifications/mark-all-read", applyMiddleware(cs.MarkAllNotificationsAsRead, middleware.RequireAuth, authMiddleware))
	// Logout route - clears cookies
	cs.Router.HandleFunc("/logout", applyMiddleware(cs.Logout, middleware.RequireAuth, authMiddleware))
}

// ListenAndServe starts the HTTP server.
func (cs *ClientServer) ListenAndServe() error {
	handler := middleware.GetClientIPMiddleware(cs.Router)

	// Get TLS configuration for the server
	var tlsConfig *tls.Config
	if cs.Config.TLSCertFile != "" && cs.Config.TLSKeyFile != "" {
		tlsConfig = getSecureTLSConfig(false) // Server should never skip verification
	}

	server := &http.Server{
		Addr:              ":" + cs.Config.Port,
		Handler:           handler,
		ReadHeaderTimeout: cs.Config.HTTPTimeouts.ReadHeader,
		ReadTimeout:       cs.Config.HTTPTimeouts.Read,
		WriteTimeout:      cs.Config.HTTPTimeouts.Write,
		IdleTimeout:       cs.Config.HTTPTimeouts.Idle,
		TLSConfig:         tlsConfig,
	}

	log.Printf("Client started on port: %s (%s environment)", cs.Config.Port, cs.Config.Environment)

	if cs.Config.TLSCertFile != "" && cs.Config.TLSKeyFile != "" {
		log.Printf("Starting HTTPS client with TLS certificates")
		log.Printf("TLS: MinVersion=TLS1.2, %d cipher suites configured", len(tlsConfig.CipherSuites))
		return server.ListenAndServeTLS(cs.Config.TLSCertFile, cs.Config.TLSKeyFile)
	}

	log.Printf("Starting HTTP client (no TLS)")
	return server.ListenAndServe()
}

func applyMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

// Standardized way to make requests to the backend server, used in handlers.
func (cs *ClientServer) newRequest(ctx context.Context, method string, url string, req any, ip string) (*http.Response, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, backendError("Failed to marshal request: " + err.Error())
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, backendError("Failed to create request:" + err.Error())
	}

	httpReq.Header.Set("Content-Type", "application/json")
	helpers.SetIPHeaders(httpReq, ip)

	resp, err := cs.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, backendError("Registration request failed: " + err.Error())
	}

	return resp, nil
}

// Makes a backend request and includes cookies from the original request, necessary for authenticated endpoints.
func (cs *ClientServer) newRequestWithCookies(ctx context.Context, method string, url string, req any, originalReq *http.Request) (*http.Response, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, backendError("Failed to marshal request: " + err.Error())
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, backendError("Failed to create request: " + err.Error())
	}

	httpReq.Header.Set("Content-Type", "application/json")

	ip := middleware.GetIPFromContext(originalReq)
	if ip == "" {
		return nil, backendError("No IP found in request")
	}

	helpers.SetIPHeaders(httpReq, ip)

	for _, cookie := range originalReq.Cookies() {
		httpReq.AddCookie(cookie)
	}

	resp, err := cs.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, backendError("Backend request failed: " + err.Error())
	}

	return resp, nil
}
