package database

import (
	"database/sql"
	"time"
)

type Usuario struct {
	identity
	Email string
	Nome  string
	// senhaHash string `json:"-"`
	Ctime time.Time
}

func VerifyUser(tx *sql.Tx, email, senha string) (*Usuario, error) {
	row := tx.QueryRow(
		`SELECT id, nome, ctime, senha_hash FROM usuario WHERE email = $1 LIMIT 1`,
		email,
	)
	if row.Err() != nil {
		return nil, row.Err()
	}
	usuario := &Usuario{
		Email: email,
	}

	var id int
	var senhaHash string
	err := row.Scan(
		&id,
		&usuario.Nome,
		&usuario.Ctime,
		&senhaHash,
	)
	if err != nil {
		return nil, err
	}

	usuario.SetId(id)

	if err = checkPasswordHash(senha, senhaHash); err != nil {
		return nil, err
	}

	return usuario, nil
}

func FindUsuarioByID(tx *sql.Tx, id int) (*Usuario, error) {
	row := tx.QueryRow(
		`SELECT id, nome, email, ctime FROM usuario WHERE id = $1 LIMIT 1`,
		id,
	)
	if row.Err() != nil {
		return nil, row.Err()
	}
	usuario := &Usuario{}
	err := row.Scan(
		&id,
		&usuario.Nome,
		&usuario.Email,
		&usuario.Ctime,
	)
	if err != nil {
		return nil, err
	}
	usuario.SetId(id)
	return usuario, nil
}
