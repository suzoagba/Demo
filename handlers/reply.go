package handlers

import (
	"database/sql"
	"net/http"
)

func ReplyHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ifIn := IsLoggedIn(r, db)
		if ifIn.User.LoggedIn {
			if r.Method == http.MethodPost {
				// Parse the form data from the request
				err := r.ParseForm()
				if err != nil {
					http.Error(w, "Failed to parse form data", http.StatusBadRequest)
					return
				}

				// Extract the reply data
				postID := r.Form.Get("postID")
				content := r.Form.Get("content")

				// Insert the reply data into the database
				_, err = db.Exec("INSERT INTO comments (postID, username, content) VALUES (?, ?, ?)", postID, ifIn.User.Username, content)
				if err != nil {
					http.Error(w, "Failed to insert reply data into database", http.StatusInternalServerError)
					return
				}

				// Redirect or display a success message
				http.Redirect(w, r, "/viewPost?id="+postID, http.StatusFound)
			}
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}
