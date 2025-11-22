package server

import "net/http"

// Logout handles user logout by clearing session cookies
func (cs *ClientServer) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear access token cookie
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // Delete the cookie
	}

	// Clear refresh token cookie
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // Delete the cookie
	}

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	// Redirect to homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
