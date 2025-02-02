package database

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Model interface {
	Id() (int, error)
	SetId(int)
	Saved() bool
}

type identity struct {
	id    int
	saved bool
}

var (
	ErrNotSaved = fmt.Errorf("primary key id não foi salva")
)

func (p *identity) Id() (int, error) {
	if !p.saved {
		return 0, ErrNotSaved
	}
	return p.id, nil
}

func (p *identity) SetId(i int) {
	p.id = i
	p.saved = true
}

func (p *identity) Saved() bool {
	return p.saved
}

func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
