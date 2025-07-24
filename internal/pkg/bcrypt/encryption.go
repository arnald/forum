package bcrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const encryptionCost = 12

type Provider interface {
	Generate(plaintextPassword string) (string, error)
	Matches(plaintextPassword string, encryptedPass string) error
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

func (p *encryptionProvider) Matches(encryptedPass string, passwordToCheck string) error {
	hashedPassword, err := p.Generate(passwordToCheck)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(encryptedPass), []byte(hashedPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return nil
		default:
			return err
		}
	}

	return nil
}
