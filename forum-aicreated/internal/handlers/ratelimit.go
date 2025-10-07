// Package handlers - ratelimit.go implements rate limiting to prevent abuse.
// This file provides request rate limiting based on client IP addresses to prevent
// spam, brute force attacks, and DoS attempts. Uses in-memory storage with automatic
// cleanup of expired entries and configurable limits per time window.
package handlers

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimiter tracks request counts per client IP address
// Uses a map protected by RWMutex for thread-safe concurrent access
// Implements a sliding window rate limiting algorithm
type RateLimiter struct {
	clients map[string]*ClientLimiter // Maps IP addresses to their request counters
	mutex   sync.RWMutex              // Protects concurrent access to clients map
}

// ClientLimiter tracks an individual client's request rate
// Each client has a request counter and a reset timestamp
type ClientLimiter struct {
	requests  int       // Number of requests made in current window
	resetTime time.Time // When the rate limit counter resets
}

// Rate limiting configuration constants
// These values define how many requests are allowed per time window
const (
	maxRequests = 100        // Maximum requests allowed per window
	resetWindow = time.Minute // Time window duration (1 minute)
)

// Global rate limiter instance
// Single instance shared across all requests for consistent rate limiting
var rateLimiter = &RateLimiter{
	clients: make(map[string]*ClientLimiter),
}

// cleanupExpired removes expired client entries from memory
// This prevents memory leaks by periodically cleaning up old rate limit data
// Should be called periodically to maintain memory efficiency
func (rl *RateLimiter) cleanupExpired() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	// Remove entries where reset time has passed
	for ip, client := range rl.clients {
		if now.After(client.resetTime) {
			delete(rl.clients, ip)
		}
	}
}

// isAllowed checks if a client IP is allowed to make a request
// Implements sliding window rate limiting with automatic reset
// Returns true if request is allowed, false if rate limit exceeded
func (rl *RateLimiter) isAllowed(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	// Check if this is a new client
	client, exists := rl.clients[ip]
	if !exists {
		// First request from this IP - create new limiter entry
		rl.clients[ip] = &ClientLimiter{
			requests:  1,
			resetTime: now.Add(resetWindow),
		}
		return true
	}

	// Check if the rate limit window has expired
	if now.After(client.resetTime) {
		// Reset the counter for a new window
		client.requests = 1
		client.resetTime = now.Add(resetWindow)
		return true
	}

	// Check if client has exceeded rate limit
	if client.requests >= maxRequests {
		return false // Rate limit exceeded
	}

	// Increment request counter and allow
	client.requests++
	return true
}

// getClientIP extracts the real client IP address from an HTTP request
// Handles proxy headers (X-Forwarded-For, X-Real-IP) for accurate identification
// This is important for rate limiting behind reverse proxies or load balancers
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (common with proxies/load balancers)
	// Format: "client, proxy1, proxy2" - we want the first IP (original client)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header (simpler proxy header)
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fallback to direct connection IP (RemoteAddr)
	// Format: "ip:port" - we only need the IP part
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	return ip
}

// RateLimitMiddleware wraps an HTTP handler with rate limiting logic
// Checks client IP against rate limits before allowing request to proceed
// Returns 429 Too Many Requests if limit exceeded
func RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Periodically clean up expired entries (runs in background)
		// This prevents memory buildup from old client entries
		go rateLimiter.cleanupExpired()

		// Extract client IP for rate limiting
		ip := getClientIP(r)

		// Check if this client is allowed to make a request
		if !rateLimiter.isAllowed(ip) {
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}

		// Rate limit check passed - proceed with request
		next.ServeHTTP(w, r)
	}
}

// RateLimitedHandler is a convenience method for applying rate limiting to handlers
// Wraps the provided handler with rate limiting middleware
func (h *Handler) RateLimitedHandler(handler http.HandlerFunc) http.HandlerFunc {
	return RateLimitMiddleware(handler)
}
