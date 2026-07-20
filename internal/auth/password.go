package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	Plaintext string
	Hash      []byte
}

const bcryptCost = 12

func (p *Password) HashPassword(passwordPlaintext string) error {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(passwordPlaintext), bcryptCost)
	if err != nil {
		return err
	}

	p.Plaintext = ""
	p.Hash = passwordHash

	return nil
}

func (p *Password) Matches(passwordPlaintext string) (bool, error) {
	if len(p.Hash) == 0 {
		return false, errors.New("missing password hash")
	}

	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(passwordPlaintext))
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
