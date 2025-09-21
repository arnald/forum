/*
Package vote contains the domain models for the voting system.

This package defines:
- Vote entity for tracking user votes on posts and comments
- VoteType enumeration for like/dislike votes
- TargetType enumeration for posts vs comments
- Vote counting and status structures
- Business logic for vote management

The voting system allows users to:
- Like or dislike posts and comments
- Change their vote (like to dislike or vice versa)
- Remove their vote entirely
- View vote counts and scores
- See their own voting status on content
*/
package vote

import (
	"time"
)

// VoteType represents the type of vote a user can cast
type VoteType int

// Vote type constants - users can either like or dislike content
const (
	VoteTypeLike    VoteType = iota + 1 // Positive vote (thumbs up)
	VoteTypeDislike                     // Negative vote (thumbs down)
)

// String returns the string representation of a VoteType
// This is useful for API responses and logging
func (v VoteType) String() string {
	switch v {
	case VoteTypeLike:
		return "like"
	case VoteTypeDislike:
		return "dislike"
	default:
		return "unknown"
	}
}

// TargetType represents what type of content is being voted on
type TargetType int

// Target type constants - votes can be on posts or comments
const (
	TargetTypePost    TargetType = iota + 1 // Vote on a forum post
	TargetTypeComment                       // Vote on a comment
)

// String returns the string representation of a TargetType
// This is useful for API responses and database queries
func (t TargetType) String() string {
	switch t {
	case TargetTypePost:
		return "post"
	case TargetTypeComment:
		return "comment"
	default:
		return "unknown"
	}
}

// Vote represents a single vote cast by a user
// Each vote is unique per (user, target) combination
type Vote struct {
	ID         string     // Unique identifier for this vote
	UserID     string     // ID of the user who cast this vote
	TargetID   string     // ID of the post or comment being voted on
	TargetType TargetType // Whether this is a vote on a post or comment
	VoteType   VoteType   // Whether this is a like or dislike
	CreatedAt  time.Time  // When the vote was first cast
	UpdatedAt  time.Time  // When the vote was last changed
}

// VoteCounts aggregates vote statistics for a piece of content
// This provides summary information about all votes on an item
type VoteCounts struct {
	Likes    int // Total number of like votes
	Dislikes int // Total number of dislike votes
	Total    int // Net score: Likes - Dislikes
}

// VoteStatus provides complete voting information for a user and content item
// This includes both the user's personal vote and overall vote statistics
type VoteStatus struct {
	HasVoted   bool       // Whether the current user has voted on this content
	VoteType   *VoteType  // Type of vote from current user (nil if no vote)
	VoteCounts VoteCounts // Aggregate vote counts for this content
}

// Business rules:
// - One vote per user per target (enforced by database constraint)
// - Users can change their vote type (like -> dislike or vice versa)
// - Users can remove their vote entirely
// - Vote changes update the UpdatedAt timestamp