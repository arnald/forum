package bcrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const encryptionCost = 12

type Provider interface {
	Generate(plaintextPassword string) (string, error)
	Matches(databasePassword string, passwordFromRequest string) error
}

func NewProvider() Provider {
	return &encryptionProvider{}
}

type encryptionProvider struct{}

func (p *encryptionProvider) Generate(plaintextPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), encryptionCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (p *encryptionProvider) Matches(databasePassword string, passwordFromRequest string) error {
	err := bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(passwordFromRequest))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return err
		default:
			return err
		}
	}

	return nil
}
