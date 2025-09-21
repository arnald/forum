/*
Package user repository defines the contract for user data operations.

This file contains the Repository interface that defines all data access
operations for user entities. The interface follows the Repository pattern
to abstract data storage concerns from business logic.

The Repository interface provides methods for:
- User registration and account creation
- User retrieval by username or email for authentication
- Administrative user listing operations
- Database transaction support through context

This abstraction allows different storage implementations (SQLite, PostgreSQL, etc.)
while maintaining consistent business logic in the application layer.
*/
package user

import (
	"context"
)

// Repository defines the contract for user data persistence operations
// This interface abstracts the data storage layer from business logic,
// allowing different database implementations while maintaining consistency
//
// All methods use context.Context for:
// - Request timeout and cancellation support
// - Transaction management in database implementations
// - Trace propagation for monitoring and debugging
//
// Repository implementations should:
// - Enforce unique constraints on username and email fields
// - Handle database errors and convert them to domain errors
// - Support atomic operations for data consistency
// - Implement proper connection pooling and resource management
type Repository interface {
	// GetAll retrieves all users from the data store
	// This method is typically used for administrative purposes and user listings
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//
	// Returns:
	//   - []User: Slice of all user entities in the system
	//   - error: Database or connection errors
	//
	// Use cases:
	//   - Admin user management interfaces
	//   - User statistics and reporting
	//   - System maintenance operations
	GetAll(ctx context.Context) ([]User, error)

	// UserRegister creates a new user account in the data store
	// This method handles new user registration with validation and constraints
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - user: User entity to be created (should have encrypted password)
	//
	// Returns:
	//   - error: Validation, constraint, or database errors
	//
	// Requirements:
	//   - Username must be unique across all users
	//   - Email must be unique across all users
	//   - Password should already be encrypted (bcrypt)
	//   - User ID should be a valid UUID
	//
	// Common errors:
	//   - Duplicate username or email violations
	//   - Database connection issues
	//   - Invalid data format or constraints
	UserRegister(ctx context.Context, user *User) error

	// GetUserByUsername retrieves a user by their username for authentication
	// This method is used during login processes and user lookups
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - username: Exact username to search for (case-sensitive)
	//
	// Returns:
	//   - *User: User entity if found, nil if not found
	//   - error: Database errors (not user-not-found errors)
	//
	// Usage:
	//   - Username-based authentication
	//   - User profile lookups
	//   - Permission and authorization checks
	//
	// Note: Returns nil user and nil error when user is not found
	GetUserByUsername(ctx context.Context, username string) (*User, error)

	// GetUserByEmail retrieves a user by their email address for authentication
	// This method supports email-based login and user recovery operations
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - email: Exact email address to search for (case-insensitive)
	//
	// Returns:
	//   - *User: User entity if found, nil if not found
	//   - error: Database errors (not user-not-found errors)
	//
	// Usage:
	//   - Email-based authentication
	//   - Password reset workflows
	//   - Account recovery processes
	//   - Email verification operations
	//
	// Note: Returns nil user and nil error when user is not found
	// Email matching should be case-insensitive for better user experience
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
