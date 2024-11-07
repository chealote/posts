package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"posts/internal/auth"
	"regexp"
	"strings"
)

type UserCredentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Handler struct {
	Mux     *http.ServeMux
	Headers map[string]string
}

type PostWithId struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type PostContent struct {
	Title string `json:"title"`
	Contents string `json:"contents"`
}

type PostDatabase interface {
	CreatePost(string, string, string) error
	ListWithId() ([]PostWithId, error)
	ContentsPost(string) (PostContent, error)
}

var (
	PostDb PostDatabase

	ignoreAuthFromPaths = []string{"/signup", "/signin"}
	rePostIdValidChars  = regexp.MustCompile("[^a-zA-Z0-9]+")
)

func replyError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Path:", r.URL.Path, "Method:", r.Method)

	for header, value := range h.Headers {
		w.Header().Set(header, value)
	}

	for _, path := range ignoreAuthFromPaths {
		if path == r.URL.Path {
			h.Mux.ServeHTTP(w, r)
			return
		}
	}

	if r.Method == "OPTIONS" {
		return
	}

	if len(r.Header["Authorization"]) <= 0 {
		fmt.Printf("ERROR: ServeHTTP: missing auth header\n")
		replyError(w, http.StatusUnauthorized)
		return
	}

	token := r.Header["Authorization"][0]

	ok, err := auth.ValidateAuthorization(token)
	if err != nil {
		fmt.Printf("ERROR: ServeHTTP: %s\n", err)
		replyError(w, http.StatusInternalServerError)
		return
	}

	if !ok {
		fmt.Printf("ERROR: ServeHTTP: token unauthorized: %s\n", token)
		replyError(w, http.StatusUnauthorized)
		return
	}

	h.Mux.ServeHTTP(w, r)
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		fmt.Println("In HandleRoot, skipping OPTIONS method")
		return
	}

	pieces := strings.Split(r.URL.Path, "/")
	if len(pieces) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}

	switch pieces[1] {
	case "signin":
		HandleSignIn(w, r)
	case "signup":
		HandleSignUp(w, r)
	case "token":
		w.Write([]byte("seems valid"))
	case "logout":
		HandleLogout(w, r)
	case "posts":
		if len(pieces) > 2 {
			HandlePostWithId(w, r, pieces[2])
		} else {
			HandlePostRoot(w, r)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
	}
}

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	user := UserCredentials{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println("ERROR: HandleSignUp:", err)
		replyError(w, http.StatusBadRequest)
		return
	}

	fmt.Println("User:", user.Name, user.Password)
	if err := auth.RegisterUser(user.Name, user.Password); err != nil {
		fmt.Println("ERROR: HandleSignUp:", err)
		replyError(w, http.StatusBadRequest)
		return
	}
	fmt.Println("HandleSignUp: success")
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	user := UserCredentials{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println("ERROR: HandleSignIn reading body:", err)
		replyError(w, http.StatusBadRequest)
		return
	}

	token, err := auth.Login(user.Name, user.Password)
	if err != nil {
		fmt.Println("ERROR: HandleSignIn from auth.Login():", err)
		replyError(w, http.StatusUnauthorized)
		return
	}

	w.Write([]byte(token))
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("attempting to logout"))

	// handler already checked for the auth token, should be safe to use here?
	token := r.Header["Authorization"][0]
	if err := auth.Logout(token); err != nil {
		fmt.Println("ERROR: HandleLogout:", err)
		replyError(w, http.StatusInternalServerError)
		return
	}
}

func HandlePostRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		post := struct {
			Id    string `json:"id"`
			Title string `json:"title"`
			Post  string `json:"post"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			fmt.Println("ERROR: HandleSignIn Decode JSON:", err)
			replyError(w, http.StatusBadRequest)
			return
		}

		post.Id = rePostIdValidChars.ReplaceAllString(post.Title, "-")

		if err := PostDb.CreatePost(post.Id, post.Title, post.Post); err != nil {
			fmt.Println("ERROR: HandleSignIn CreatePost:", err)
			replyError(w, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	case "GET":
		postWithIds, err := PostDb.ListWithId()
		if err != nil {
			replyError(w, http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(postWithIds)
		if err != nil {
			replyError(w, http.StatusInternalServerError)
			return
		}

		w.Write(b)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
	}
}

func HandlePostWithId(w http.ResponseWriter, r *http.Request, postId string) {
	switch r.Method {
	case "GET":
		post, err := PostDb.ContentsPost(postId)
		if err != nil {
			replyError(w, http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(post)
		if err != nil {
			replyError(w, http.StatusInternalServerError)
			return
		}

		w.Write(b)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))

	}
}
