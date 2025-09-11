package topiccommands

import (
	"context"
	"errors"
	"testing"

	"github.com/arnald/forum/internal/domain/user"
	testhelpers "github.com/arnald/forum/internal/pkg/testing"
)

func TestUpdateTopicHandler_Handle(t *testing.T) {
	t.Run("group: update topic", func(t *testing.T) {
		testCases := newUpdateTopicTestCases()
		for _, tt := range testCases {
			t.Run(tt.name, runUpdateTopicTest(tt))
		}
	})
}

type updateTopicTestCase struct {
	name       string
	request    UpdateTopicRequest
	setupMocks func(*testhelpers.MockRepository)
	wantTopic  *user.Topic
	wantError  error
}

func newUpdateTopicTestCases() []updateTopicTestCase {
	return []updateTopicTestCase{
		{
			name: "valid request",
			request: UpdateTopicRequest{
				User: &user.User{
					ID:       "test-user-id",
					Username: "testuser",
					Email:    "testuser@example.com",
					Password: "testpassword",
				},
				TopicID:    1,
				CategoryID: 10,
				Title:      "Updated Title",
				Content:    "Updated Content",
				ImagePath:  "",
			},
			setupMocks: func(repo *testhelpers.MockRepository) {
				repo.UpdateTopicFunc = func(ctx context.Context, topic *user.Topic) error {
					return nil
				}
			},
			wantTopic: &user.Topic{
				UserID:     "test-user-id",
				CategoryID: 10,
				ID:         1,
				Title:      "Updated Title",
				Content:    "Updated Content",
				ImagePath:  "",
			},
			wantError: nil,
		},
		{
			name: "invalid request",
			request: UpdateTopicRequest{
				User: &user.User{
					ID:       "test-user-id",
					Username: "testuser",
					Email:    "testuser@example.com",
					Password: "testpassword",
				},
				Title:     "Test Title",
				Content:   "Test Content",
				ImagePath: "",
			},
			setupMocks: func(repo *testhelpers.MockRepository) {
				repo.CreateTopicFunc = func(ctx context.Context, topic *user.Topic) error {
					return testhelpers.ErrTest
				}
			},
			wantTopic: nil,
			wantError: testhelpers.ErrTest,
		},
	}
}

func runUpdateTopicTest(tt updateTopicTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		repo := &testhelpers.MockRepository{}
		tt.setupMocks(repo)

		handler := NewUpdateTopicHandler(repo)
		got, err := handler.Handle(context.Background(), tt.request)

		if !errors.Is(err, tt.wantError) {
			t.Errorf("expected error %v, got %v", tt.wantError, err)
		}

		testhelpers.AssertTopicMatch(t, got, tt.wantTopic)
	}
}

func TestNewUpdateTopicHandler(t *testing.T) {
	repo := &testhelpers.MockRepository{}
	handler := NewUpdateTopicHandler(repo)

	if handler == nil {
		t.Fatal("expected non-nil handler")
	}
}
