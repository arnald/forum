package helpers

import (
	"context"
	"errors"
	"os/user"
)

var ErrNoUserFound = errors.New("no user found in context")

func GetUserFromContext(ctx context.Context, key string) (*user.User, error) {
	user, ok := ctx.Value(key).(*user.User)
	if !ok {
		return nil, ErrNoUserFound
	}
	return user, nil
}
