package posts

import (
	"text/template"
	"posts/internal/database"
	"os"
)

var (
	DB Database
	postSummaryTemplate = `
<b>{{ .Title }}</b>
<p>{{ .Summary }}</p>
`
)

type PostSummary struct {
	Title string
	Summary string
}

type Database interface {
	ListTopPosts() ([]database.Post, error)
}

func ListTopPosts() ([]byte, error) {
	// TODO it replies with links, those links should be downloaded at the frontend?
	// or maybe download links here and like run the template thing here
	posts, err := DB.ListTopPosts()
	if err != nil {
		return []byte{}, err
	}
	tmpl, err := template.New("post").Parse(postSummaryTemplate)
	if err != nil {
		return []byte{}, err
	}
	for _, post := range posts {
		if err := tmpl.Execute(os.Stdout, PostSummary{
			Title: post.Title,
			Summary: "This should download partial text from the file or something else, maybe this part could live in the DB while the complete thing is in a static file, along with images and videos",
		}); err != nil {
			return []byte{}, err
		}
	}
	return []byte{}, err
}
