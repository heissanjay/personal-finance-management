package storage

import (
	"context"
)

type DataStore interface {
	SaveUser(ctx context.Context, user User) error
	GetUserByUsername(ctx context.Context, username string) (User, error)
}

type User struct {
	ID       int
	Username string
	Password string
}
