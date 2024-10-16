package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"posts/internal/auth"
	"time"
)

type UserCredentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

var (
	IgnoreAuthFromPaths = []string{"/signup", "/signin"}
)

type Handler struct {
	Mux     *http.ServeMux
	Headers map[string]string
}

func replyError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Path:", r.URL.Path, "Method:", r.Method)

	for header, value := range h.Headers {
		w.Header().Set(header, value)
	}

	for _, path := range IgnoreAuthFromPaths {
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

	switch r.URL.Path {
	case "/signin":
		HandleSignIn(w, r)
	case "/signup":
		HandleSignUp(w, r)
	case "/token":
		w.Write([]byte("seems valid"))
	case "/logout":
		HandleLogout(w, r)
	default:
		// TODO pretending to read a DB so adding some sleepy time
		time.Sleep(time.Second)
		back := fmt.Sprintf("Hello World from path %s", r.URL.Path)
		w.Write([]byte(back))
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
		fmt.Println("ERROR: HandleSignIn:", err)
		replyError(w, http.StatusBadRequest)
		return
	}

	token, err := auth.Login(user.Name, user.Password)
	if err != nil {
		fmt.Println("ERROR: HandleSignIn:", err)
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
