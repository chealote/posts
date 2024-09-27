package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"errors"
	"posts/internal/database"

	_ "github.com/mattn/go-sqlite3"
)

var (
	InvalidConfigError = errors.New("Invalid config")
)

type Config struct {
	Filename string `json:"filename"`
	ScriptsPath string `json:"scriptsPath"`
}

type SQLite struct {
	conn *sql.DB

	filename string
	scriptsPath string
}

func (d SQLite) Initialize() error {
	query, err := d.getQuery("initialize")
	if err != nil {
		return err
	}

	_, err = d.conn.Exec(query)
	return err
}

func Connect(config Config) (database.Implementation, error) {
	if config.Filename == "" || config.ScriptsPath == "" {
		return nil, InvalidConfigError
	}

	conn, err := sql.Open("sqlite3", config.Filename)
	if err != nil {
		return SQLite{}, err
	}

	return SQLite{conn, config.Filename, config.ScriptsPath}, err
}

func (d SQLite) getQuery(script string) (string, error) {
	fmt.Println("SOME CONFIG:", d.filename, d.scriptsPath)
	filepath := fmt.Sprintf("%s/%s.sql", d.scriptsPath, script)
	content, err := os.ReadFile(filepath)
	fmt.Printf("getQuery(): FILE: %s -- QUERY: %s\n", filepath, string(content))
	return string(content), err
}

func (d SQLite) LookupSession(session string) (bool, error) {
	query, err := d.getQuery("lookup-session")
	if err != nil {
		return false, err
	}
	rows, err := d.conn.Query(query, session)
	if err != nil {
		fmt.Println("ERROR: LookupSession:", err)
		return false, err
	}
	defer rows.Close()
	token := ""
	expires := ""
	now := ""
	for rows.Next() {
		rows.Scan(&token, &expires, &now)
		fmt.Printf("LookupSession: session=%s expires='%s' now='%s'\n", session, expires, now)
		return true, nil
	}
	fmt.Printf("LookupSession: no session found: %s\n", session)
	return false, err
}

func (d SQLite) RegisterUser(username, password string) error {
	query, err := d.getQuery("register-user")
	if err != nil {
		return err
	}
	_, err = d.conn.Exec(query, username, password)
	return err
}

func (d SQLite) checkUserCredentials(username, password string) (bool, error) {
	query, err := d.getQuery("check-user-credentials")
	if err != nil {
		return false, err
	}

	rows, err := d.conn.Query(query, username)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if !rows.Next() {
		return false, database.ErrUnauthorized
	}

	dbPassword := ""
	rows.Scan(&dbPassword)
	if dbPassword != password {
		return false, database.ErrUnauthorized
	}

	return true, nil
}

func (d SQLite) checkExistingSession(username string) (bool, error) {
	query, err := d.getQuery("check-existing-session")
	if err != nil {
		return false, err
	}
	rows, err := d.conn.Query(query, username)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (d SQLite) deleteCurrentSession(username string) error {
	query, err := d.getQuery("delete-current-session")
	if err != nil {
		return err
	}
	_, err = d.conn.Exec(query, username)
	return err
}

func (d SQLite) CreateSession(username string, password string) (string, error) {
	if ok, err := d.checkUserCredentials(username, password); err != nil || !ok {
		return "", err
	}

	ok, err := d.checkExistingSession(username)
	if err != nil {
		return "", err
	}
	if ok {
		if err := d.deleteCurrentSession(username); err != nil {
			return "", err
		}
	}

	// TODO maybe change this
	token := "crypticToken"
	query, err := d.getQuery("create-session")
	if err != nil {
		return "", err
	}
	_, err = d.conn.Exec(query, username, token)
	return token, err
}
