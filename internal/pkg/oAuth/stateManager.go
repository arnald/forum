package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

var (
	ErrStateNotFound = errors.New("state not found")
	ErrStateExpired  = errors.New("state expired")
)

type StateManager struct {
	states map[string]int64
	mu     sync.RWMutex
	ttl    time.Duration
}

func NewStateManager(ttl time.Duration) *StateManager {
	sm := &StateManager{
		states: make(map[string]int64),
		ttl:    ttl,
	}

	go sm.cleanup()

	return sm
}

func (sm *StateManager) Generate() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(b)

	sm.mu.Lock()
	sm.states[state] = time.Now().Unix()
	sm.mu.Unlock()

	return state, nil
}

func (sm *StateManager) Verify(state string) error {
	sm.mu.Lock()
	createdAt, exists := sm.states[state]
	sm.mu.RUnlock()

	if !exists {
		return ErrStateNotFound
	}

	if time.Now().Unix()-createdAt > int64(sm.ttl.Seconds()) {
		sm.mu.Lock()
		delete(sm.states, state)
		sm.mu.Unlock()
		return ErrStateExpired
	}

	sm.mu.Lock()
	delete(sm.states, state)
	sm.mu.Unlock()

	return nil
}

func (sm *StateManager) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now().Unix()
		sm.mu.Lock()
		for state, createdAt := range sm.states {
			if now-createdAt > int64(sm.ttl.Seconds()) {
				delete(sm.states, state)
			}
		}
		sm.mu.Unlock()
	}
}
