package categories

import "github.com/arnald/forum/internal/domain/user"

type Category struct {
	ID          int
	Name        string
	Description string
	CreatedAt   string
	Topics      []user.Topic
	CreatedBy   string
}
