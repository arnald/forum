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

type RateLimiter struct {
	clients map[string]*ClientLimiter
	mutex   sync.RWMutex
}

type ClientLimiter struct {
	requests  int
	resetTime time.Time
}

const (
	maxRequests = 100
	resetWindow = time.Minute
)

var rateLimiter = &RateLimiter{
	clients: make(map[string]*ClientLimiter),
}

func (rl *RateLimiter) cleanupExpired() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	for ip, client := range rl.clients {
		if now.After(client.resetTime) {
			delete(rl.clients, ip)
		}
	}
}

func (rl *RateLimiter) isAllowed(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	client, exists := rl.clients[ip]
	if !exists {
		rl.clients[ip] = &ClientLimiter{
			requests:  1,
			resetTime: now.Add(resetWindow),
		}
		return true
	}

	if now.After(client.resetTime) {
		client.requests = 1
		client.resetTime = now.Add(resetWindow)
		return true
	}

	if client.requests >= maxRequests {
		return false
	}

	client.requests++
	return true
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fallback to RemoteAddr
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	return ip
}

func RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clean up expired entries periodically
		go rateLimiter.cleanupExpired()

		ip := getClientIP(r)

		if !rateLimiter.isAllowed(ip) {
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (h *Handler) RateLimitedHandler(handler http.HandlerFunc) http.HandlerFunc {
	return RateLimitMiddleware(handler)
}
