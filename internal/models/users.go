package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (UserWithoutPassword, error)
	PasswordUpdate(id int, currentPassword, newPassword string) error
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserWithoutPassword struct {
	ID      int
	Name    string
	Email   string
	Created time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	var hashedPassword []byte
	stmt := `SELECT hashed_password FROM users WHERE id = ?`
	row := u.DB.QueryRow(stmt, id)
	err := row.Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		} else {
			return err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(currentPassword))
	if err != nil {
		return ErrInvalidCredentials
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	stmt = `UPDATE users SET hashed_password = ? WHERE id = ?`
	_, err = u.DB.Exec(stmt, hashedNewPassword, id)

	return nil
}

func (u *UserModel) Get(id int) (UserWithoutPassword, error) {
	stmt := `SELECT name, email, created FROM users WHERE id = ?`
	row := u.DB.QueryRow(stmt, id)

	var user UserWithoutPassword

	err := row.Scan(&user.Name, &user.Email, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserWithoutPassword{}, ErrNoRecord
		} else {
			return UserWithoutPassword{}, err
		}
	}

	return user, nil
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

func (u *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := u.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (u *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := u.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
