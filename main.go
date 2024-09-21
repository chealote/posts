package main

import (
	"fmt"
	"posts/internal/database/sqlite"
	"posts/internal/database"
)

/* TODO also bring this back
var (
	initDbFlag = flag.Bool("i", false, "initialize DB and exit")
	configFilepath = flag.String("c", "config.json", "config variables")
	cfg = sqlite.Cfg{}
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

	if *initDbFlag {
		conn.Initialize()
		os.Exit(0)
	}

	// TODO where should DB be?
	auth.DB = conn

	m := http.NewServeMux()
	h := handler.Handler{Mux: m}
	m.HandleFunc("/", handler.HandleRoot)
	m.HandleFunc("/signup", handler.HandleSignUp)
	m.HandleFunc("/signin", handler.HandleSignIn)
	http.ListenAndServe(":8080", h)
}
*/

func main() {
	// instance database
	// sqlite injection
	config := sqlite.Config{
		Username: "",
		Password: "",
	}

	sqlite, err := sqlite.Connect(config)
	if err != nil {
		panic(err)
	}

	db := database.Database{
		Conn: sqlite,
		ScriptsPath: "./scripts",
	}

	fmt.Println(db.RegisterUser("ale", "secret"))
	fmt.Println(db.CreateSession("ale", "secret"))
}
