package oauthrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arnald/forum/internal/domain/oauth"
	"github.com/arnald/forum/internal/domain/user"
)

type Repo struct {
	db *sql.DB
}

func NewOAuthRepository(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetUserByProviderID(ctx context.Context, provider oauth.Provider, providerUserID string) (*user.User, error) {
	query := `
	SELECT u.id, u.username, u.email, u.password_hash, u.created_at
	FROM users u
	INNER JOIN oauth_providers op ON u.id = op.user_id
	WHERE op.provider = ? AND op.provider_user_id = ?
	`

	var u user.User
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, string(provider), providerUserID).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repo) CreateOAuthUser(ctx context.Context, oauthUser *oauth.User) (*user.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	insertUserQuery := `
        INSERT INTO users (id, username, email, password_hash)
        VALUES (?, ?, ?, '')
    `

	_, err = tx.ExecContext(ctx, insertUserQuery,
		oauthUser.UserId,
		oauthUser.Username,
		oauthUser.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	insertUserQuery = `
	INSERT INTO oauth_providers (user_id, provider, provider_user_id, email, username, avatar_url)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = tx.ExecContext(ctx, insertUserQuery,
		oauthUser.UserId,
		string(oauthUser.Provider),
		oauthUser.ProviderID,
		oauthUser.Email,
		oauthUser.Username,
		oauthUser.AvatarURL,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &user.User{
		ID:       oauthUser.UserId,
		Username: oauthUser.Username,
		Email:    oauthUser.Email,
		Password: "",
	}, nil
}

func (r *Repo) LinkOAuthProvider(ctx context.Context, userID string, oauthUser *oauth.User) error {
	query := `
	INSERT INTO oath_providers (user_id, provider, provider_user_id, email, username, avatar_url, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?,)
	ON CONFLICT(provider, provider_user_id) DO UPDATE SET
		email = excluded.username,
		avatar_url = excluded.avatar_url,
		updated_at = excluded.updated_at
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		userID,
		string(oauthUser.Provider),
		oauthUser.Provider,
		oauthUser.Email,
		oauthUser.Username,
		oauthUser.AvatarURL,
	)

	return err
}

func (r *Repo) GetOAuthProvider(ctx context.Context, userID string, provider oauth.Provider) (*oauth.User, error) {
	query := `
		SELECT provider_user_id, email, username, avatar_url
		FROM oauth_providers
		WHERE user_id = ? AND provider = ?
	`

	var oauthUser oauth.User
	oauthUser.Provider = provider

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, userID, string(provider)).Scan(
		ctx,
		&oauthUser.ProviderID,
		&oauthUser.Email,
		&oauthUser.Username,
		&oauthUser.AvatarURL,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &oauthUser, nil
}
