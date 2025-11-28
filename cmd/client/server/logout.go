package server

import "net/http"

// Logout handles user logout by clearing session cookies.
func (cs *ClientServer) Logout(w http.ResponseWriter, r *http.Request) {
	cs.clearSessionCookies(w)

	// Redirect to homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// clearSessionCookies clears both session cookies (for logout).
func (cs *ClientServer) clearSessionCookies(w http.ResponseWriter) {
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
}
