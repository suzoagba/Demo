package pages

import (
	"database/sql"
	"fmt"
	"forum/handlers"
	"forum/structs"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is logged in
		ifIn := handlers.IsLoggedIn(r, db)
		if !ifIn.User.LoggedIn {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		// Handle GET request to render the create post page
		if r.Method == http.MethodGet {
			handlers.RenderTemplates("createPost", structs.ForPage{
				User: ifIn.User,
				Tags: structs.Tags,
			}, w, r)
			return
		}

		// Handle POST request
		if r.Method != http.MethodPost {
			return
		}

		// Parse the form data from the request
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			handlers.ErrorHandler(w, http.StatusBadRequest, "Failed to parse form data")
			return
		}

		// Extract the post data
		title := r.FormValue("title")
		description := r.FormValue("description")
		selectedTags := r.Form["tags"]

		// Check for empty fields
		if title == "" || description == "" || len(selectedTags) == 0 {
			handlers.ErrorHandler(w, http.StatusBadRequest, "Forbidden empty fields.")
			return
		}

		// Check if an image file was uploaded
		file, header, err := r.FormFile("image")
		if err != nil {
			handlers.ErrorHandler(w, http.StatusBadRequest, "Failed to retrieve the image file.")
			return
		}
		defer file.Close()

		// Validate the image size
		if header.Size > 20*1024*1024 { // Max size: 20MB
			handlers.ErrorHandler(w, http.StatusBadRequest, "Image size exceeds the maximum limit of 20MB.")
			return
		}

		// Validate the image format
		fileExt := filepath.Ext(header.Filename)
		allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
		if !allowedExts[fileExt] {
			handlers.ErrorHandler(w, http.StatusBadRequest, "Invalid image format. Only JPEG, PNG, and GIF formats are allowed.")
			return
		}

		// Prepare the SQL statement for inserting post data
		stmt := "INSERT INTO posts (username, title, description, imageFilename) VALUES (?, ?, ?, ?)"

		// Generate a unique filename for the image
		imageFilename := generateUniqueFilename(header.Filename)

		// Execute the SQL statement to insert post data into the database
		result, err := db.Exec(stmt, ifIn.User.Username, title, description, imageFilename)
		if err != nil {
			handlers.ErrorHandler(w, http.StatusInternalServerError, "Failed to insert post data into the database.")
			return
		}

		// Get the ID of the newly created post
		postID, err := result.LastInsertId()
		if err != nil {
			handlers.ErrorHandler(w, http.StatusInternalServerError, "Failed to get post ID.")
			return
		}

		// Save the image file to the server
		imagePath := filepath.Join("./assets/uploads", imageFilename)
		err = saveImageToFile(file, imagePath)
		if err != nil {
			handlers.ErrorHandler(w, http.StatusInternalServerError, "Failed to save image.")
			return
		}

		// Insert the selected tags into the post_tags table
		for _, tagID := range selectedTags {
			_, err = db.Exec("INSERT INTO post_tags (postID, tagID) VALUES (?, ?)", postID, tagID)
			if err != nil {
				handlers.ErrorHandler(w, http.StatusInternalServerError, "Failed to insert tag into post_tags table.")
				return
			}
		}

		// Redirect to the homepage or display a success message
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func saveImageToFile(file multipart.File, imagePath string) error {
	// Create a new file at the specified path
	f, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Copy the uploaded file to the new file
	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	return nil
}

func generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	filename := fmt.Sprintf("%d%s", timestamp, ext)
	return filename
}
