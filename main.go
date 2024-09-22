package main

import (
	"flag"
	"net/http"
	"os"
	"encoding/json"
	"posts/internal/auth"
	"posts/internal/database"
	"posts/internal/database/sqlite"
	"posts/internal/handler"
)

var (
	initDbFlag = flag.Bool("i", false, "initialize DB and exit")
	configFilepath = flag.String("c", "config.json", "config variables")
	config = database.Config{}
)

func init() {
	flag.Parse()

	contents, err := os.ReadFile(*configFilepath)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(contents, &config); err != nil {
		panic(err)
	}
}

func main() {
	conn, err := sqlite.Connect(config)
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

	m := http.NewServeMux()
	h := handler.Handler{Mux: m}
	m.HandleFunc("/", handler.HandleRoot)
	m.HandleFunc("/signup", handler.HandleSignUp)
	m.HandleFunc("/signin", handler.HandleSignIn)
	http.ListenAndServe(":8080", h)
}
