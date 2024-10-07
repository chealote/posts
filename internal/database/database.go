package database

import (
	"errors"
)

var (
	NotesTablename   = "notes"
	SessionTablename = "sessions"
	UsersTablename   = "users"

	ErrUnauthorized = errors.New("Unauthorized")
	Database  Implementation
)

type Implementation interface {
	Initialize() error

	// auth
	LookupSession(session string) (bool, error)
	RegisterUser(username string, password string) error
	CreateSession(username string, password string) (string, error)
	DeleteSession(token string) error
}
