package model

import (
	"fmt"
	"time"
)

type RedirectURL string

type Email string

func (u RedirectURL) Valid() bool {
	return len(u) > 0 && u[0] == '/'
}

type Identity interface {
	Id() (int, error)
	SetId(int)
	Saved() bool
}

type identity struct {
	id    int
	saved bool
}

func newIdentity(id int) identity {
	return identity{
		id:    id,
		saved: true,
	}
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

type User struct {
	identity
	email Email
	nome  string
	ctime time.Time
}

func (u User) Name() string {
	return u.nome
}

func NewUser(id int, email Email, nome string, ctime time.Time) User {
	return User{
		identity: newIdentity(id),
		email:    email,
		nome:     nome,
		ctime:    ctime,
	}
}
