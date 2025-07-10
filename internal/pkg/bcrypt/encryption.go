package bcrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const encryptionCost = 12

type Provider interface {
	Generate(plaintextPassword string) ([]byte, error)
	Matches(plaintextPassword string, encryptedPass []byte) (bool, error)
}

func NewProvider() Provider {
	return &encryptionProvider{}
}

type encryptionProvider struct{}

func (p *encryptionProvider) Generate(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), encryptionCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func (p *encryptionProvider) Matches(plaintextPassword string, encryptedPass []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(encryptedPass, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
