/*
Package bcrypt provides password encryption and verification functionality.

This package wraps the golang.org/x/crypto/bcrypt library to provide
secure password handling for the forum application. It implements:

- Password hashing using bcrypt algorithm
- Secure cost factor for computational difficulty
- Password verification against stored hashes
- Provider pattern for testability and dependency injection
- Protection against timing attacks

Bcrypt is specifically designed for password hashing and includes
built-in salt generation and timing attack protection.
*/
package bcrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// encryptionCost defines the computational cost for bcrypt hashing
// Cost of 12 provides a good balance between security and performance
// This results in approximately 250ms hashing time on modern hardware
const encryptionCost = 12

// Provider defines the interface for password encryption services
// This interface allows for easy testing and different encryption implementations
// while maintaining a consistent API across the application
type Provider interface {
	// Generate creates a bcrypt hash from a plaintext password
	// Parameters:
	//   - plaintextPassword: The plain text password to encrypt
	// Returns:
	//   - string: Bcrypt hash suitable for database storage
	//   - error: Any encryption errors that occurred
	Generate(plaintextPassword string) (string, error)

	// Matches verifies a password against a stored bcrypt hash
	// Parameters:
	//   - databasePassword: The stored bcrypt hash from database
	//   - passwordFromRequest: The plaintext password from user input
	// Returns:
	//   - error: nil if passwords match, error if they don't match or other issues
	Matches(databasePassword string, passwordFromRequest string) error
}

// NewProvider creates a new bcrypt encryption provider instance
// This factory function returns a concrete implementation of the Provider interface
//
// Returns:
//   - Provider: Encryption provider ready to hash and verify passwords
//
// Usage:
//   encryptor := bcrypt.NewProvider()
//   hash, err := encryptor.Generate("userPassword123")
//   err = encryptor.Matches(hash, "userPassword123")
//
// The provider is thread-safe and can be used concurrently
func NewProvider() Provider {
	return &encryptionProvider{} // Return concrete implementation
}

// encryptionProvider is the concrete implementation of the Provider interface
// It wraps the golang.org/x/crypto/bcrypt library for password operations
// This struct contains no fields as bcrypt operations are stateless
type encryptionProvider struct{}

// Generate creates a bcrypt hash from a plaintext password
// This method uses bcrypt with a configured cost factor to create
// a cryptographically secure hash suitable for database storage
//
// Parameters:
//   - plaintextPassword: The user's plaintext password to encrypt
//
// Returns:
//   - string: Bcrypt hash in the format "$2a$12$hash"
//   - error: Any errors during hash generation
//
// Security Features:
//   - Automatic salt generation (unique per password)
//   - Configurable cost factor for computational difficulty
//   - Designed to be slow to prevent brute force attacks
//   - Constant time verification to prevent timing attacks
//
// Example:
//   hash, err := provider.Generate("mySecurePassword123")
//   // hash might be: "$2a$12$exampleHashString..."
func (p *encryptionProvider) Generate(plaintextPassword string) (string, error) {
	// Generate bcrypt hash with configured cost factor
	// The cost factor determines how many rounds of hashing are performed
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), encryptionCost)
	if err != nil {
		return "", err // Return error if hashing fails
	}

	return string(hash), nil // Return hash as string for database storage
}

// Matches verifies a plaintext password against a stored bcrypt hash
// This method uses constant-time comparison to prevent timing attacks
// and properly handles bcrypt-specific errors
//
// Parameters:
//   - databasePassword: The bcrypt hash stored in the database
//   - passwordFromRequest: The plaintext password provided by the user
//
// Returns:
//   - error: nil if passwords match, specific error if they don't match
//
// Error Handling:
//   - ErrMismatchedHashAndPassword: Passwords don't match (invalid login)
//   - Other errors: Malformed hash or other bcrypt issues
//
// Security Features:
//   - Constant time comparison prevents timing attacks
//   - Handles salt extraction and verification automatically
//   - Works with hashes generated at different cost factors
//
// Usage:
//   err := provider.Matches(storedHash, userInputPassword)
//   if err != nil {
//       // Handle authentication failure
//   }
func (p *encryptionProvider) Matches(databasePassword string, passwordFromRequest string) error {
	// Compare hash and password using bcrypt's secure comparison
	// This function automatically extracts salt and cost from the hash
	err := bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(passwordFromRequest))
	if err != nil {
		// Handle bcrypt-specific errors
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return err // Password doesn't match (authentication failure)
		default:
			return err // Other bcrypt errors (malformed hash, etc.)
		}
	}

	return nil // Passwords match - authentication successful
}
