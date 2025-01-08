package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, email, password string) (int, error) {
	stmt := `INSERT into users (name, email, hashedPassword, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	result, err := u.DB.Exec(stmt, name, email, password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (u *UserModel) Authenticate(email, password string) error {
	return nil
}

func (u *UserModel) Exists(id int) (bool, error) {
	return true, nil
}
