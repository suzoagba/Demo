package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func StartDB(input bool) {
	if input {
		err := createDatabase()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createDatabase() error {
	// Create the database folder if it doesn't exist
	err := os.MkdirAll("./database", os.ModePerm)
	if err != nil {
		return err
	}

	// Create the database file
	_, err = os.Create("./database/forum.db")
	if err != nil {
		return err
	}

	return nil
}

func CreateTables(db *sql.DB) error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
		uuid TEXT PRIMARY KEY,
		email TEXT UNIQUE,
		username TEXT UNIQUE,
		password TEXT,
		level INTEGER DEFAULT 0
	);

		CREATE TABLE IF NOT EXISTS posts (
			postID INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT,
			title TEXT,
			description TEXT,
			likes INTEGER DEFAULT 0,
			dislikes INTEGER DEFAULT 0,
			creationDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			imageFilename TEXT,
			FOREIGN KEY (username) REFERENCES users(username)
		);

		CREATE TABLE IF NOT EXISTS comments (
			commentID INTEGER PRIMARY KEY AUTOINCREMENT,
			postID INTEGER,
			username TEXT,
			content TEXT,
			likes INTEGER DEFAULT 0,
			dislikes INTEGER DEFAULT 0,
			creationDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (postID) REFERENCES posts(postID),
			FOREIGN KEY (username) REFERENCES users(username)
		);

		CREATE TABLE IF NOT EXISTS authenticated_users (
			session_id TEXT PRIMARY KEY,
			username TEXT UNIQUE,
			FOREIGN KEY(username) REFERENCES users(username)
		);

		CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE
		);

		CREATE TABLE IF NOT EXISTS post_tags (
			postID INTEGER,
			tagID INTEGER,
			FOREIGN KEY (postID) REFERENCES posts(postID),
			FOREIGN KEY (tagID) REFERENCES tags(id),
			PRIMARY KEY (postID, tagID)
		);

		INSERT OR IGNORE INTO tags (name) VALUES ('Cooking'), ('Mechanics'), ('Travel'), ('IT'), ('Random'), ('Market');

		CREATE TABLE IF NOT EXISTS post_reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			user_id UUID,
			reaction_type BOOLEAN,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(post_id, user_id),
			FOREIGN KEY (post_id) REFERENCES posts(postID),
			FOREIGN KEY (user_id) REFERENCES users(uuid)
		);
		
		CREATE TABLE IF NOT EXISTS comment_reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			comment_id INTEGER,
			user_id UUID,
			reaction_type BOOLEAN,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(comment_id, user_id),
			FOREIGN KEY (comment_id) REFERENCES comments(commentID),
			FOREIGN KEY (user_id) REFERENCES users(uuid)
		);	`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	fmt.Println("Database and tables created successfully!")
	return nil
}
