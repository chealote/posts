package auth

import (
	"errors"
	"fmt"
	"posts/internal/database"
	"posts/internal/utils"
	"time"
)

type AuthDatabase interface {
	CreateReplaceSession(string, string) error
	LookupSession(string) (bool, error)
	RegisterUser(string, string) error
	DeleteSession(string) error
	CheckValidUserCredentials(string, string) (bool, error)
}

var (
	AuthDb AuthDatabase
)

func ValidateAuthorization(session string) (bool, error) {
	ok, err := AuthDb.LookupSession(session)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func RegisterUser(username string, password string) error {
	err := AuthDb.RegisterUser(username, password)
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
