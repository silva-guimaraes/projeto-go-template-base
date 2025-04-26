package database

import (
	"database/sql"
	"foobar/model"
	"time"
)

func VerifyUser(tx *sql.Tx, email model.Email, senha string) (model.User, error) {
	row := tx.QueryRow(
		`SELECT id, nome, ctime, senha_hash FROM usuario WHERE email = $1 LIMIT 1`,
		email,
	)
	if row.Err() != nil {
		return model.User{}, row.Err()
	}

	var (
		id        int
		senhaHash string
		nome      string
		ctime     time.Time
	)
	err := row.Scan(&id, &nome, &ctime, &senhaHash)
	if err != nil {
		return model.User{}, err
	}

	if err = checkPasswordHash(senha, senhaHash); err != nil {
		return model.User{}, err
	}

	return model.NewUser(id, email, nome, ctime), nil
}

func FindUsuarioByID(tx *sql.Tx, id int) (model.User, error) {
	row := tx.QueryRow(
		`SELECT id, nome, email, ctime FROM usuario WHERE id = $1 LIMIT 1`,
		id,
	)
	if row.Err() != nil {
		return model.User{}, row.Err()
	}
	var (
		nome  string
		email model.Email
		ctime time.Time
	)
	err := row.Scan(&id, &nome, &email, &ctime)
	if err != nil {
		return model.User{}, err
	}
	return model.NewUser(id, email, nome, ctime), nil
}
