package vote

import (
	"encoding/json"
	"net/http"

	"github.com/arnald/forum/internal/app/vote/queries"
	"github.com/arnald/forum/internal/domain/vote"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type VoteHandler struct {
	castVoteHandler     queries.CastVoteRequestHandler
	getVoteStatusHandler queries.GetVoteStatusRequestHandler
	getUserVotesHandler  queries.GetUserVotesRequestHandler
}

func NewVoteHandler(
	castVoteHandler queries.CastVoteRequestHandler,
	getVoteStatusHandler queries.GetVoteStatusRequestHandler,
	getUserVotesHandler queries.GetUserVotesRequestHandler,
) *VoteHandler {
	return &VoteHandler{
		castVoteHandler:      castVoteHandler,
		getVoteStatusHandler: getVoteStatusHandler,
		getUserVotesHandler:  getUserVotesHandler,
	}
}

func (h *VoteHandler) CastVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		TargetID   string `json:"target_id"`
		TargetType string `json:"target_type"` // "post" or "comment"
		VoteType   string `json:"vote_type"`   // "like" or "dislike"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		helpers.RespondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse target type
	var targetType vote.TargetType
	switch req.TargetType {
	case "post":
		targetType = vote.TargetTypePost
	case "comment":
		targetType = vote.TargetTypeComment
	default:
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid target type. Must be 'post' or 'comment'")
		return
	}

	// Parse vote type
	var voteType vote.VoteType
	switch req.VoteType {
	case "like":
		voteType = vote.VoteTypeLike
	case "dislike":
		voteType = vote.VoteTypeDislike
	default:
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid vote type. Must be 'like' or 'dislike'")
		return
	}

	castVoteReq := queries.CastVoteRequest{
		UserID:     userID,
		TargetID:   req.TargetID,
		TargetType: targetType,
		VoteType:   voteType,
	}

	voteStatus, err := h.castVoteHandler.Handle(r.Context(), castVoteReq)
	if err != nil {
		switch err {
		case queries.ErrEmptyUserID:
			helpers.RespondWithError(w, http.StatusUnauthorized, "User authentication required")
		case queries.ErrEmptyTargetID:
			helpers.RespondWithError(w, http.StatusBadRequest, "Target ID is required")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to cast vote")
		}
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, map[string]interface{}{
		"success":     true,
		"vote_status": voteStatus,
	})
}

func (h *VoteHandler) GetVoteStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	targetID := r.URL.Query().Get("target_id")
	targetTypeStr := r.URL.Query().Get("target_type")
	userID := r.Header.Get("X-User-ID")

	if targetID == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Target ID is required")
		return
	}

	// Parse target type
	var targetType vote.TargetType
	switch targetTypeStr {
	case "post":
		targetType = vote.TargetTypePost
	case "comment":
		targetType = vote.TargetTypeComment
	default:
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid target type. Must be 'post' or 'comment'")
		return
	}

	getVoteStatusReq := queries.GetVoteStatusRequest{
		UserID:     userID, // Can be empty for non-authenticated users
		TargetID:   targetID,
		TargetType: targetType,
	}

	voteStatus, err := h.getVoteStatusHandler.Handle(r.Context(), getVoteStatusReq)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to get vote status")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, map[string]interface{}{
		"success":     true,
		"vote_status": voteStatus,
	})
}

func (h *VoteHandler) GetUserVotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		helpers.RespondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	targetTypeStr := r.URL.Query().Get("target_type")

	// Parse target type
	var targetType vote.TargetType
	switch targetTypeStr {
	case "post":
		targetType = vote.TargetTypePost
	case "comment":
		targetType = vote.TargetTypeComment
	default:
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid target type. Must be 'post' or 'comment'")
		return
	}

	getUserVotesReq := queries.GetUserVotesRequest{
		UserID:     userID,
		TargetType: targetType,
	}

	votes, err := h.getUserVotesHandler.Handle(r.Context(), getUserVotesReq)
	if err != nil {
		switch err {
		case queries.ErrEmptyUserID:
			helpers.RespondWithError(w, http.StatusUnauthorized, "User authentication required")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to get user votes")
		}
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, map[string]interface{}{
		"success": true,
		"votes":   votes,
	})
}