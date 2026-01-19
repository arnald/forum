package votes

import "errors"

var (
	ErrVoteNotFound      = errors.New("vote not found")
	ErrInvalidVoteTarget = errors.New("invalid vote target: exactly one of topicID or commentID must be provided")
)
