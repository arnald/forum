package topiccommands

import (
	"context"
	"errors"
	"testing"

	"github.com/arnald/forum/internal/domain/topic"
	"github.com/arnald/forum/internal/domain/user"
	testhelpers "github.com/arnald/forum/internal/pkg/testing"
)

func TestCreateTopicHandler_Handle(t *testing.T) {
	t.Run("group: create topic", func(t *testing.T) {
		testCases := newCreateTopicTestCases()
		for _, tt := range testCases {
			t.Run(tt.name, runCreateTopicTest(tt))
		}
	})
}

type createTopicTestCase struct {
	name       string
	request    CreateTopicRequest
	setupMocks func(*testhelpers.MockRepository)
	wantTopic  *topic.Topic
	wantError  error
}

func newCreateTopicTestCases() []createTopicTestCase {
	return []createTopicTestCase{
		{
			name: "valid request",
			request: CreateTopicRequest{
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
				repo.CreateTopicFunc = func(ctx context.Context, topic *topic.Topic) error {
					return nil
				}
			},
			wantTopic: &topic.Topic{
				UserID:    "test-user-id",
				Title:     "Test Title",
				Content:   "Test Content",
				ImagePath: "",
			},
			wantError: nil,
		},
		{
			name: "invalid request",
			request: CreateTopicRequest{
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
				repo.CreateTopicFunc = func(ctx context.Context, topic *topic.Topic) error {
					return testhelpers.ErrTest
				}
			},
			wantTopic: nil,
			wantError: testhelpers.ErrTest,
		},
	}
}

func runCreateTopicTest(tt createTopicTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		repo := &testhelpers.MockRepository{}
		tt.setupMocks(repo)

		handler := NewCreateTopicHandler(repo)
		got, err := handler.Handle(context.Background(), tt.request)

		if !errors.Is(err, tt.wantError) {
			t.Errorf("expected error %v, got %v", tt.wantError, err)
		}

		testhelpers.AssertTopicMatch(t, got, tt.wantTopic)
	}
}

func TestNewTopicHandler(t *testing.T) {
	repo := &testhelpers.MockRepository{}
	handler := NewCreateTopicHandler(repo)

	if handler == nil {
		t.Fatal("expected non-nil handler")
	}
}
