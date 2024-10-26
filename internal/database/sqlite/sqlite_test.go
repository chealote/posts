package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"posts/internal/utils"
	"testing"
	"time"
)

var (
	conn   SQLite
	dbConn *sql.DB

	dbFilepath = "/tmp/test.db"
)

func init() {
	os.Remove(dbFilepath)

	sqlite, err := Connect(Config{
		Filename:    dbFilepath,
		ScriptsPath: "../../../scripts",
	})

	if err != nil {
		fmt.Errorf("failed to open database: %v", err)
		return
	}

	conn = sqlite

	if err := conn.Initialize(); err != nil {
		fmt.Errorf("failed to initialize db: %v", err)
		return
	}

	dbconn, err := sql.Open("sqlite3", dbFilepath)
	if err != nil {
		fmt.Errorf("cannot establish my own testing connection: %v", err)
		return
	}

	dbConn = dbconn
}

// TODO replace this with something more fancy
func cleanup() {
	_, err := dbConn.Exec("DELETE FROM posts")
	if err != nil {
		panic(err)
	}
}

func Test_CreatePost(t *testing.T) {
	expectedPosts := []string{
		"First post",
		"Second post",
		"Third post",
	}

	for _, post := range expectedPosts {
		if err := conn.CreatePost(post); err != nil {
			t.Fatalf("failed to create post %s: %s", post, err)
		}
	}

	i := 0
	rows, err := dbConn.Query("SELECT title FROM posts")
	if err != nil {
		t.Fatalf("error querying posts")
	}
	defer rows.Close()

	for rows.Next() {
		post := ""
		rows.Scan(&post)
		if post != expectedPosts[i] {
			t.Fatalf("post %d don't match: %s != %s", i, post, expectedPosts)
		}

		i++
	}

	if len(expectedPosts) != i {
		t.Fatalf("not same number of posts")
	}

	cleanup()
}

func Test_ListPostTitles(t *testing.T) {
	expectedPosts := []string{
		"First post",
		"Second post",
		"Third post",
	}

	for _, post := range expectedPosts {
		query := fmt.Sprintf("INSERT INTO posts(title) VALUES ('%s')", post)
		_, err := dbConn.Exec(query)
		if err != nil {
			t.Errorf("failed to create post %s: %s", post, err)
		}
	}

	titles, err := conn.ListPostTitles()
	if err != nil {
		t.Errorf("ListPostTitles returned an error: %v", err)
	}

	if len(titles) != len(expectedPosts) {
		t.Errorf("expected %d titles, got %d", len(expectedPosts), len(titles))
	}

	for i, title := range titles {
		if title != expectedPosts[i] {
			t.Errorf("posts don't match %d to be %q, got %q", i, expectedPosts[i], title)
		}
	}

	cleanup()
}

func Test_RegisterUser(t *testing.T) {
	expectedUser := "admin"
	expectedPass := "admin"
	if err := conn.RegisterUser(expectedUser, expectedPass); err != nil {
		t.Errorf("failed to register new user")
	}

	rows, err := dbConn.Query("SELECT username, password, salt FROM users")
	if err != nil {
		t.Errorf("failed querying for user")
	}
	defer rows.Close()

	if !rows.Next() {
		t.Errorf("no user brought back")
	}

	user := ""
	salted := ""
	salt := ""
	if err := rows.Scan(&user, &salted, &salt); err != nil {
		t.Errorf("failed scanning rows")
	}

	expectedSalted := utils.Sha512Sum(fmt.Sprintf("%s%s", expectedPass, salt))

	if user != expectedUser {
		t.Errorf("users don't match, got=%s, expect=%s", user, expectedUser)
	}

	if salted != expectedSalted {
		t.Errorf("salted pass don't match, got=%s, expect=%s", salted, expectedSalted)
	}
}

func Test_CreateReplaceSession(t *testing.T) {
	expectedUser := "admin"
	expectedSession := "session"

	if err := conn.CreateReplaceSession(expectedUser, expectedSession); err != nil {
		t.Errorf("failed to create a session: %v", err)
	}

	rows, err := dbConn.Query("SELECT username, expires FROM sessions WHERE expires > DATETIME('now')")
	if err != nil {
		t.Errorf("failed querying for user")
	}
	defer rows.Close()

	if !rows.Next() {
		t.Errorf("no active session currently")
	}

	user := ""
	expires := ""
	if err := rows.Scan(&user, &expires); err != nil {
		t.Errorf("error reading user from response")
	}

	layout := "2006-01-02T15:04:05Z"
	parsedExpire, err := time.Parse(layout, expires)
	if err != nil {
		t.Errorf("error parsing date: %s: %v", expires, err)
	}

	now := time.Now().UTC()

	if now.After(parsedExpire) {
		t.Errorf("invalid expire time for session: %s", parsedExpire)
	}
}

func Test_CreateReplaceSession_InvalidSession(t *testing.T) {
	expectedUser := "admin"
	expectedSession := "session"

	if err := conn.CreateReplaceSession(expectedUser, expectedSession); err != nil {
		t.Errorf("failed to create a session: %v", err)
	}

	rows, err := dbConn.Query("SELECT username, expires FROM sessions WHERE expires > DATETIME('now')")
	if err != nil {
		t.Errorf("failed querying for user")
	}
	defer rows.Close()

	if !rows.Next() {
		t.Errorf("no active session currently")
	}

	user := ""
	expires := ""
	if err := rows.Scan(&user, &expires); err != nil {
		t.Errorf("error reading user from response")
	}

	layout := "2006-01-02T15:04:05Z"
	parsedExpire, err := time.Parse(layout, expires)
	if err != nil {
		t.Errorf("error parsing date: %s: %v", expires, err)
	}

	future := time.Now().UTC().Add(time.Second * time.Duration(300))

	if future.Before(parsedExpire) {
		t.Errorf("invalid expire time for session: %s, %s", parsedExpire, future)
	}
}

