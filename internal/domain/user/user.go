/*
Package user contains the core domain models and interfaces for user management.

This package defines:
- User entity with all user properties
- User repository interface for data operations
- Session management interfaces
- Authentication and authorization logic

The User domain is central to the forum application as it handles:
- User registration and authentication
- Session management and security
- User profile information
- Authorization for protected operations
*/
package user

import (
	"time"
)

// User represents a forum user with all their properties and metadata
// This is the core domain entity for user management
type User struct {
	ID        string     // Unique identifier (UUID) for the user
	Email     string     // User's email address (must be unique, used for login)
	Username  string     // User's display name (must be unique, used for login)
	Password  string     // Encrypted password hash (never store plain text)
	Role      string     // User role for authorization (e.g., "user", "admin")
	AvatarURL *string    // Optional profile picture URL (pointer allows nil)
	CreatedAt time.Time  // Timestamp when user account was created
}

// Note: Password field contains bcrypt hash, never plain text password
// Email and Username fields must be unique across all users
// ID field is generated using UUID to ensure global uniqueness
