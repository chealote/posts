package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"posts/internal/auth"
)

var (
	IgnoreAuthPaths = []string{"/signup", "/signin"}
)

type Handler struct {
	Mux *http.ServeMux
}

func replyError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path:", r.URL.Path)
	for _, path := range IgnoreAuthPaths {
		if path == r.URL.Path {
			h.Mux.ServeHTTP(w, r)
			return
		}
	}

	if len(r.Header["Authorization"]) <= 0 {
		replyError(w, http.StatusUnauthorized)
		return
	}
	token := r.Header["Authorization"][0]
	ok, err := auth.ValidateAuthorization(token)
	if err != nil {
		replyError(w, http.StatusInternalServerError)
		return
	}
	if !ok {
		replyError(w, http.StatusUnauthorized)
		return
	}
	h.Mux.ServeHTTP(w, r)
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	user := struct {
		Name   string `json:"name"`
		Secret string `json:"secret"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println("HandleSignUp:", err)
		replyError(w, http.StatusBadRequest)
		return
	}
	fmt.Println("User:", user.Name, user.Secret)
	if err := auth.RegisterUser(user.Name, user.Secret); err != nil {
		fmt.Println("HandleSignUp:", err)
		replyError(w, http.StatusBadRequest)
		return
	}
	fmt.Println("HandleSignUp: success")
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	user := struct {
		Name   string `json:"name"`
		Secret string `json:"secret"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println("HandleSignIn:", err)
		replyError(w, http.StatusBadRequest)
		return
	}
	token, err := auth.Login(user.Name, user.Secret)
	if err != nil {
		fmt.Println("HandleSignIn:", err)
		replyError(w, http.StatusUnauthorized)
		return
	}

	resToken := struct{
		Token string `json:"token"`
	}{
		Token: token,
	}
	res, err := json.Marshal(resToken)
	if err != nil {
		fmt.Println("HandleSignIn:", err)
		replyError(w, http.StatusInternalServerError)
		return
	}
	w.Write(res)
	fmt.Println("HandleSignIn: success")
}
