package user

import (
	"database/sql"
	"fmt"
	"forum/handlers"
	"forum/structs"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		forPage := structs.ForPage{}
		forPage.User = handlers.IsLoggedIn(r, db).User
		forPage.Error.Error = false
		if r.Method == http.MethodGet {
			if forPage.User.LoggedIn {
				http.Redirect(w, r, "/", http.StatusFound)
			} else {
				handlers.RenderTemplates("login", forPage, w, r)
			}
		} else if r.Method == http.MethodPost {
			login(w, r, db, forPage)
		}
	}
}

func login(w http.ResponseWriter, r *http.Request, db *sql.DB, forPage structs.ForPage) {
	err := r.ParseForm()
	if err != nil {
		forPage.Error.Error = true
		forPage.Error.Message = "Failed to parse form data"
		handlers.RenderTemplates("login", forPage, w, r)
		return
	}

	// Extract the user input values
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if passwordCorrect("SELECT password FROM users WHERE username = ?", username, password, db) {
		sessionID := uuid.New().String()
		expiration := time.Now().Add(24 * time.Hour) // Set the expiration time for the cookie
		createCookie(w, "forum-session", sessionID, expiration)
		err := updateSessionID(username, sessionID, db)
		if err != nil {
			forPage.Error.Error = true
			forPage.Error.Message = "Unable to update session ID: " + err.Error() + "."
			handlers.RenderTemplates("login", forPage, w, r)
			return
		}
		err = addActiveSession(db, sessionID, username)
		if err != nil {
			forPage.Error.Error = true
			forPage.Error.Message = "Unable to add active session"
			handlers.RenderTemplates("login", forPage, w, r)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		forPage.Error.Error = true
		forPage.Error.Message = "Check your username/password"
		handlers.RenderTemplates("login", forPage, w, r)
		return
	}

}

func updateSessionID(uuid, sessionID string, db *sql.DB) error {
	// Prepare the SQL statement
	fmt.Println(uuid, sessionID)
	stmt, err := db.Prepare("UPDATE authenticated_users SET session_id = ? WHERE username = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(sessionID, uuid)
	if err != nil {
		return err
	}

	return nil
}

func passwordCorrect(query string, value string, password string, db *sql.DB) bool {
	row := db.QueryRow(query, value)

	var storedPassword string
	// Scan the retrieved password into the variable
	err := row.Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found in the database
			return false
		}
		return false
	}

	// Compare the stored password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		// Password does not match
		return false
	}

	return true
}

// Create a new session cookie
func createCookie(w http.ResponseWriter, name, value string, expiration time.Time) {
	cookie := &http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expiration,
	}

	http.SetCookie(w, cookie)
}

// Add new active session to the database
func addActiveSession(db *sql.DB, sessionID, username string) error {
	insertSQL := `
        INSERT OR REPLACE INTO authenticated_users (session_id, username)
        VALUES (?, ?);
    `

	_, err := db.Exec(insertSQL, sessionID, username)
	if err != nil {
		return err
	}

	return nil
}
