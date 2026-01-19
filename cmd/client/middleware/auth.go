package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/helpers"
)

type contextKey string

const (
	userContextKey contextKey = "user"
	backendMeURL   string     = "http://localhost:8080/api/v1/me"
)

var (
	ErrUserNotAuthorized = errors.New("user not authorized")
	ErrTooManyRequests   = errors.New("too many requests")
)

// AuthMiddleware wraps a handler and injects authenticated user data into context.
func AuthMiddleware(httpClient *http.Client) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Try to get user from /me endpoint.
			user, err := getCurrentUser(ctx, httpClient, r)
			if err == nil && user != nil {
				// User authenticated, add to context.
				ctx = context.WithValue(ctx, userContextKey, user)
				// TO DO: CHECK FOR 429 IN EVERY BACKEND RESPONSE (IN EACH HANDLER) e.x LIKE BELOW
			} else if errors.Is(err, ErrTooManyRequests) {
				http.Error(w, err.Error(), http.StatusTooManyRequests)
				return
			}
			// If error or no user, continue without user context.
			// This allows optional auth for certain pages.

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// getCurrentUser fetches the current user from the backend /me endpoint.
func getCurrentUser(ctx context.Context, httpClient *http.Client, r *http.Request) (*domain.LoggedInUser, error) {
	// Create a new request to the backend /me endpoint.
	meReq, err := http.NewRequestWithContext(ctx, http.MethodGet, backendMeURL, nil)
	if err != nil {
		log.Printf("Failed to create /me request: %v", err)
		return nil, err
	}

	// Copy cookies from the original request to the /me request.
	// This includes access_token and refresh_token cookies.
	for _, cookie := range r.Cookies() {
		meReq.AddCookie(cookie)
	}

	meReq.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(meReq)
	if err != nil {
		log.Printf("Failed to fetch /me: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// If not authorized, return nil (user not authenticated).
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUserNotAuthorized
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrTooManyRequests
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status from /me: %d", resp.StatusCode)
		return nil, err
	}

	var meResp domain.BackendMeResponse
	err = helpers.DecodeBackendResponse(resp, &meResp)
	if err != nil {
		log.Printf("Failed to decode /me response: %v", err)
		return nil, err
	}

	// Convert backend response to LoggedInUser domain model.
	user := &domain.LoggedInUser{
		ID:        meResp.ID,
		Username:  meResp.Username,
		Email:     meResp.Email,
		AvatarURL: meResp.AvatarURL,
	}

	return user, nil
}

// GetUserFromContext retrieves the authenticated user from context.
func GetUserFromContext(ctx context.Context) *domain.LoggedInUser {
	user, ok := ctx.Value(userContextKey).(*domain.LoggedInUser)
	if !ok {
		return nil
	}
	return user
}
