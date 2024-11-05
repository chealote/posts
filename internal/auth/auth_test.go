package auth

import (
	"fmt"
	"testing"
)

var (
	sessionStorage map[string]string
	userStorage    map[string]string
)

type mockDbStruct struct{}

func (m mockDbStruct) LookupSession(session string) (bool, error) {
	for _, s := range sessionStorage {
		if session == s {
			return true, nil
		}
	}
	return false, nil
}

func (m mockDbStruct) CreateReplaceSession(username string, session string) error {
	sessionStorage[username] = session
	return nil
}

func (m mockDbStruct) RegisterUser(username string, password string) error {
	userStorage[username] = password
	return nil
}

func (m mockDbStruct) DeleteSession(session string) error {
	for _, s := range sessionStorage {
		if s == session {
			delete(sessionStorage, session)
			return nil
		}
	}
	return fmt.Errorf("session not found")
}

func (m mockDbStruct) CheckValidUserCredentials(user string, pass string) (bool, error) {
	for u, p := range userStorage {
		if u == user && p == pass {
			return true, nil
		}
	}
	return false, nil
}

func init() {
	mockDb := mockDbStruct{}
	sessionStorage = make(map[string]string)
	userStorage = make(map[string]string)
	AuthDb = mockDb
}

func Test_ValidateAuthorization(t *testing.T) {
	session := "somesession"
	sessionStorage["invalid"] = session
	defer delete(sessionStorage, "invalid")

	ok, err := ValidateAuthorization(session)
	if err != nil {
		t.Fatalf("some error: %s", err)
	}
	if !ok {
		t.Fatalf("not authorized")
	}
}

func Test_RegisterUser(t *testing.T) {
	user := "user"
	pass := "pass"

	RegisterUser(user, pass)
	for u, p := range userStorage {
		if u == user && p == pass {
			return
		}
	}
	t.Fatalf("user not found")
}

func Test_Login(t *testing.T) {
	user := "user"
	pass := "pass"
	userStorage[user] = pass
	defer delete(userStorage, user)

	session, err := Login(user, pass)
	if err != nil {
		t.Fatalf("failed storing token: %s", err)
	}

	for u, s := range sessionStorage {
		if u == user && s == session {
			return
		}
	}
	t.Fatalf("session not found for logged in user")
}

func Test_Login_InvalidCredentials(t *testing.T) {
	user := "user"
	goodPass := "elaboratepassword"
	invalidPass := "thewrongpassword"
	userStorage[user] = goodPass
	defer delete(userStorage, user)

	_, err := Login(user, invalidPass)
	if err == nil {
		t.Fatalf("login should be invalid")
	}
}

func Test_Logout(t *testing.T) {
	user := "user"
	session := "sessionToken"
	sessionStorage[user] = session
	defer delete(sessionStorage, user)

	if err := Logout(session); err != nil {
		t.Fatalf("failed loggin out")
	}
}
