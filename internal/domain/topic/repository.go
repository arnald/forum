package topic

import "context"

type Repository interface {
	CreateTopic(ctx context.Context, topic *Topic) error
	UpdateTopic(ctx context.Context, topic *Topic) error
	DeleteTopic(ctx context.Context, userID string, topicID int) error
	GetTopicByID(ctx context.Context, topicID int) (*Topic, error)
	GetAllTopics(ctx context.Context, page, size, categoryID int, orderBy, order, filter string) ([]Topic, error)
	GetTotalTopicsCount(ctx context.Context, filter string, categoryID int) (int, error)
}
