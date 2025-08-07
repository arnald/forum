package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/uuid"
)

const contextTimeout = 15 * time.Second

type CreateSessionRequest struct {
	UserID    string
	IPAddress string
	Token     string
}

type Manager struct {
	db             *sql.DB
	tokenGenerator tokenGenerator
	sessionConfig  config.SessionManagerConfig
}

func NewSessionManager(db *sql.DB, sessionConfig config.SessionManagerConfig) user.SessionManager {
	return &Manager{
		db:             db,
		sessionConfig:  sessionConfig,
		tokenGenerator: uuid.NewProvider(),
	}
}

type tokenGenerator interface {
	NewUUID() string
}

func (sm *Manager) CreateSession(ctx context.Context, userID string) (*user.Session, error) {
	query := `
	INSERT INTO sessions (token, user_id, expires_at, refresh_token, refresh_token_expires_at)
	VALUES (?, ?, ?, ?, ?)`

	stmt, err := sm.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	newSessionToken := sm.tokenGenerator.NewUUID()
	newrefreshToken := sm.tokenGenerator.NewUUID()

	expiry := time.Now().Add(sm.sessionConfig.DefaultExpiry)
	refreshExpiry := expiry.Add(sm.sessionConfig.RefreshTokenExpiry)

	_, err = stmt.ExecContext(
		ctx,
		newSessionToken,
		userID,
		expiry.Format("2006-01-02 15:04:05"),

		newrefreshToken,
		refreshExpiry.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return nil, err
	}

	session := &user.Session{
		AccessToken:        newSessionToken,
		UserID:             userID,
		Expiry:             expiry,
		RefreshToken:       newrefreshToken,
		RefreshTokenExpiry: refreshExpiry,
	}

	return session, nil
}

func (sm *Manager) GetSession(sessionID string) (*user.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	query := `SELECT token, user_id, expires_at FROM sessions WHERE token = ?`

	stmt, err := sm.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, sessionID)

	var session user.Session

	err = row.Scan(&session.AccessToken, &session.UserID, &session.Expiry)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	if session.Expiry.Before(time.Now()) {
		_ = sm.DeleteSession(sessionID)
		return nil, ErrSessionExpired
	}
	return &session, nil
}

func (sm *Manager) GetUserFromSession(sessionID string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	query := `
    SELECT 
        u.id,
        u.email,
        u.username,
        u.created_at,
        u.avatar_url,
        u.password_hash
    FROM users u
    INNER JOIN sessions s ON s.user_id = u.id
    WHERE s.token = ? AND s.expires_at > CURRENT_TIMESTAMP
	`

	stmt, err := sm.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, sessionID)

	var User user.User

	err = row.Scan(
		&User.ID,
		&User.Email,
		&User.Username,
		&User.CreatedAt,
		&User.AvatarURL,
		&User.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &User, nil
}

func (sm *Manager) DeleteSession(sessionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	query := `DELETE FROM sessions WHERE id = ?`

	stmt, err := sm.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, sessionID)
	return err
}

func (sm *Manager) NewSessionCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     sm.sessionConfig.CookieName,
		Value:    token,
		Path:     sm.sessionConfig.CookiePath,
		Domain:   sm.sessionConfig.CookieDomain,
		HttpOnly: sm.sessionConfig.HTTPOnlyCookie,
		Secure:   sm.sessionConfig.SecureCookie,
		SameSite: parseSameSite(sm.sessionConfig.SameSite),
		MaxAge:   int(sm.sessionConfig.DefaultExpiry.Seconds()),
	}
}

func (sm *Manager) ValidateSession(sessionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	query := `
        SELECT expires_at 
        FROM sessions 
        WHERE token = ?
    `

	stmt, err := sm.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	var expiresAt time.Time
	err = stmt.QueryRowContext(ctx, sessionID).Scan(&expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSessionNotFound
		}
		return fmt.Errorf("scanning session failed: %w", err)
	}

	if expiresAt.Before(time.Now()) {
		_ = sm.DeleteSession(sessionID)
		return ErrSessionExpired
	}

	return nil
}

func parseSameSite(s string) http.SameSite {
	switch s {
	case "Strict":
		return http.SameSiteStrictMode
	case "None":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}
