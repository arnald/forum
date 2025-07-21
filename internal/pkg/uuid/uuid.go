package uuid

import (
	"github.com/google/uuid"
)

type Provider interface {
	NewUUID() string
}

func NewProvider() Provider {
	return uuidProvider{}
}

type uuidProvider struct{}

func (u uuidProvider) NewUUID() string {
	return uuid.New().String()
}
