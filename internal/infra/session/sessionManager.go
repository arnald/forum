package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/uuid"
)

type CreateSessionRequest struct {
	UserID    string
	IPAddress string
	Token     string
}

type SessionManager struct {
	db             *sql.DB
	sessionConfig  config.SessionManagerConfig
	tokenGenerator tokenGenerator
}

func NewSessionManager(db *sql.DB, sessionConfig config.SessionManagerConfig) *SessionManager {
	return &SessionManager{
		db:             db,
		sessionConfig:  sessionConfig,
		tokenGenerator: uuid.NewProvider(),
	}
}

type tokenGenerator interface {
	NewUUID() string
}

func (sm *SessionManager) CreateSession(userID string) (*user.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	INSERT INTO sessions (token, user_id, expires_at)
	VALUES (?, ?, ?)`

	stmt, err := sm.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	newSessionToken := sm.tokenGenerator.NewUUID()

	_, err = stmt.ExecContext(
		ctx,
		newSessionToken,
		userID,
		sm.sessionConfig.DefaultExpiry,
	)
	if err != nil {
		return nil, err
	}

	session := &user.Session{
		Token:  []byte(newSessionToken),
		UserID: userID,
	}

	return session, nil
}

func (sm *SessionManager) GetSession(sessionID string) (*user.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `SELECT token, user_id, expires_at, ip_address FROM sessions WHERE id = ?`

	stmt, err := sm.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, sessionID)

	var session user.Session

	err = row.Scan(&session.Token, &session.UserID, &session.Expiry, &session.IPAddress)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if session.Expiry.Before(time.Now()) {
		_ = sm.DeleteSession(sessionID)
		return nil, nil

	}
	return &session, nil
}

func (sm *SessionManager) DeleteSession(sessionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
