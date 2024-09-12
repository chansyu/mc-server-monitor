package mocks

import "github.com/itzsBananas/mc-server-monitor/internal/models"

type UserModel struct{}

func (m *UserModel) Insert(username, password string) error {
	switch username {
	case "dupe@example.com":
		return models.ErrDuplicateUsername
	default:
		return nil
	}
}
func (m *UserModel) Authenticate(username, password string) (int, error) {
	if username == "alice@example.com" && password == "pa$$word" {
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
