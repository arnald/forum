package helpers

import (
	"context"
	"fmt"
	"os/user"
)

func GetUserFromContext(ctx context.Context, key string) (*user.User, error) {
	user, ok := ctx.Value(key).(*user.User)
	if !ok {
		return nil, fmt.Errorf("no user found in context")
	}
	return user, nil
}
