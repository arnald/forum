/*
Package post repository defines the contract for post data operations.

This file contains the Repository interface that defines all data access
operations for post entities and their relationships. The interface follows
the Repository pattern to abstract data storage concerns from business logic.

The Repository interface provides methods for:
- Post CRUD operations (Create, Read, Update, Delete)
- Post filtering and retrieval by various criteria
- Category management and post-category relationships
- User-post association queries
- Vote count integration and denormalization

This abstraction supports different storage implementations while maintaining
consistent business logic for forum post management.
*/
package post

import (
	"context"
)

// Repository defines the contract for post data persistence operations
// This interface abstracts the data storage layer from business logic,
// enabling different database implementations while maintaining consistency
//
// All methods use context.Context for:
// - Request timeout and cancellation support
// - Transaction management for data consistency
// - Trace propagation for monitoring and debugging
//
// Repository implementations should:
// - Maintain referential integrity with users and categories
// - Update denormalized vote counts when votes change
// - Handle category associations efficiently
// - Support pagination for large post lists (future enhancement)
// - Implement proper indexing for query performance
type Repository interface {
	// CreatePost creates a new post in the data store
	// This method handles post creation with category associations
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - post: Post entity to be created with all required fields
	//
	// Returns:
	//   - error: Validation, constraint, or database errors
	//
	// Requirements:
	//   - Post ID should be a valid UUID
	//   - UserID must reference an existing user
	//   - Categories must exist or be created
	//   - Title and Content must not be empty
	//
	// Side effects:
	//   - Creates category associations in junction table
	//   - Updates category post counts (if denormalized)
	CreatePost(ctx context.Context, post *Post) error

	// GetPostByID retrieves a specific post by its unique identifier
	// This method includes vote counts and category information
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - id: UUID of the post to retrieve
	//
	// Returns:
	//   - *Post: Post entity with complete information, nil if not found
	//   - error: Database errors (not post-not-found errors)
	//
	// Includes:
	//   - Post content and metadata
	//   - Associated categories
	//   - Current vote counts (likes, dislikes, score)
	//   - Author username (denormalized)
	GetPostByID(ctx context.Context, id string) (*Post, error)

	// GetAllPosts retrieves all posts from the data store
	// This method returns posts with complete information for listings
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//
	// Returns:
	//   - []*Post: Slice of all posts ordered by creation date (newest first)
	//   - error: Database or connection errors
	//
	// Performance considerations:
	//   - Should implement pagination for large datasets
	//   - Includes denormalized vote counts to avoid N+1 queries
	//   - May limit returned fields for list views vs detail views
	GetAllPosts(ctx context.Context) ([]*Post, error)

	// GetPostsByUserID retrieves all posts created by a specific user
	// This method supports user profile pages and author-specific listings
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - userID: UUID of the user whose posts to retrieve
	//
	// Returns:
	//   - []*Post: Slice of posts by the specified user, ordered by date
	//   - error: Database errors
	//
	// Use cases:
	//   - User profile pages showing their posts
	//   - Author activity tracking
	//   - Content moderation by user
	GetPostsByUserID(ctx context.Context, userID string) ([]*Post, error)

	// GetPostsByCategory retrieves posts associated with a category name
	// This method supports category browsing by name (for URL-friendly routing)
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - category: Name of the category to filter by
	//
	// Returns:
	//   - []*Post: Slice of posts in the specified category
	//   - error: Database errors or category not found
	//
	// Usage:
	//   - Category page browsing (/category/technology)
	//   - Content filtering by topic
	//   - Category-specific RSS feeds
	GetPostsByCategory(ctx context.Context, category string) ([]*Post, error)

	// GetPostsByCategoryID retrieves posts associated with a category ID
	// This method supports efficient category filtering using primary keys
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - categoryID: UUID of the category to filter by
	//
	// Returns:
	//   - []*Post: Slice of posts in the specified category
	//   - error: Database errors
	//
	// Performance:
	//   - More efficient than name-based lookup
	//   - Used internally when category ID is already known
	//   - Supports JOIN operations for complex queries
	GetPostsByCategoryID(ctx context.Context, categoryID string) ([]*Post, error)

	// UpdatePost modifies an existing post in the data store
	// This method handles post content updates and category changes
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - post: Post entity with updated fields
	//
	// Returns:
	//   - error: Validation, authorization, or database errors
	//
	// Business rules:
	//   - Only post author or admins can update posts
	//   - Title and content cannot be empty
	//   - Category associations can be modified
	//   - UpdatedAt timestamp should be refreshed
	//
	// Side effects:
	//   - Updates category associations if changed
	//   - Maintains vote counts during updates
	UpdatePost(ctx context.Context, post *Post) error

	// DeletePost removes a post and its associations from the data store
	// This method handles cascading deletes for related data
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - id: UUID of the post to delete
	//
	// Returns:
	//   - error: Authorization or database errors
	//
	// Cascading effects:
	//   - Removes all votes on the post
	//   - Removes all comments on the post
	//   - Removes category associations
	//   - Updates category post counts
	//
	// Business rules:
	//   - Only post author or admins can delete posts
	//   - Soft delete may be preferred for audit trails
	DeletePost(ctx context.Context, id string) error

	// GetCategories retrieves all available categories
	// This method supports category selection and browsing interfaces
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//
	// Returns:
	//   - []*Category: Slice of all categories with metadata
	//   - error: Database errors
	//
	// Usage:
	//   - Category dropdown in post creation forms
	//   - Category navigation menus
	//   - Site-wide category statistics
	//
	// Note: This method is also available in category repository
	// but included here for convenience in post-related operations
	GetCategories(ctx context.Context) ([]*Category, error)
}