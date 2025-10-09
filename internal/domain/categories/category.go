package categories

import "github.com/arnald/forum/internal/domain/topic"

type Category struct {
	ID          int
	Name        string
	Description string
	CreatedAt   string
	Topics      []topic.Topic
	CreatedBy   string
}
