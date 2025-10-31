package notifications

import (
	"context"
	"sync"

	"github.com/arnald/forum/internal/domain/notification"
)

const ChannelCapacity int = 10

type NotificationService struct {
	repo    notification.Repository
	clients map[string][]chan *notification.Notification
	mu      sync.RWMutex
}

func NewNotificationService(repo notification.Repository) *NotificationService {
	return &NotificationService{
		repo:    repo,
		clients: make(map[string][]chan *notification.Notification),
	}
}

func (s *NotificationService) RegisterClient(userID string) chan *notification.Notification {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan *notification.Notification, ChannelCapacity)
	s.clients[userID] = append(s.clients[userID], ch)

	return ch
}

func (s *NotificationService) UnregisterClient(userID string, ch chan *notification.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()

	clients := s.clients[userID]
	for i, client := range clients {
		if client == ch {
			s.clients[userID] = append(clients[:i], clients[i+1:]...)
			close(ch)
			break
		}
	}

	if len(s.clients[userID]) == 0 {
		delete(s.clients, userID)
	}
}

func (s *NotificationService) CreateNotification(ctx context.Context, notification *notification.Notification) error {
	err := s.repo.Create(ctx, notification)
	if err != nil {
		return err
	}

	s.broadcastToUser(notification.UserID, notification)

	return nil
}

func (s *NotificationService) GetNotifications(ctx context.Context, userID string, limit int) ([]*notification.Notification, error) {
	return s.repo.GetByUserID(ctx, userID, limit)
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID int, userID string) error {
	return s.repo.MarkAsRead(ctx, notificationID, userID)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *NotificationService) broadcastToUser(userID string, notification *notification.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()

	clients := s.clients[userID]
	for _, ch := range clients {
		select {
		case ch <- notification:
			// sent
		default:
			// full channel, skip
		}
	}
}
