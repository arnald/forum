package server

import (
	"log"
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

	http.Redirect(w, r, cs.BackendURLs.GithubRegisterURL(), http.StatusTemporaryRedirect)
}

func (cs *ClientServer) Callback(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("access_token")
	refreshToken := r.URL.Query().Get("refresh_token")

	if accessToken == "" || refreshToken == "" {
		log.Printf("OAuth callback missing tokens - access_token: %v, refresh_token: %v", accessToken != "", refreshToken != "")
		http.Error(w, "OAuth session not found", http.StatusBadRequest)
		return
	}

	log.Printf("Setting OAuth cookies for callback - access_token present: %v, refresh_token present: %v", accessToken != "", refreshToken != "")

	// Set cookies before redirect
	cs.setSessionCookies(w, accessToken, refreshToken)

	log.Println("Redirecting to homepage after OAuth callback")
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

	http.Redirect(w, r, cs.BackendURLs.GoogleRegisterURL(), http.StatusTemporaryRedirect)
}
