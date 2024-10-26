package auth

var (
	DB Database
)

type Database interface {
	ListPosts() ([]string, error)
}

func ListPostTitles() ([]string, error) {
	return DB.ListPosts()
}
