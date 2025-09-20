package vote

import (
	"time"
)

type VoteType int

const (
	VoteTypeLike VoteType = iota + 1
	VoteTypeDislike
)

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

type TargetType int

const (
	TargetTypePost TargetType = iota + 1
	TargetTypeComment
)

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

type Vote struct {
	ID         string
	UserID     string
	TargetID   string
	TargetType TargetType
	VoteType   VoteType
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type VoteCounts struct {
	Likes    int
	Dislikes int
	Total    int // Likes - Dislikes
}

type VoteStatus struct {
	HasVoted  bool
	VoteType  *VoteType // nil if user hasn't voted
	VoteCounts VoteCounts
}