package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func getConnectionString(cfg Cfg) string {
	if cfg.Database == "" {
		return fmt.Sprintf("postgresql://%s:%s@%s/?sslmode=disable",
			cfg.Username, cfg.Password, cfg.Hostname)
	}
	return fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Hostname, cfg.Database)
}

func Initialize(cfg Cfg) error {
	connStr := getConnectionString(cfg)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	_, err = conn.Exec(fmt.Sprintf(`CREATE DATABASE %s`, cfg.Database))
	if err != nil {
		return err
	}

	_, err = conn.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (username text, token text, expires timestamp)`, sessionTablename))
	_, err = conn.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (username text, secret text)`, usersTablename))
	return err
}

func Connect(cfg Cfg) (Database, error) {
	connStr := getConnectionString(cfg)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return Database{}, err
	}

	return Database{conn}, err
}
