package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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

func (u *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT into users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = u.DB.Exec(stmt, name, email, hashedPassword)
	if err != nil {
		var mySQLErr *mysql.MySQLError

		if errors.As(err, &mySQLErr) {
			if mySQLErr.Number == 1062 && strings.Contains(mySQLErr.Message, "users_uc_email") {
				return ErrDuplicatedEmail
			}
		}
		return err
	}

	return nil
}

func (u *UserModel) Authenticate(email, password string) error {
	return nil
}

func (u *UserModel) Exists(id int) (bool, error) {
	return true, nil
}
