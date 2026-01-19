package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextIP string

const userIP contextIP = "ip"

func GetClientIPMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		ctx := context.WithValue(r.Context(), userIP, ip)
		r.Header.Set("X-Forwarded-For", ip)
		r.Header.Set("X-Real-IP", ip)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetIPFromContext(r *http.Request) string {
	ip := r.Context().Value(userIP)
	if ip == nil {
		return ""
	}
	ipStr, ok := ip.(string)
	if !ok {
		return ""
	}

	return ipStr
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
