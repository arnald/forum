/*
Package uuid provides UUID generation functionality for the forum application.

This package wraps the Google UUID library to provide a clean interface
for generating unique identifiers throughout the application. It uses:

- RFC 4122 compliant UUID version 4 (random)
- String representation for easy storage and transmission
- Provider pattern for testability and dependency injection
- Thread-safe UUID generation

UUIDs are used as primary keys for all entities in the forum
including users, posts, comments, votes, and categories.
*/
package uuid

import (
	"github.com/google/uuid"
)

// Provider defines the interface for UUID generation services
// This interface allows for easy testing and different UUID implementations
// while maintaining a consistent API across the application
type Provider interface {
	// NewUUID generates a new random UUID string
	// Returns a RFC 4122 compliant UUID v4 in string format
	// Example: "550e8400-e29b-41d4-a716-446655440000"
	NewUUID() string
}

// NewProvider creates a new UUID provider instance
// This factory function returns a concrete implementation of the Provider interface
//
// Returns:
//   - Provider: UUID provider ready to generate unique identifiers
//
// Usage:
//   uuidProvider := uuid.NewProvider()
//   userID := uuidProvider.NewUUID()
//
// The provider is thread-safe and can be used concurrently
// across multiple goroutines without synchronization
func NewProvider() Provider {
	return uuidProvider{} // Return concrete implementation
}

// uuidProvider is the concrete implementation of the Provider interface
// It wraps the Google UUID library to provide UUID generation functionality
// This struct contains no fields as UUID generation is stateless
type uuidProvider struct{}

// NewUUID generates a new random UUID and returns it as a string
// This method creates RFC 4122 compliant version 4 UUIDs using
// cryptographically secure random numbers
//
// Returns:
//   - string: UUID in canonical string format (36 characters with hyphens)
//
// Technical Details:
//   - Uses crypto/rand for random number generation
//   - Version 4 UUID with 122 random bits
//   - Formatted as 8-4-4-4-12 hexadecimal digits
//   - Thread-safe and suitable for concurrent use
//
// Example Output: "f47ac10b-58cc-4372-a567-0e02b2c3d479"
//
// Usage:
//   userID := provider.NewUUID()  // For user entity primary key
//   postID := provider.NewUUID()  // For post entity primary key
//   voteID := provider.NewUUID()  // For vote entity primary key
func (u uuidProvider) NewUUID() string {
	return uuid.New().String() // Generate and format UUID as string
}
