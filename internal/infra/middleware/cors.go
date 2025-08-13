package middleware

import "net/http"

type corsMiddleware struct {
	handler http.Handler
}

func (c *corsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	c.handler.ServeHTTP(w, r)
}

func NewCorsMiddleware(handler http.Handler) http.Handler {
	return &corsMiddleware{handler: handler}
}
