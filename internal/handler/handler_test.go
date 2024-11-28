package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type mockAuthService struct{}

func (a mockAuthService) ValidateAuthorization(s string) (bool, error) {
	return true, nil
}

func (a mockAuthService) RegisterUser(user, pass string) error {
	return nil
}

func (a mockAuthService) Login(user, pass string) (string, error) {
	return "token", nil
}

func (a mockAuthService) Logout(token string) error {
	return nil
}

type responseWriter struct {
	*bytes.Buffer
	status *int
}

type postDb struct {
	posts []PostWithId
}

func (p postDb) CreatePost(id, title, content string) error {
	return nil
}

func (p postDb) ContentsPost(id string) (PostContent, error) {
	return PostContent{}, nil
}

func (p postDb) ListWithId() ([]PostWithId, error) {
	return p.posts, nil
}

func (rw responseWriter) Header() http.Header {
	return map[string][]string{}
}

func (rw responseWriter) WriteHeader(s int) {
	*rw.status = s
}

func slicesAreEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i, _ := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func valueEndIndex(s []byte) int {
	for i, v := range s {
		if v == byte(0) {
			return i
		}
	}
	return -1
}

func Test_HandlePostRoot_WhenGetPostRoot_GetListOfPosts(t *testing.T) {
	n := 0
	r := http.Request{
		Method: "GET",
	}
	w := responseWriter{
		bytes.NewBuffer(nil),
		&n,
	}

	expectedPosts := []PostWithId{
		{"1", "One post"},
		{"2", "Two post"},
		{"3", "Three post"},
	}

	PostDb = postDb{
		expectedPosts,
	}

	HandlePostRoot(w, &r)

	buffer := make([]byte, 1000)
	_, err := w.Read(buffer)
	if err != nil {
		t.Fatalf("cannot read well")
	}

	b, err := json.Marshal(expectedPosts)
	if err != nil {
		t.Fatalf("I can't marshal the posts: %s", err)
	}

	buffer = buffer[0:valueEndIndex(buffer)]
	if !slicesAreEqual(b, buffer) {
		t.Fatalf("slices are not equal: %v != %v", b, buffer)
	}
	if *w.status != http.StatusOK {
		t.Fatalf("response code is not %d: %d", http.StatusOK, *w.status)
	}
}

func Test_HandlePostRoot_WhenCreatePostValid_PostGetsCreated(t *testing.T) {
	expectedPost := PostCreate{
		Id:    "123",
		Title: "The post title",
		Post:  "The post contents",
	}

	expectedPosts := []PostWithId{
		{expectedPost.Id, expectedPost.Title},
	}

	PostDb = postDb{
		expectedPosts,
	}

	b, err := json.Marshal(expectedPost)
	if err != nil {
		t.Fatalf("failed to marshal post: %v", err)
	}

	r := http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewBuffer(b)),
	}

	n := 0
	w := responseWriter{
		bytes.NewBuffer(b),
		&n,
	}
	HandlePostRoot(w, &r)

	if *w.status != http.StatusCreated {
		t.Fatalf("Failed creating post: %d", *w.status)
	}

	r = http.Request{
		Method: "GET",
	}
	n = 0
	w = responseWriter{
		bytes.NewBuffer(nil),
		&n,
	}
	HandlePostRoot(w, &r)

	buffer := make([]byte, 1000)
	if _, err := w.Read(buffer); err != nil {
		t.Fatalf("cannot read well: %v", err)
	}

	b, err = json.Marshal(expectedPosts)
	if err != nil {
		t.Fatalf("I can't marshal the posts: %s", err)
	}

	buffer = buffer[0:valueEndIndex(buffer)]
	if !slicesAreEqual(b, buffer) {
		t.Fatalf("slices are not equal: %v != %v (%s != %s)", b, buffer, string(b), string(buffer))
	}
	if *w.status != http.StatusOK {
		t.Fatalf("response code is not %d: %d", http.StatusOK, *w.status)
	}
}
