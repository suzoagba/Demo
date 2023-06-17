package user

import (
	"database/sql"
	"fmt"
	"forum/handlers"
	"forum/structs"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		forPage := structs.ForPage{}
		forPage.User = handlers.IsLoggedIn(r, db).User
		forPage.Error.Error = false
		if r.Method == http.MethodGet {
			if forPage.User.LoggedIn {
				http.Redirect(w, r, "/", http.StatusFound)
			} else {
				handlers.RenderTemplates("register", forPage, w, r)
			}
		} else if r.Method == http.MethodPost {
			register(w, r, db, forPage)
		}
	}
}

func register(w http.ResponseWriter, r *http.Request, db *sql.DB, forPage structs.ForPage) {
	// Parse the form data from the request
	err := r.ParseForm()
	if err != nil {
		forPage.Error.Error = true
		forPage.Error.Message = "Failed to parse form data"
		handlers.RenderTemplates("register", forPage, w, r)
		return
	}

	// Extract the user input values
	email := r.Form.Get("email")
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	invalidInput := false // Checks if the credentials that user wrote are valid

	// Check if the email or username is already taken
	if rowExists("SELECT email from users WHERE email = ?", email, db) { // if the email exists
		invalidInput = true
		forPage.Error.Error = true
		forPage.Error.Message = "Email already taken"
		forPage.Error.Field1 = email
		forPage.Error.Field2 = username
		handlers.RenderTemplates("register", forPage, w, r)
		return
	} else if rowExists("SELECT username from users WHERE username = ?", username, db) { // if the email exists
		invalidInput = true
		forPage.Error.Error = true
		forPage.Error.Message = "Username already taken"
		forPage.Error.Field1 = email
		forPage.Error.Field2 = username
		handlers.RenderTemplates("register", forPage, w, r)
		return
	}

	// Encrypt the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to encrypt password", http.StatusInternalServerError)
		return
	}

	if !invalidInput {
		// Insert the user into the database
		err = insertUser(email, username, string(hashedPassword), db)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		fmt.Println("- Registration failed!")
		fmt.Fprintln(w, "Registration failed!")
	}
}

// Function to check if the email is already taken (example implementation)
func rowExists(query string, value string, db *sql.DB) bool {
	row := db.QueryRow(query, value)
	switch err := row.Scan(&value); err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		return false
	}
}

// Function to insert the user into the database (example implementation)
func insertUser(email, username, password string, db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO users (uuid, username, password, email) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(uuid.New().String(), username, password, email)
	return nil
}
