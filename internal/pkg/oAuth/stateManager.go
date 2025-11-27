package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

const (
	cleanUpInterval = 3
	bufferSize      = 32
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
	b := make([]byte, bufferSize)
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
	defer sm.mu.Unlock()

	createdAt, exists := sm.states[state]
	if !exists {
		return ErrStateNotFound
	}

	if time.Now().Unix()-createdAt > int64(sm.ttl.Seconds()) {
		delete(sm.states, state)
		return ErrStateExpired
	}

	delete(sm.states, state)
	return nil
}

func (sm *StateManager) cleanup() {
	ticker := time.NewTicker(cleanUpInterval * time.Minute)
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
