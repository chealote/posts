package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	Title string
	Link string
}

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Database string `json:"database"`
}

type SQLite struct {
	conn *sql.DB
	config Config
}

// only scan a single row
func (d SQLite) RunQuery(query string, values ...any) error {
	rows, err := d.conn.Query(query)
	if err != nil {
		return err
	}
	if rows.Next() {
		return rows.Scan(values...)
	}
	return nil
}

func (d SQLite) Exec(query string) error {
	_, err := d.conn.Exec(query)
	return err
}

func Connect(config Config) (SQLite, error) {
	conn, err := sql.Open("sqlite3", "posts.db")
	if err != nil {
		return SQLite{}, err
	}

	return SQLite{conn, config}, err
}

/*
func (d SQLite) ListTopPosts() ([]Post, error) {
	query, err := d.getQuery("list-top-posts")
	if err != nil {
		return []Post{}, err
	}

	rows, err := d.conn.Query(query)
	if err != nil {
		return []Post{}, err
	}
	defer rows.Close()

	posts := []Post{}
	title := ""
	link := ""
	for rows.Next() {
		rows.Scan(&title, &link)
		posts = append(posts, Post{
			Title: title,
			Link: link,
		})
	}

	return posts, err
}
*/
