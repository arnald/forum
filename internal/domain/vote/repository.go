/*
Package vote repository defines the contract for vote data operations.

This file contains the Repository interface that defines all data access
operations for vote entities and voting statistics. The interface follows
the Repository pattern to abstract data storage concerns from business logic.

The Repository interface provides methods for:
- Vote CRUD operations (Create, Read, Update, Delete)
- Vote status checking and user voting history
- Vote count aggregation for posts and comments
- Voting business logic support (toggle, change, remove)
- Denormalized vote count maintenance

This abstraction enables different storage implementations while maintaining
consistent business logic for the forum voting system.
*/
package vote

import (
	"context"
)

// Repository defines the contract for vote data persistence operations
// This interface abstracts the data storage layer from business logic,
// enabling different database implementations while maintaining consistency
//
// All methods use context.Context for:
// - Request timeout and cancellation support
// - Transaction management for atomic vote operations
// - Trace propagation for monitoring and debugging
//
// Repository implementations should:
// - Enforce unique constraints (one vote per user per target)
// - Update denormalized vote counts in posts/comments
// - Handle vote changes atomically (create/update/delete)
// - Support efficient aggregation queries for vote counts
// - Maintain referential integrity with users and content
type Repository interface {
	// CreateVote creates a new vote in the data store
	// This method handles initial vote casting on posts or comments
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - vote: Vote entity to be created with all required fields
	//
	// Returns:
	//   - error: Constraint violations, validation, or database errors
	//
	// Requirements:
	//   - Vote ID should be a valid UUID
	//   - UserID must reference an existing user
	//   - TargetID must reference an existing post or comment
	//   - TargetType must be valid (Post or Comment)
	//   - VoteType must be valid (Like or Dislike)
	//
	// Business rules:
	//   - Only one vote per user per target (enforced by unique constraint)
	//   - Vote creation should update denormalized counts in target entity
	//
	// Side effects:
	//   - Updates LikeCount/DislikeCount/VoteScore in posts/comments table
	//   - May trigger notifications for content authors
	CreateVote(ctx context.Context, vote *Vote) error

	// UpdateVote modifies an existing vote in the data store
	// This method handles vote type changes (like to dislike or vice versa)
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - vote: Vote entity with updated VoteType and UpdatedAt fields
	//
	// Returns:
	//   - error: Validation or database errors
	//
	// Business rules:
	//   - Only VoteType and UpdatedAt can be changed
	//   - User cannot change vote on someone else's vote
	//   - Must update denormalized counts in target entity
	//
	// Side effects:
	//   - Updates LikeCount/DislikeCount/VoteScore in posts/comments table
	//   - Atomic operation to prevent count inconsistencies
	UpdateVote(ctx context.Context, vote *Vote) error

	// DeleteVote removes a vote from the data store
	// This method handles vote removal (user "unlikes" or "undislikes")
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - userID: ID of the user whose vote to delete
	//   - targetID: ID of the post or comment being voted on
	//   - targetType: Whether the target is a post or comment
	//
	// Returns:
	//   - error: Authorization or database errors
	//
	// Business rules:
	//   - User can only delete their own votes
	//   - Must update denormalized counts in target entity
	//   - Should be idempotent (no error if vote doesn't exist)
	//
	// Side effects:
	//   - Updates LikeCount/DislikeCount/VoteScore in posts/comments table
	//   - Removes vote record from votes table
	DeleteVote(ctx context.Context, userID, targetID string, targetType TargetType) error

	// GetVote retrieves a specific vote by user and target
	// This method checks if a user has voted on specific content
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - userID: ID of the user whose vote to retrieve
	//   - targetID: ID of the post or comment to check
	//   - targetType: Whether the target is a post or comment
	//
	// Returns:
	//   - *Vote: Vote entity if found, nil if user hasn't voted
	//   - error: Database errors (not vote-not-found errors)
	//
	// Usage:
	//   - Checking user's current vote before allowing vote changes
	//   - Displaying user's vote status in UI (highlighted buttons)
	//   - Implementing vote toggle logic (same vote = remove, different = change)
	//
	// Note: Returns nil vote and nil error when no vote exists
	GetVote(ctx context.Context, userID, targetID string, targetType TargetType) (*Vote, error)

	// GetVoteCounts retrieves aggregated vote statistics for content
	// This method provides vote summaries without user-specific information
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - targetID: ID of the post or comment to get counts for
	//   - targetType: Whether the target is a post or comment
	//
	// Returns:
	//   - *VoteCounts: Aggregated vote statistics (likes, dislikes, total)
	//   - error: Database errors
	//
	// Performance:
	//   - Should use denormalized counts from posts/comments table when available
	//   - Fallback to real-time aggregation if denormalization is not implemented
	//   - Critical for displaying vote counts without authentication
	//
	// Usage:
	//   - Public vote displays on posts and comments
	//   - Sorting content by popularity (vote score)
	//   - Analytics and reporting on content engagement
	GetVoteCounts(ctx context.Context, targetID string, targetType TargetType) (*VoteCounts, error)

	// GetVoteStatus retrieves complete voting information for a user and content
	// This method combines user's personal vote with overall vote statistics
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - userID: ID of the user to get vote status for (can be empty for anonymous)
	//   - targetID: ID of the post or comment to check
	//   - targetType: Whether the target is a post or comment
	//
	// Returns:
	//   - *VoteStatus: Complete voting information (user vote + counts)
	//   - error: Database errors
	//
	// Returns:
	//   - HasVoted: Whether the user has voted on this content
	//   - VoteType: Type of vote (like/dislike) if user has voted
	//   - VoteCounts: Aggregated vote statistics for all users
	//
	// Usage:
	//   - Displaying complete vote information in UI
	//   - Determining which vote buttons to highlight
	//   - Single query to get all voting information needed
	//
	// Performance optimization:
	//   - Combines multiple queries into a single database call
	//   - Essential for efficient vote display in content lists
	GetVoteStatus(ctx context.Context, userID, targetID string, targetType TargetType) (*VoteStatus, error)

	// GetUserVotes retrieves all votes cast by a specific user
	// This method supports user voting history and activity tracking
	//
	// Parameters:
	//   - ctx: Context for timeout/cancellation and transaction management
	//   - userID: ID of the user whose votes to retrieve
	//   - targetType: Filter by target type (post or comment), or all if not specified
	//
	// Returns:
	//   - []*Vote: Slice of votes cast by the user, ordered by date
	//   - error: Database errors
	//
	// Use cases:
	//   - User profile pages showing voting activity
	//   - Content recommendation based on user preferences
	//   - Moderation tools for reviewing user behavior
	//   - Analytics on user engagement patterns
	//
	// Privacy considerations:
	//   - Should only be accessible to the user themselves or admins
	//   - May need to respect user privacy settings
	//   - Consider anonymizing votes for public statistics
	GetUserVotes(ctx context.Context, userID string, targetType TargetType) ([]*Vote, error)
}