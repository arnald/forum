/*
Package queries contains the application layer use cases for vote operations.

This package implements the Command Query Responsibility Segregation (CQRS) pattern
for vote-related business operations. It contains:

- CastVote: Handle user voting on posts and comments
- GetVoteStatus: Retrieve voting information for content
- GetUserVotes: Get all votes cast by a specific user

The vote system supports:
- Like/dislike voting on posts and comments
- Vote toggling (click same vote to remove it)
- Vote changing (like to dislike or vice versa)
- Vote counting and status tracking
*/
package queries

import (
	"context"
	"time"

	"github.com/arnald/forum/internal/domain/vote"
	"github.com/arnald/forum/internal/pkg/uuid"
)

// CastVoteRequest contains the data needed to cast a vote
// This represents a user's intent to vote on a piece of content
type CastVoteRequest struct {
	UserID     string           // ID of the user casting the vote
	TargetID   string           // ID of the post or comment being voted on
	TargetType vote.TargetType  // Whether voting on a post or comment
	VoteType   vote.VoteType    // Whether this is a like or dislike
}

// CastVoteRequestHandler defines the interface for vote casting use case
// This follows the Command pattern for handling business operations
type CastVoteRequestHandler interface {
	Handle(ctx context.Context, req CastVoteRequest) (*vote.VoteStatus, error)
}

// castVoteRequestHandler implements the vote casting business logic
// It encapsulates the rules around voting behavior and persistence
type castVoteRequestHandler struct {
	repo         vote.Repository // Repository for vote data operations
	uuidProvider uuid.Provider   // UUID generator for new vote IDs
}

// NewCastVoteHandler creates a new instance of the vote casting handler
// It injects the required dependencies for vote operations
//
// Parameters:
//   - repo: Repository interface for vote data persistence
//   - uuidProvider: UUID generator for creating unique vote IDs
//
// Returns:
//   - CastVoteRequestHandler: Handler ready to process vote casting requests
func NewCastVoteHandler(repo vote.Repository, uuidProvider uuid.Provider) CastVoteRequestHandler {
	return &castVoteRequestHandler{
		repo:         repo,         // Store vote repository
		uuidProvider: uuidProvider, // Store UUID generator
	}
}

// Handle processes a vote casting request with full business logic
// This method implements the core voting behavior including vote toggling and changing
//
// Business Logic:
// 1. If user hasn't voted: Create new vote
// 2. If user clicks same vote type: Remove vote (toggle off)
// 3. If user clicks different vote type: Change vote (like<->dislike)
//
// Parameters:
//   - ctx: Context for request cancellation and timeouts
//   - req: Vote request containing user, target, and vote type
//
// Returns:
//   - *vote.VoteStatus: Updated voting status after the operation
//   - error: Any error that occurred during processing
//
// Errors:
//   - ErrEmptyUserID: When user ID is not provided
//   - ErrEmptyTargetID: When target ID is not provided
//   - Repository errors: Database or validation errors
func (h *castVoteRequestHandler) Handle(ctx context.Context, req CastVoteRequest) (*vote.VoteStatus, error) {
	// Validate required fields
	if req.UserID == "" {
		return nil, ErrEmptyUserID
	}
	if req.TargetID == "" {
		return nil, ErrEmptyTargetID
	}

	// Check if user has already voted on this target
	existingVote, err := h.repo.GetVote(ctx, req.UserID, req.TargetID, req.TargetType)

	if err != nil && existingVote == nil {
		// Case 1: No existing vote - create a new vote
		newVote := &vote.Vote{
			ID:         h.uuidProvider.NewUUID(), // Generate unique ID
			UserID:     req.UserID,               // User casting the vote
			TargetID:   req.TargetID,             // Post or comment being voted on
			TargetType: req.TargetType,           // Post vs comment
			VoteType:   req.VoteType,             // Like vs dislike
			CreatedAt:  time.Now(),               // Timestamp for creation
			UpdatedAt:  time.Now(),               // Timestamp for last update
		}

		// Persist the new vote to database
		err = h.repo.CreateVote(ctx, newVote)
		if err != nil {
			return nil, err
		}
	} else if existingVote != nil {
		// Case 2: User has already voted - handle vote change or removal
		if existingVote.VoteType == req.VoteType {
			// Same vote type clicked - toggle off (remove vote)
			// This allows users to "unlike" or "undislike" content
			err = h.repo.DeleteVote(ctx, req.UserID, req.TargetID, req.TargetType)
			if err != nil {
				return nil, err
			}
		} else {
			// Different vote type clicked - change vote (like<->dislike)
			// This allows users to change from like to dislike or vice versa
			existingVote.VoteType = req.VoteType
			existingVote.UpdatedAt = time.Now() // Update timestamp
			err = h.repo.UpdateVote(ctx, existingVote)
			if err != nil {
				return nil, err
			}
		}
	}

	// Return the updated vote status including counts and user's current vote
	// This provides all information needed to update the UI
	return h.repo.GetVoteStatus(ctx, req.UserID, req.TargetID, req.TargetType)
}