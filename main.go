package main

import (
	"database/sql"
	"fmt"
	"forum/database"
	"forum/handlers"
	"forum/handlers/pages"
	"forum/handlers/user"
	"forum/structs"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	// Create new empty database
	database.StartDB(false)

	// Open a connection to the database
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the users table if it doesn't exist
	err = database.CreateTables(db)
	if err != nil {
		log.Fatal(err)
	}
	structs.GetTags(db)

	// Page handlers
	http.HandleFunc("/", pages.HomeHandler(db))
	http.HandleFunc("/createPost", pages.CreatePostHandler(db))
	http.HandleFunc("/viewPost", pages.ViewPostHandler(db))
	http.HandleFunc("/reply", handlers.ReplyHandler(db))

	http.HandleFunc("/register", user.RegisterHandler(db))
	http.HandleFunc("/login", user.LoginHandler(db))
	http.HandleFunc("/logout", user.LogoutHandler)

	// Handling assets
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/favicon.ico", ignoreFavicon)

	fmt.Println("Go to: http://localhost:8080")
	handlers.Open("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ignoreFavicon(_ http.ResponseWriter, _ *http.Request) {}
