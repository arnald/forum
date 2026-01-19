package server

import (
	"net/http"

	"github.com/arnald/forum/cmd/client/helpers"
	"github.com/arnald/forum/cmd/client/middleware"
)

func (cs *ClientServer) GitHubRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip := middleware.GetIPFromContext(r)
	if ip == "" {
		http.Error(w, "Error no IP found in request", http.StatusInternalServerError)
	}

	helpers.SetIPHeaders(r, ip)

	http.Redirect(w, r, backendGithubRegister, http.StatusTemporaryRedirect)
}

func (cs *ClientServer) Callback(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("access_token")
	refreshToken := r.URL.Query().Get("refresh_token")

	if accessToken == "" || refreshToken == "" {
		http.Error(w, "Github session not found", http.StatusBadRequest)
		return
	}

	cs.setSessionCookies(w, accessToken, refreshToken)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (cs *ClientServer) GoogleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip := middleware.GetIPFromContext(r)
	if ip == "" {
		http.Error(w, "Error no IP found in request", http.StatusInternalServerError)
	}

	helpers.SetIPHeaders(r, ip)

	http.Redirect(w, r, backendGooglebRegister, http.StatusTemporaryRedirect)
}
