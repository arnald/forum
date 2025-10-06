// Package auth handles all authentication and authorization logic for the forum.
// It provides secure user registration, login, session management, and permission checking.
// Uses bcrypt for password hashing and UUID-based sessions with expiration.
package auth

import (
	"errors"
	"forum/internal/database"
	"forum/internal/models"
	"net/http"
	"time"

	"github.com/gofrs/uuid"           // UUID generation for session IDs
	"golang.org/x/crypto/bcrypt"      // Secure password hashing
)

// Auth handles authentication operations and requires database access
type Auth struct {
	db *database.DB // Database connection for user and session operations
}

// NewAuth creates a new Auth instance with database dependency injection
func NewAuth(db *database.DB) *Auth {
	return &Auth{db: db}
}

// ===== PASSWORD SECURITY =====

// HashPassword securely hashes a plaintext password using bcrypt
// bcrypt automatically handles salt generation and is resistant to timing attacks
func (a *Auth) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword verifies a plaintext password against a bcrypt hash
// Returns true if the password matches, false otherwise
// Uses constant-time comparison to prevent timing attacks
func (a *Auth) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ===== USER CREATION =====

// CreateUser creates a new local user account with email/password authentication
// Automatically hashes the password and assigns the default 'user' role
func (a *Auth) CreateUser(email, username, password string) (*models.User, error) {
	// Hash the password before storing it
	hashedPassword, err := a.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Insert new user into database with default role
	query := `INSERT INTO users (email, username, password, role) VALUES (?, ?, ?, ?)`
	result, err := a.db.Exec(query, email, username, hashedPassword, models.RoleUser)
	if err != nil {
		return nil, err
	}

	// Get the auto-generated user ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Return the complete user object
	user := &models.User{
		ID:       int(id),
		Email:    email,
		Username: username,
		Password: hashedPassword,
		Role:     models.RoleUser,
	}

	return user, nil
}

// CreateOAuthUser creates a new user account from OAuth provider (Google, GitHub, etc.)
// These users don't have local passwords since they authenticate via external providers
func (a *Auth) CreateOAuthUser(email, username, provider, providerID string) (*models.User, error) {
	// Insert OAuth user with empty password and provider information
	query := `INSERT INTO users (email, username, password, role, provider, provider_id) VALUES (?, ?, '', ?, ?, ?)`
	result, err := a.db.Exec(query, email, username, models.RoleUser, provider, providerID)
	if err != nil {
		return nil, err
	}

	// Get the auto-generated user ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Return the complete OAuth user object
	user := &models.User{
		ID:         int(id),
		Email:      email,
		Username:   username,
		Role:       models.RoleUser,
		Provider:   provider,
		ProviderID: providerID,
	}

	return user, nil
}

// ===== USER RETRIEVAL =====

// GetUserByEmail retrieves a user from the database by their email address
// Used during login to find the user account for authentication
func (a *Auth) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, username, password, role, provider, provider_id, created_at FROM users WHERE email = ?`
	err := a.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Role, &user.Provider, &user.ProviderID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user from the database by their unique ID
// Used for session validation and retrieving current user information
func (a *Auth) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, username, password, role, provider, provider_id, created_at FROM users WHERE id = ?`
	err := a.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Role, &user.Provider, &user.ProviderID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByProvider retrieves a user by their OAuth provider information
