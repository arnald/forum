package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

const (
	contextTimeout = 10 * time.Second
)

// Logout handles user logout by clearing session cookies and backend session.
func (cs *ClientServer) Logout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), contextTimeout)
	defer cancel()

	err := cs.logoutFromBackend(ctx, r)
	if err != nil {
		log.Printf("Failed to logout from backend: %v", err)
	}

	cs.clearSessionCookies(w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// logoutFromBackend calls the backend logout endpoint to delete the session.
func (cs *ClientServer) logoutFromBackend(ctx context.Context, r *http.Request) error {
	logoutReq, err := http.NewRequestWithContext(ctx, http.MethodPost, backendLogoutURL, nil)
	if err != nil {
		return err
	}

	for _, cookie := range r.Cookies() {
		logoutReq.AddCookie(cookie)
	}

	resp, err := cs.HTTPClient.Do(logoutReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Backend logout returned status: %d", resp.StatusCode)
		return backendError("logout failed")
	}

	return nil
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
