package mocks

import "github.com/itzsBananas/mc-server-monitor/internal/models"

type UserModel struct{}

var Username = "alice@example.com"
var Password = "pa$$word"

func (m *UserModel) Insert(username, password string) error {
	switch username {
	case "dupe@example.com":
		return models.ErrDuplicateUsername
	default:
		return nil
	}
}
func (m *UserModel) Authenticate(username, password string) (int, error) {
	if username == Username && password == Password {
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
