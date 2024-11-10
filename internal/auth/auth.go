package auth

import (
	"errors"
	"fmt"
	"posts/internal/database"
	"posts/internal/utils"
	"strings"
	"time"
)

type AuthDatabase interface {
	CreateReplaceSession(string, string) error
	LookupSession(string) (bool, error)
	RegisterUser(string, string, string) error
	DeleteSession(string) error
	CheckValidUserCredentials(string, string) (bool, error)
	RolesFromUser(string) (string, error)
}

const (
	RoleAdmin  = "admin"
	RolePoster = "poster"
	RoleReader = "reader"

	ActionPostRead   = "postRead"
	ActionPostCreate = "postCreate"
)

var (
	AuthDb AuthDatabase

	Permissions = map[string]string{
		RoleAdmin:  fmt.Sprintf("%s,%s", ActionPostRead, ActionPostCreate),
		RolePoster: fmt.Sprintf("%s,%s", ActionPostRead, ActionPostCreate),
		RoleReader: fmt.Sprintf("%s", ActionPostRead),
	}
)

func ValidateAuthorization(session string) (bool, error) {
	ok, err := AuthDb.LookupSession(session)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func RegisterUser(username string, password string) error {
	roles := ""
	// should store comma-based
	if username == "admin" {
		roles = fmt.Sprintf("%s", RoleAdmin)
	} else {
		roles = fmt.Sprintf("%s", RoleReader)
	}

	err := AuthDb.RegisterUser(username, password, roles)
	// TODO this is a debug check only, just return error without checking
	if errors.Is(err, database.ErrConstraintKey) {
		fmt.Println("user already exists!")
	}
	return err
}

func Login(username string, password string) (string, error) {
	ok, err := AuthDb.CheckValidUserCredentials(username, password)
	if err != nil {
		return "", err
	}
	if !ok {
		// TODO return 500 or 401 if !ok
		return "", fmt.Errorf("login is invalid")
	}

	epochNow := fmt.Sprintf("%d", time.Now().Unix())
	token := utils.Sha512Sum(fmt.Sprintf("%s%s", username, epochNow))
	return token, AuthDb.CreateReplaceSession(username, token)
}

func Logout(token string) error {
	return AuthDb.DeleteSession(token)
}

func IsUserAllowed(user string, action string) (bool, error) {
	role, err := AuthDb.RolesFromUser(user)
	if err != nil {
		return false, err
	}

	allowedActions, ok := Permissions[role]
	if !ok {
		return false, fmt.Errorf("user '%s' has no known allowed actions", user)
	}

	for _, a := range strings.Split(allowedActions, ",") {
		if a == action {
			return true, nil
		}
	}

	return false, nil
}
