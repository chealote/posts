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

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Database string `json:"database"`
	ScriptsPath string `json:"scriptsPath"`
}

type Implementation interface {
	Initialize() error

	// auth
	LookupSession(session string) (bool, error)
	RegisterUser(username string, secret string) error
	CreateSession(username string, secret string) (string, error)
}
