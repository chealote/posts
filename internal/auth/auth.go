package auth

import (
	"errors"
	"fmt"
	"posts/internal/database"
	"posts/internal/utils"
	"time"
)

var (
	DB Database
)

type Database interface {
	CreateReplaceSession(string, string) error
	LookupSession(string) (bool, error)
	RegisterUser(string, string) error
	DeleteSession(string) error
	CheckValidUserCredentials(string, string) (bool, error)
}

func ValidateAuthorization(session string) (bool, error) {
	ok, err := DB.LookupSession(session)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func RegisterUser(username string, password string) error {
	err := DB.RegisterUser(username, password)
	// TODO this is a debug check only, just return error without checking
	if errors.Is(err, database.ErrConstraintKey) {
		fmt.Println("user already exists!")
	}
	return err
}

func Login(username string, password string) (string, error) {
	ok, err := DB.CheckValidUserCredentials(username, password)
	if err != nil || !ok {
		// TODO return 500 or 401 if !ok
		return "", err
	}

	epochNow := fmt.Sprintf("%d", time.Now().Unix())
	token := utils.Sha512Sum(fmt.Sprintf("%s%s", username, epochNow))
	return token, DB.CreateReplaceSession(username, token)
}

func Logout(token string) error {
	return DB.DeleteSession(token)
}