// Used during OAuth login to find existing accounts linked to external providers
func (a *Auth) GetUserByProvider(provider, providerID string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, username, password, role, provider, provider_id, created_at FROM users WHERE provider = ? AND provider_id = ?`
	err := a.db.QueryRow(query, provider, providerID).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Role, &user.Provider, &user.ProviderID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ===== SESSION MANAGEMENT =====

// CreateSession creates a new user session with a random UUID and 24-hour expiration
// Called after successful login to maintain user authentication state
func (a *Auth) CreateSession(userID int) (*models.Session, error) {
	// Generate a unique session ID using UUID v4
	sessionID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	// Set session to expire in 24 hours for security
	expiresAt := time.Now().Add(24 * time.Hour)

	// Store session in database
	query := `INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`
	_, err = a.db.Exec(query, sessionID.String(), userID, expiresAt)
	if err != nil {
		return nil, err
	}

	// Return complete session object
	session := &models.Session{
		ID:        sessionID.String(),
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	return session, nil
}

// GetSession retrieves a session by ID and validates its expiration
// Automatically deletes expired sessions and returns an error if session is invalid
func (a *Auth) GetSession(sessionID string) (*models.Session, error) {
	session := &models.Session{}
	query := `SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = ?`
	err := a.db.QueryRow(query, sessionID).Scan(&session.ID, &session.UserID, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Check if session has expired and clean it up if so
	if time.Now().After(session.ExpiresAt) {
		a.DeleteSession(sessionID)
		return nil, errors.New("session expired")
	}

	return session, nil
}

// DeleteSession removes a session from the database
// Called during logout or when cleaning up expired sessions
func (a *Auth) DeleteSession(sessionID string) error {
	query := `DELETE FROM sessions WHERE id = ?`
	_, err := a.db.Exec(query, sessionID)
	return err
}

// GetUserFromRequest extracts the current user from an HTTP request
// Reads the session cookie, validates the session, and returns the user object
// Returns error if no valid session exists
func (a *Auth) GetUserFromRequest(r *http.Request) (*models.User, error) {
	// Extract session ID from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, errors.New("no session found")
	}

	// Validate session and check expiration
	session, err := a.GetSession(cookie.Value)
	if err != nil {
		return nil, err
	}

	// Get the user associated with the session
	user, err := a.GetUserByID(session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ===== VALIDATION HELPERS =====

// EmailExists checks if an email address is already registered
// Used during registration to prevent duplicate accounts
func (a *Auth) EmailExists(email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	err := a.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UsernameExists checks if a username is already taken
// Used during registration to ensure unique usernames
func (a *Auth) UsernameExists(username string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = ?`
	err := a.db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ===== ADMIN FUNCTIONS =====

// UpdateUserRole changes a user's role (admin/moderator/user)
// Only accessible by administrators for user management
// Automatically invalidates all user sessions to force re-login with new permissions
func (a *Auth) UpdateUserRole(userID int, role models.UserRole) error {
	query := `UPDATE users SET role = ? WHERE id = ?`
	_, err := a.db.Exec(query, role, userID)
	if err != nil {
		return err
	}

	// Invalidate all sessions for this user to force them to log back in
	// This ensures they get the new role immediately without cached permissions
	return a.InvalidateUserSessions(userID)
}

// InvalidateUserSessions deletes all active sessions for a specific user
// Used when user roles change or for security purposes
func (a *Auth) InvalidateUserSessions(userID int) error {
	query := `DELETE FROM sessions WHERE user_id = ?`
	_, err := a.db.Exec(query, userID)
	return err
}

// ===== PERMISSION SYSTEM =====

// HasPermission checks if a user has a specific permission based on their role
// Implements role-based access control with hierarchical permissions
// Permissions: "view", "create", "edit_own", "moderate", "delete_posts"
func (a *Auth) HasPermission(user *models.User, permission string) bool {
	// Anonymous users can only view content
	if user == nil {
		return permission == "view"
	}

	// Role-based permission checking
	switch user.Role {
	case models.RoleAdmin:
		return true // Admins have all permissions
	case models.RoleModerator:
		return permission == "view" || permission == "moderate" || permission == "delete_posts"
	case models.RoleUser:
		return permission == "view" || permission == "create" || permission == "edit_own"
	default:
		return permission == "view" // Default to view-only for unknown roles
	}
}

// IsOwner checks if a user owns a specific resource (post, comment, etc.)
// Used to determine if a user can edit or delete their own content
func (a *Auth) IsOwner(user *models.User, ownerID int) bool {
	return user != nil && user.ID == ownerID
}
