package core

import (
	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	Generate(string) (string, error)
	Verify(string, string) bool
}

func NewHasher() Hasher {
	return &Crypto{}
}

type Crypto struct{}

func (c *Crypto) Generate(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}

	return string(hash), nil
}

func (c *Crypto) Verify(hash, passowrd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passowrd)); err != nil {
		return false
	}

	return true
}
