package auth

var (
	DB Database
)

type Database interface {
	CreateSession(string, string) (string, error)
	LookupSession(string) (bool, error)
	RegisterUser(string, string) error
	DeleteSession(string) error
}

func ValidateAuthorization(session string) (bool, error) {
	ok, err := DB.LookupSession(session)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func RegisterUser(username string, password string) error {
	// TODO check if user exists first? in the same query?
	return DB.RegisterUser(username, password)
}

func Login(username string, password string) (string, error) {
	return DB.CreateSession(username, password)
}

func Logout(token string) error {
	return DB.DeleteSession(token)
}
