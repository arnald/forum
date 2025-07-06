package uuid

import (
	"github.com/google/uuid"
)

type Provider interface {
	NewUUID() uuid.UUID
}

func NewProvider() Provider {
	return uuidProvider{}
}

type uuidProvider struct{}

func (u uuidProvider) NewUUID() uuid.UUID {
	return uuid.New()
}
