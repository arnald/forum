package uuid_test

import (
	"testing"

	u "github.com/arnald/forum/internal/pkg/uuid"

	"github.com/google/uuid"
)

func TestNewUUIDSuccess(t *testing.T) {
	provider := u.NewProvider()
	id := provider.NewUUID()
	if id.String() == "" {
		t.Error("Expected non-nil UUI, got nil")
	}
}

func TestUUIDUniqueness(t *testing.T) {
	provider := u.NewProvider()
	id1 := provider.NewUUID()
	id2 := provider.NewUUID()

	if id1.String() == id2.String() {
		t.Errorf("Expected different UUIDs, got duplicates: %s and %s", id1, id2)
	}
}

func TestUUIDFormat(t *testing.T) {
	provider := u.NewProvider()
	id := provider.NewUUID()
	_, err := uuid.Parse(id.String())
	if err != nil {
		t.Errorf("UUID %s has invalid format: %v", id, err)
	}
}
