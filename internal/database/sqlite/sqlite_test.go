package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"posts/internal/utils"
	"testing"
	"time"
	"posts/internal/handler"
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
	expectedPosts := []struct {
		id    string
		title string
		post  string
	}{
		{
			"123",
			"First post",
			"Content of the first post",
		},
		{
			"456",
			"Second post",
			"Content of the second post",
		},
	}

	for _, post := range expectedPosts {
		if err := conn.CreatePost(post.id, post.title, post.post); err != nil {
			t.Fatalf("failed to create post %s: %s", post.id, err)
		}
	}

	i := 0
	rows, err := dbConn.Query("SELECT id, title, post FROM posts")
	if err != nil {
		t.Fatalf("error querying posts: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		id := ""
		title := ""
		post := ""
		rows.Scan(&id, &title, &post)
		if id != expectedPosts[i].id &&
			title != expectedPosts[i].title &&
			post != expectedPosts[i].post {
			t.Fatalf("post %d don't match: %s != %s", i, id, expectedPosts[i].id)
		}

		i++
	}

	if len(expectedPosts) != i {
		t.Fatalf("not same number of posts")
	}

	cleanup()
}

func Test_ListPostWithIds(t *testing.T) {
	expectedPosts := []handler.PostWithId{
		{
			Id: "1",
			Title: "First Post",
		},
		{
			Id: "2",
			Title: "Second post",
		},
		{
			Id: "3",
			Title: "Third post",
		},
	}

	for _, post := range expectedPosts {
		query := fmt.Sprintf("INSERT INTO posts(id, title) VALUES ('%s', '%s')", post.Id, post.Title)
		_, err := dbConn.Exec(query)
		if err != nil {
			t.Errorf("failed to create post %s: %s", post, err)
		}
	}

	rows, err := dbConn.Query("SELECT id FROM posts")
	if err != nil {
		t.Errorf("failed retrieving ids from created posts")
	}

	ids := []string{}
	for rows.Next() {
		id := ""
		rows.Scan(&id)
		ids = append(ids, id)
	}

	if len(ids) != len(expectedPosts) {
		t.Errorf("len of ids don't match with len of expected posts")
	}

	posts, err := conn.ListWithId()
	if err != nil {
		t.Errorf("ListPostTitles returned an error: %v", err)
	}

	if len(posts) != len(expectedPosts) {
		t.Errorf("expected %d titles, got %d", len(expectedPosts), len(posts))
	}

	for i, post := range posts {
		if post.Title != expectedPosts[i].Title || post.Id != ids[i] {
			t.Errorf("posts %d don't match: expect=%q, got=%q", i, expectedPosts[i], post)
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
