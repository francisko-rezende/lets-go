package mocks

import (
	"time"

	"snippetbox.francisko/internal/models"
)

type UserModel struct{}

var mockUserWithoutPassword = models.UserWithoutPassword{
	ID:      1,
	Name:    "Alice",
	Email:   "alice@example.com",
	Created: time.Now(),
}

func (u *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	return nil
}

func (u *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicatedEmail
	default:
		return nil
	}
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (u *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (u *UserModel) Get(id int) (models.UserWithoutPassword, error) {
	switch id {
	case 1:
		return mockUserWithoutPassword, nil
	default:
		return models.UserWithoutPassword{}, models.ErrNoRecord
	}
}
