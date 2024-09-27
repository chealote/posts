package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func getConnectionString(config config) string {
	if config.Database == "" {
		return fmt.Sprintf("postgresql://%s:%s@%s/?sslmode=disable",
			config.Username, config.Password, config.Hostname)
	}
	return fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		config.Username, config.Password, config.Hostname, config.Database)
}

func Initialize(config config) error {
	connStr := getConnectionString(config)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	_, err = conn.Exec(fmt.Sprintf(`CREATE DATABASE %s`, config.Database))
	if err != nil {
		return err
	}

	_, err = conn.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (username text, token text, expires timestamp)`, sessionTablename))
	_, err = conn.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (username text, password text)`, usersTablename))
	return err
}

func Connect(config config) (Database, error) {
	connStr := getConnectionString(config)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return Database{}, err
	}

	return Database{conn}, err
}
