package mocks

import "snippetbox.francisko/internal/models"

type UserModel struct{}

func (u *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicatedEmail
	default:
		return nil
	}
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "password" {
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
