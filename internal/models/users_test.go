package models_test

import (
	"testing"

	"github.com/itzsBananas/mc-server-monitor/internal/models"
)

func TestUserModelExists(t *testing.T) {
	testCases := []struct {
		userID int
		desc   string
		want   bool
	}{
		{
			userID: 0,
			desc:   "Zero ID",
			want:   false,
		},
		{
			userID: 1,
			desc:   "Valid ID",
			want:   true,
		},
		{
			userID: 100,
			desc:   "Not available ID",
			want:   false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			db := newTestDB(t)

			m := models.UserModel{db}

			exists, err := m.Exists(tC.userID)
			if exists != tC.want {
				t.Fatalf("expected %t; got %t: %v", tC.want, exists, err)
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUserCreation(t *testing.T) {
	username := "junior"
	password := "ilovecheese"
	db := newTestDB(t)

	m := models.UserModel{db}

	_, err := m.Authenticate(username, password)
	if err != models.ErrInvalidCredentials {
		t.Fatalf("Authentication should've failed...")
	}

	err = m.Insert(username, password)
	if err != nil {
		t.Fatal("Unexpected error when creating account: ", err)
	}

	id, err := m.Authenticate(username, password)
	if err != nil {
		t.Fatalf("Authentication should've passed after signing up...")
	}

	created, err := m.Exists(id)
	if err != nil {
		t.Fatal("Unexpected error when finding existing account: ", err)
	}
	if !created {
		t.Fatal("Expected to find the account, but didn't...")
	}
}
