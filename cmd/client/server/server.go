package server

import (
	"log"
	"net/http"
	"net/http/cookiejar"

	"github.com/arnald/forum/cmd/client/config"
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

	// Homepage
	cs.Router.HandleFunc("/", cs.HomePage)

	// Register page
	cs.Router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cs.RegisterPage(w, r)
		case http.MethodPost:
			cs.RegisterPost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Login page
	cs.Router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cs.LoginPage(w, r)
		case http.MethodPost:
			cs.LoginPost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
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
