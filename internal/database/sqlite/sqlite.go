package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"posts/internal/database"
	"posts/internal/handler"
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

func (d SQLite) Initialize() error {
	query, err := d.getQuery("initialize")
	if err != nil {
		return err
	}

	_, err = d.conn.Exec(query)
	return err
}

func (d SQLite) Close() {
	d.conn.Close()
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

func (d SQLite) RegisterUser(username, password, roles string) error {
	query, err := d.getQuery("register-user")
	if err != nil {
		return err
	}

	// TODO salt the password before save
	salt := "randomString"
	saltedPassword := utils.Sha512Sum(fmt.Sprintf("%s%s", password, salt))

	_, err = d.conn.Exec(query, username, saltedPassword, salt, roles)
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

func (d SQLite) ListWithId() ([]handler.PostWithId, error) {
	query, err := d.getQuery("list-post-titles-ids")
	if err != nil {
		return []handler.PostWithId{}, err
	}

	rows, err := d.conn.Query(query)
	if err != nil {
		return []handler.PostWithId{}, err
	}
	defer rows.Close()

	titles := []handler.PostWithId{}
	title := ""
	id := ""
	for rows.Next() {
		rows.Scan(&id, &title)
		titles = append(titles, handler.PostWithId{
			Id:    id,
			Title: title,
		})
	}

	return titles, nil
}

func (d SQLite) CreatePost(postId string, title string, post string) error {

	query, err := d.getQuery("create-post")
	if err != nil {
		return err
	}

	_, err = d.conn.Exec(query, postId, title, post)
	return err
}

func (d SQLite) ContentsPost(id string) (handler.PostContent, error) {
	query, err := d.getQuery("get-post-contents")
	if err != nil {
		return handler.PostContent{}, err
	}

	rows, err := d.conn.Query(query, id)
	if err != nil {
		return handler.PostContent{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return handler.PostContent{}, fmt.Errorf("empty response")
	}

	title := ""
	contents := ""
	err = rows.Scan(&title, &contents)

	return handler.PostContent{
		Title:    title,
		Contents: contents,
	}, err
}

func (d SQLite) RolesFromUser(user string) (string, error) {
	query, err := d.getQuery("roles-from-user")
	if err != nil {
		return "", err
	}

	rows, err := d.conn.Query(query, user)
	if err != nil {
		return "", err
	}

	if !rows.Next() {
		return "", fmt.Errorf("couldn't find roles for user '%s'", user)
	}

	roles := ""
	err = rows.Scan(&roles)
	return roles, nil
}
