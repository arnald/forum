/*
Package queries contains user-related application layer use cases.

This package implements the Command Query Responsibility Segregation (CQRS) pattern
for user operations. It contains:

- UserRegister: Handle new user account creation with validation and encryption
- UserLoginEmail: Handle user authentication using email credentials
- UserLoginUsername: Handle user authentication using username credentials

The user query handlers implement business logic for user management including:
- Account creation with email validation
- Password encryption using bcrypt
- User authentication and credential verification
- Unique constraint enforcement for usernames and emails
*/
package queries

import (
	"context"
	"time"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/bcrypt"
	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/uuid"
)

// UserRegisterRequest contains the data needed to create a new user account
// This represents the information collected from a registration form
type UserRegisterRequest struct {
	Name     string // Username for the new account (must be unique)
	Password string // Plain text password (will be encrypted before storage)
	Email    string // Email address for the account (must be unique and valid)
}

// UserRegisterRequestHandler defines the interface for user registration use case
// This follows the Command pattern for handling user account creation
type UserRegisterRequestHandler interface {
	Handle(ctx context.Context, req UserRegisterRequest) (*user.User, error)
}

// userRegisterRequestHandler implements the user registration business logic
// It encapsulates the rules around user creation, validation, and security
type userRegisterRequestHandler struct {
	uuidiProvider      uuid.Provider   // UUID generator for unique user IDs
	encryptionProvider bcrypt.Provider // Password encryption provider for security
	repo               user.Repository // Repository for user data persistence
}

// NewUserRegisterHandler creates a new instance of the user registration handler
// It injects the required dependencies for user registration operations
//
// Parameters:
//   - repo: Repository interface for user data persistence
//   - uuidProvider: UUID generator for creating unique user IDs
//   - en: Encryption provider for password hashing
//
// Returns:
//   - UserRegisterRequestHandler: Handler ready to process registration requests
func NewUserRegisterHandler(repo user.Repository, uuidProvider uuid.Provider, en bcrypt.Provider) UserRegisterRequestHandler {
	return userRegisterRequestHandler{
		repo:               repo,           // Store user repository
		uuidiProvider:      uuidProvider,   // Store UUID generator
		encryptionProvider: en,             // Store encryption provider
	}
}

// Handle processes a user registration request with full validation and security
// This method implements the core user registration business logic
//
// Registration Process:
// 1. Create user entity with generated UUID and timestamp
// 2. Validate email format and constraints
// 3. Encrypt password using bcrypt for security
// 4. Persist user to database through repository
// 5. Return created user (with encrypted password)
//
// Parameters:
//   - ctx: Context for request cancellation and timeouts
//   - req: Registration request containing user details
//
// Returns:
//   - *user.User: Created user entity with encrypted password
//   - error: Any validation, encryption, or persistence error
//
// Security Features:
//   - Email validation prevents invalid addresses
//   - Password encryption using bcrypt (never stores plain text)
//   - UUID generation ensures unique user identifiers
//   - Repository enforces unique constraints for username/email
func (h userRegisterRequestHandler) Handle(ctx context.Context, req UserRegisterRequest) (*user.User, error) {
	// Create new user entity with provided data and generated fields
	user := &user.User{
		CreatedAt: time.Now(),                // Set creation timestamp
		Password:  req.Password,              // Temporary plain text (will be encrypted)
		AvatarURL: nil,                       // No avatar on registration
		Username:  req.Name,                  // Username from request
		Email:     req.Email,                 // Email from request
		ID:        h.uuidiProvider.NewUUID(), // Generate unique ID
	}

	// Validate email format and constraints
	// This prevents invalid email addresses from being registered
	err := helpers.ValidateEmail(user.Email)
	if err != nil {
		return nil, err // Return validation error
	}

	// Encrypt password using bcrypt for security
	// Never store plain text passwords in the database
	encryptedPass, err := h.encryptionProvider.Generate(user.Password)
	if err != nil {
		return nil, err // Return encryption error
	}

	// Replace plain text password with encrypted version
	user.Password = encryptedPass

	// Persist user to database through repository layer
	// Repository will enforce unique constraints for username/email
	err = h.repo.UserRegister(ctx, user)
	if err != nil {
		return nil, err // Return persistence error
	}

	// Return successfully created user
	return user, err
}
