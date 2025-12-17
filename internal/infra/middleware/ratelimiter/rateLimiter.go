package ratelimiter

import (
	"sync"
	"time"
)

type ClientInfo struct {
	currentWindow  WindowInfo
	previousWindow WindowInfo
}

type WindowInfo struct {
	count     int
	startTime int64
}

type RateLimiter struct {
	clients         map[string]*ClientInfo
	mu              sync.RWMutex
	Limit           int
	windowSize      int64
	cleanupInterval time.Duration
}

func NewRateLimiter(limit int, windowSeconds int64, cleanup time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients:         make(map[string]*ClientInfo),
		Limit:           limit,
		windowSize:      windowSeconds,
		cleanupInterval: cleanup,
	}

	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) Allow(ip string) (bool, int, int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now().Unix()

	client, exists := rl.clients[ip]
	if !exists {
		client = &ClientInfo{
			currentWindow: WindowInfo{
				count:     0,
				startTime: now,
			},
		}
		rl.clients[ip] = client
	}

	currentWindowStart := client.currentWindow.startTime
	windowElapsed := now - currentWindowStart

	if windowElapsed >= rl.windowSize {
		client.previousWindow = client.currentWindow

		client.currentWindow = WindowInfo{
			count:     0,
			startTime: now,
		}

		currentWindowStart = now
		windowElapsed = 0
	}

	previousWeight := float64(rl.windowSize-windowElapsed) / float64(rl.windowSize)
	weightedPreviousCount := float64(client.previousWindow.count) * previousWeight
	totalCount := int(weightedPreviousCount) + client.currentWindow.count

	allowed := totalCount < rl.Limit

	if allowed {
		client.currentWindow.count++
	}

	remaining := rl.Limit - totalCount - 1
	if remaining < 0 {
		remaining = 0
	}

	resetTime := currentWindowStart + rl.windowSize

	return allowed, remaining, resetTime
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now().Unix()

		for ip, client := range rl.clients {
			if now-client.currentWindow.startTime > rl.windowSize*2 {
				delete(rl.clients, ip)
			}
		}

		rl.mu.Unlock()
	}
}
