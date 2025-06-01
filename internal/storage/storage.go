package storage

import (
	"denet-app/internal/domain/users"
	"denet-app/internal/storage/postgres"

	"fmt"
)

type Storage interface {
	GetLeaderboard() ([]users.LeaderboardEntry, error)
	GetUserStatus(id string) (users.StatusInfo, error)
	PostTaskComplete(id, taskType string) error
	PostReferrer(id, referralCode, refferedBy string) error
}

func NewStorage(URL string) (Storage, error) {
	var storage Storage
	var err error

	switch {
	case URL[:11] == "postgres://":
		storage, err = postgres.NewPostgresStorage(URL)
		if err != nil {
			return nil, fmt.Errorf("failed to create Postgres storage: %w", err)
		}
		return storage, nil
		// Add other storage types here if needed
	}
	return nil, fmt.Errorf("unsupported storage type")
}
