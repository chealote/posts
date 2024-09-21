package main

import (
	"flag"
	"net/http"
	"os"
	"encoding/json"
	"posts/internal/auth"
	"posts/internal/posts"
	"posts/internal/database"
	"posts/internal/database/sqlite"
	"posts/internal/handler"
)

var (
	initDbFlag = flag.Bool("i", false, "initialize DB and exit")
	configFilepath = flag.String("c", "config.json", "config variables")
	cfg = database.Cfg{}
)

func init() {
	flag.Parse()

	contents, err := os.ReadFile(*configFilepath)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(contents, &cfg); err != nil {
		panic(err)
	}
}

func main() {
	conn, err := sqlite.Connect(cfg)
	if err != nil {
		panic(err)
	}
	database.Database = conn

	if *initDbFlag {
		database.Database.Initialize()
		os.Exit(0)
	}

	// TODO where should DB be?
	auth.DB = database.Database
	posts.DB = database.Database

	m := http.NewServeMux()
	h := handler.Handler{Mux: m}
	m.HandleFunc("/", handler.HandleRoot)
	m.HandleFunc("/signup", handler.HandleSignUp)
	m.HandleFunc("/signin", handler.HandleSignIn)
	http.ListenAndServe(":8080", h)
}
