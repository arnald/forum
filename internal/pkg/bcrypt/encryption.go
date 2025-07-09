package bcrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Provider interface {
	Set(plaintextPassword string) error
	Matches(plaintextPassword string) (bool, error)
}

func NewProvider() Provider {
	return &encryptionProvider{}
}

type encryptionProvider struct {
	plaintext *string
	hash      []byte
}

func (p *encryptionProvider) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *encryptionProvider) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
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
