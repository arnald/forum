package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/arnald/forum/internal/infra/middleware/ratelimiter"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type rateLimitMiddleware struct {
	limiter *ratelimiter.RateLimiter
	handler http.Handler
}

func NewRateLimiterMiddleware(handler http.Handler, limit int, windowSeconds int64, cleanup time.Duration) http.Handler {
	return &rateLimitMiddleware{
		limiter: ratelimiter.NewRateLimiter(limit, windowSeconds, cleanup),
		handler: handler,
	}
}

func (rl *rateLimitMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := getClientIP(r)

	allowed, remaining, resetTime := rl.limiter.Allow(ip)

	w.Header().Set("X-Rateimit-Limit", strconv.Itoa(rl.limiter.Limit))
	w.Header().Set("X-Rateimit-Remaining", strconv.Itoa(remaining))
	w.Header().Set("X-Rateimit-Reset", strconv.FormatInt(resetTime, 10))

	if !allowed {
		retryAfter := resetTime - time.Now().Unix()
		w.Header().Set("Retry-After", strconv.FormatInt(retryAfter, 10))

		helpers.RespondWithError(
			w,
			http.StatusTooManyRequests,
			"Rate limit exceeded, try again later",
		)

		return
	}

	rl.handler.ServeHTTP(w, r)
}

func getClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	ip := r.RemoteAddr

	idx := strings.LastIndex(ip, ":")
	if idx != -1 {
		ip = ip[:idx]
	}

	return ip
}
