package helpers

import (
	"net/http"
)

func SetIPHeaders(r *http.Request, ip string) {
	r.Header.Set("X-Forwarded-For", ip)
	r.Header.Set("X-Real-IP", ip)
}
