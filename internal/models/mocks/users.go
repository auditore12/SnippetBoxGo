package mocks

import (
	"cmd/web/main.go/internal/models"
	"time"
)

type UserModel struct{}

var mockUsers = models.User{
	ID:      1,
	Name:    "An old silent pond",
	Email:   "blablab@gmail.com",
	Created: time.Now(),
}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}
func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) ShowUser() ([]models.User, error) {
	return []models.User{mockUsers}, nil
}
func (m *UserModel) GetUser(id int, email string) (models.User, error) {
	switch email {
	case "":
		return mockUsers, nil
	default:
		return models.User{}, models.ErrNoRecord
	}
}

func (m *UserModel) UpdateUsers(id int ,email string, name string) (int, error) {
	return 2, nil
}