package database

import (
	"fmt"
)

func (d Database) deleteCurrentSession(username string) error {
	query, err := d.getQuery("delete-current-session")
	if err != nil {
		return err
	}
	query = fmt.Sprintf(query, username)
	err = d.Conn.Exec(query)
	return err
}

func (d Database) checkExistingSession(username string) (bool, error) {
	query, err := d.getQuery("check-existing-session")
	if err != nil {
		return false, err
	}

	session := ""
	query = fmt.Sprintf(query, username)
	if err := d.Conn.RunQuery(query, &session); err != nil {
		return false, err
	}

	return session != "", nil
}

func (d Database) checkUserCredentials(username, secret string) (bool, error) {
	query, err := d.getQuery("check-user-credentials")
	if err != nil {
		return false, err
	}

	dbSecret := ""
	query = fmt.Sprintf(query, username)
	if err := d.Conn.RunQuery(query, &secret); err != nil {
		return false, err
	}

	if dbSecret == "" || dbSecret != secret {
		return false, ErrUnauthorized
	}

	return true, nil
}

func (d Database) LookupSession(session string) (bool, error) {
	query, err := d.getQuery("lookup-session")
	if err != nil {
		return false, err
	}

	token := ""
	query = fmt.Sprintf(query, session)
	fmt.Println("query:", query)
	if err := d.Conn.RunQuery(query, &token); err != nil {
		fmt.Println("LookupSession:", err)
		return false, err
	}

	return token != "", err
}

func (d Database) RegisterUser(username, secret string) error {
	query, err := d.getQuery("register-user")
	if err != nil {
		return err
	}

	query = fmt.Sprintf(query, username, secret)
	return d.Conn.RunQuery(query)
}



func (d Database) CreateSession(username string, secret string) (string, error) {
	if ok, err := d.checkUserCredentials(username, secret); err != nil || !ok {
		return "", err
	}

	_, err := d.checkExistingSession(username)
	if err != nil {
		return "", err
	}
	/* TODO bring this back
	if ok {
		if err := d.deleteCurrentSession(username); err != nil {
			return "", err
		}
	}
	*/

	token := "crypticToken" // TODO change for valid token
	query, err := d.getQuery("create-session")
	if err != nil {
		return "", err
	}
	query = fmt.Sprintf(query, username, token)
	err = d.Conn.Exec(query)
	return token, err
}

