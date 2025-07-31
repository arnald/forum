package uuid_test

import (
	"testing"

	"github.com/google/uuid"

	u "github.com/arnald/forum/internal/pkg/uuid"
)

func TestNewUUIDSuccess(t *testing.T) {
	provider := u.NewProvider()
	id := provider.NewUUID()
	if id == "" {
		t.Error("Expected non-nil UUI, got nil")
	}
}

func TestUUIDUniqueness(t *testing.T) {
	provider := u.NewProvider()
	id1 := provider.NewUUID()
	id2 := provider.NewUUID()

	if id1 == id2 {
		t.Errorf("Expected different UUIDs, got duplicates: %s and %s", id1, id2)
	}
}

func TestUUIDFormat(t *testing.T) {
	provider := u.NewProvider()
	id := provider.NewUUID()
	_, err := uuid.Parse(id)
	if err != nil {
		t.Errorf("UUID %s has invalid format: %v", id, err)
	}
}
