package database

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrUnauthorized = errors.New("error unauthorized")
)

type Connection interface {
	RunQuery(string, ...any) error
	Exec(string) error
}

type Database struct {
	Conn Connection
	ScriptsPath string
}

func (d Database) RunQuery(query string) error {
	return d.Conn.RunQuery(query)
}

func (d Database) getQuery(script string) (string, error) {
	content, err := os.ReadFile(fmt.Sprintf("%s/%s.sql", d.ScriptsPath, script))
	return string(content), err
}

