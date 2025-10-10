package category

import "github.com/arnald/forum/internal/domain/topic"

type Category struct {
	Name        string
	Description string
	CreatedAt   string
	CreatedBy   string
	Topics      []topic.Topic
	ID          int
}
