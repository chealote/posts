package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"posts/internal/database"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	conn *sql.DB
	cfg database.Cfg
}

func (d SQLite) Initialize() error {
	query, err := d.getQuery("initialize")
	if err != nil {
		return err
	}

	fmt.Println("running the query:", query)
	_, err = d.conn.Exec(query)
	return err
}

func Connect(cfg database.Cfg) (database.Implementation, error) {
	conn, err := sql.Open("sqlite3", "posts.db")
	if err != nil {
		return SQLite{}, err
	}

	return SQLite{conn, cfg}, err
}

func (d SQLite) getQuery(script string) (string, error) {
	content, err := os.ReadFile(fmt.Sprintf("%s/%s.sql", d.cfg.ScriptsPath, script))
	return string(content), err
}

func (d SQLite) LookupSession(session string) (bool, error) {
	query, err := d.getQuery("lookup-session")
	if err != nil {
		return false, err
	}
	rows, err := d.conn.Query(query, session)
	if err != nil {
		fmt.Println("LookupSession:", err)
		return false, err
	}
	defer rows.Close()
	token := ""
	for rows.Next() {
		rows.Scan(&token)
		fmt.Println("LookupSession: found token")
		return true, nil
	}
	return false, err
}

func (d SQLite) RegisterUser(username, secret string) error {
	query, err := d.getQuery("register-user")
	if err != nil {
		return err
	}
	_, err = d.conn.Exec(query, username, secret)
	return err
}

func (d SQLite) checkUserCredentials(username, secret string) (bool, error) {
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

	dbSecret := ""
	rows.Scan(&dbSecret)
	if dbSecret != secret {
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

func (d SQLite) CreateSession(username string, secret string) (string, error) {
	if ok, err := d.checkUserCredentials(username, secret); err != nil || !ok {
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

func (d SQLite) ListTopPosts() ([]database.Post, error) {
	query, err := d.getQuery("list-top-posts")
	if err != nil {
		return []database.Post{}, err
	}

	rows, err := d.conn.Query(query)
	if err != nil {
		return []database.Post{}, err
	}
	defer rows.Close()

	posts := []database.Post{}
	title := ""
	link := ""
	for rows.Next() {
		rows.Scan(&title, &link)
		posts = append(posts, database.Post{
			Title: title,
			Link: link,
		})
	}

	return posts, err
}
