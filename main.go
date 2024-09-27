package main

import (
	"flag"
	"net/http"
	"os"
	"fmt"
	"encoding/json"
	"posts/internal/auth"
	"posts/internal/database"
	"posts/internal/database/sqlite"
	"posts/internal/handler"
)

type Config struct {
	SQLite sqlite.Config `json:"sqlite"`
	Headers map[string]string `json:"headers"`
}

var (
	initDbFlag = flag.Bool("i", false, "initialize DB and exit")
	configFilepath = flag.String("c", "config.json", "config variables")
	config = Config{}
)

func init() {
	flag.Parse()
}

func main() {
	contents, err := os.ReadFile(*configFilepath)
	if err != nil {
		fmt.Println("ERROR: failed opening config file:", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(contents, &config); err != nil {
		fmt.Println("ERROR: failed parsing config file:", err)
		os.Exit(1)
	}

	conn, err := sqlite.Connect(config.SQLite)
	if err != nil {
		fmt.Println("ERROR: failed init():", err)
		return
	}
	database.Database = conn

	if *initDbFlag {
		database.Database.Initialize()
		os.Exit(0)
	}

	// TODO where should DB be?
	auth.DB = database.Database

	m := http.NewServeMux()
	h := handler.Handler{
		Mux: m,
		Headers: config.Headers,
	}
	m.HandleFunc("/", handler.HandleRoot)
	http.ListenAndServe(":8080", h)
}
