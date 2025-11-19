package server

import (
	"net/http"

	"github.com/arnald/forum/cmd/client/domain"
)

// LoginPage handles GET requests to /login.
func (cs *ClientServer) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	renderTemplate(w, "login", domain.LoginFormErrors{})
}

func (cs *ClientServer) LoginPost(w http.ResponseWriter, r *http.Request) {}
