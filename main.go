package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"posts/internal/auth"
	"posts/internal/database/sqlite"
	"posts/internal/handler"
)

type Config struct {
	SQLite  sqlite.Config     `json:"sqlite"`
	Headers map[string]string `json:"headers"`
}

const (
	PORT = ":8080"
)

var (
	initDbFlag     = flag.Bool("i", false, "initialize DB and exit")
	configFilepath = flag.String("c", "config.json", "config variables")
	config         = Config{}
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

	// TODO this is a special case when databases can be Open() and
	// Close(); maybe the database doesn't need that or something
	conn, err := sqlite.Connect(config.SQLite)
	if err != nil {
		fmt.Println("ERROR: failed init():", err)
		return
	}
	defer conn.Close()

	if *initDbFlag {
		conn.Initialize()
		os.Exit(0)
	}

	auth.AuthDb = conn
	handler.PostDb = conn

	m := http.NewServeMux()
	h := handler.Handler{
		Mux:     m,
		Headers: config.Headers,
	}
	m.HandleFunc("/", handler.HandleRoot)

	fmt.Println("Listening on port", PORT)
	panic(http.ListenAndServe(PORT, h))
}
