package auth

var (
	DB Database
)

type Database interface {
	CreateSession(string, string) (string, error)
	LookupSession(string) (bool, error)
	RegisterUser(string, string) error
}

func ValidateAuthorization(session string) (bool, error) {
	ok, err := DB.LookupSession(session)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func RegisterUser(username string, secret string) error {
	// TODO check if user exists first? in the same query?
	return DB.RegisterUser(username, secret)
}

func Login(username string, secret string) (string, error) {
	return DB.CreateSession(username, secret)
}
