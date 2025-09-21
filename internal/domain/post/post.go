/*
Package post contains the core domain models for forum posts and categories.

This package defines:
- Post entity with content, metadata, and voting information
- Category entity for post organization
- Post repository interface for data operations
- Business logic for post management

Posts are the main content of the forum where users can:
- Create discussions with titles and content
- Associate posts with categories for organization
- Vote on posts (like/dislike)
- Track engagement through vote counts
*/
package post

import (
	"time"
)

// Post represents a forum post with all its properties and metadata
// This is the core content entity of the forum application
type Post struct {
	ID           string    // Unique identifier (UUID) for the post
	Title        string    // Post title/subject line
	Content      string    // Main post content/body text
	UserID       string    // ID of the user who created this post
	Username     string    // Username of the post author (denormalized for display)
	Categories   []string  // List of category names associated with this post
	CreatedAt    time.Time // Timestamp when post was created
	UpdatedAt    time.Time // Timestamp when post was last modified
	LikeCount    int       // Number of like votes on this post
	DislikeCount int       // Number of dislike votes on this post
	VoteScore    int       // Net score: LikeCount - DislikeCount
}

// Category represents a topic category for organizing posts
// Categories help users find posts on specific topics
type Category struct {
	ID   string // Unique identifier for the category
	Name string // Display name of the category (e.g., "Technology", "Sports")
}

// Note: Posts can be associated with multiple categories (many-to-many relationship)
// VoteScore provides a quick way to rank posts by popularity
// Username is denormalized to avoid joins when displaying posts