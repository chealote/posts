package auth

import (
	"fmt"
	"testing"
)

type userStorageValue struct {
	password string
	roles    string
}

var (
	sessionStorage map[string]string
	userStorage    map[string]userStorageValue
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

func (m mockDbStruct) RegisterUser(username string, password string, roles string) error {
	userStorage[username] = userStorageValue{
		password,
		roles,
	}
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
		if u == user && p.password == pass {
			return true, nil
		}
	}
	return false, nil
}

func (m mockDbStruct) RolesFromUser(user string) (string, error) {
	for u, p := range userStorage {
		if u == user {
			return p.roles, nil
		}
	}
	return "", fmt.Errorf("user not found")
}

func init() {
	mockDb := mockDbStruct{}
	sessionStorage = make(map[string]string)
	userStorage = make(map[string]userStorageValue)
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

	if err := RegisterUser(user, pass); err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	for u, p := range userStorage {
		if u == user && p.password == pass {
			return
		}
	}
	t.Fatalf("user not found")
}

func Test_Login(t *testing.T) {
	user := "user"
	pass := "pass"
	userStorage[user] = userStorageValue{
		pass,
		"",
	}
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
	userStorage[user] = userStorageValue{
		goodPass,
		"",
	}
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

func Test_PermissionsUser(t *testing.T) {
	user := "user"
	pass := "user"
	actionShouldAllow := ActionPostRead
	actionShouldDeny := ActionPostCreate

	if err := RegisterUser(user, pass); err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	ok, err := IsUserAllowed(user, actionShouldAllow)
	if err != nil {
		t.Fatalf("failed to get roles for user '%s': %v", user, err)
	}

	if !ok {
		t.Fatalf("user %s is not allowed to perform action %s", user, actionShouldAllow)
	}

	ok, err = IsUserAllowed(user, actionShouldDeny)
	if err != nil {
		t.Fatalf("failed to get roles for user '%s': %v", user, err)
	}

	if ok {
		t.Fatalf("user %s is allowed to perform action %s", user, actionShouldDeny)
	}
}

func Test_PermissionsAdmin(t *testing.T) {
	user := "admin"
	pass := "admin"
	actionShouldAllow := ActionPostRead
	actionShouldDeny := ActionPostCreate

	if err := RegisterUser(user, pass); err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	ok, err := IsUserAllowed(user, actionShouldAllow)
	if err != nil {
		t.Fatalf("failed to get roles for user '%s': %v", user, err)
	}

	if !ok {
		t.Fatalf("user %s is not allowed to perform action %s", user, actionShouldAllow)
	}

	ok, err = IsUserAllowed(user, actionShouldDeny)
	if err != nil {
		t.Fatalf("failed to get roles for user '%s': %v", user, err)
	}

	if !ok {
		t.Fatalf("user %s is not allowed to perform action %s", user, actionShouldDeny)
	}
}
