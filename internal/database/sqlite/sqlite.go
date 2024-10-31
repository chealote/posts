package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"posts/internal/database"
	"posts/internal/utils"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	InvalidConfigError = errors.New("Invalid config")
	ErrUnauthorized    = errors.New("Unauthorized")
)

type Config struct {
	Filename    string `json:"filename"`
	ScriptsPath string `json:"scriptsPath"`
}

type SQLite struct {
	conn *sql.DB

	filename    string
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

func Connect(config Config) (SQLite, error) {
	if config.Filename == "" || config.ScriptsPath == "" {
		return SQLite{}, InvalidConfigError
	}

	conn, err := sql.Open("sqlite3", config.Filename)
	if err != nil {
		return SQLite{}, err
	}

	return SQLite{conn, config.Filename, config.ScriptsPath}, err
}

func (d SQLite) getQuery(script string) (string, error) {
	fmt.Println("getQuery: config info: ", d.filename, d.scriptsPath)
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
	if !rows.Next() {
		fmt.Printf("LookupSession: no session found: %s\n", session)
		return false, err
	}

	rows.Scan(&token, &expires, &now)
	fmt.Printf("LookupSession: session=%s expires='%s' now='%s'\n", session, expires, now)
	return true, nil
	}
}

func (d SQLite) RegisterUser(username, password string) error {
	query, err := d.getQuery("register-user")
	if err != nil {
		return err
	}

	// TODO salt the password before save
	salt := "randomString"
	saltedPassword := utils.Sha512Sum(fmt.Sprintf("%s%s", password, salt))

	_, err = d.conn.Exec(query, username, saltedPassword, salt)
	if err != nil && strings.Contains(err.Error(), "CONSTRAINT") {
		return database.ErrConstraintKey
	}
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
		return false, ErrUnauthorized
	}

	dbSaltedPassword := ""
	dbSalt := ""
	rows.Scan(&dbSaltedPassword, &dbSalt)

	saltedPassword := utils.Sha512Sum(fmt.Sprintf("%s%s", password, dbSalt))

	if dbSaltedPassword != saltedPassword {
		return false, ErrUnauthorized
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

func (d SQLite) deleteUserSession(username string) error {
	query, err := d.getQuery("delete-user-session")
	if err != nil {
		return err
	}
	_, err = d.conn.Exec(query, username)
	return err
}

func (d SQLite) CheckValidUserCredentials(username string, password string) (bool, error) {
	return d.checkUserCredentials(username, password)
}

func (d SQLite) CreateReplaceSession(username string, session string) error {
	ok, err := d.checkExistingSession(username)
	if err != nil {
		return err
	}
	if ok {
		if err := d.deleteUserSession(username); err != nil {
			return err
		}
	}

	query, err := d.getQuery("create-session")
	if err != nil {
		return err
	}
	_, err = d.conn.Exec(query, username, session)
	return err
}

func (d SQLite) DeleteSession(token string) error {
	query, err := d.getQuery("delete-token-session")
	if err != nil {
		return err
	}

	_, err = d.conn.Exec(query, token)
	return err
}

func (d SQLite) ListPostTitles() ([]string, error) {
	query, err := d.getQuery("list-post-titles")
	if err != nil {
		return []string{}, err
	}

	rows, err := d.conn.Query(query)
	if err != nil {
		return []string{}, err
	}

	titles := []string{}
	title := ""
	for rows.Next() {
		rows.Scan(&title)
		titles = append(titles, title)
	}

	return titles, nil
}

func (d SQLite) CreatePost(title string) error {
	query, err := d.getQuery("create-post")
	if err != nil {
		return err
	}

	_, err = d.conn.Exec(query, title)
	return err
}
