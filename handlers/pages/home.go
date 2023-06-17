package pages

import (
	"database/sql"
	"forum/handlers"
	"forum/structs"
	"net/http"
	"strconv"
)

func HomeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagIDStr := r.URL.Query().Get("id")
		tagID, err := strconv.Atoi(tagIDStr)
		if err != nil {
			// Set tagID to 0 to retrieve all posts
			tagID = 0
		}

		forPage := structs.ForPage{}
		forPage.User = handlers.IsLoggedIn(r, db).User
		forPage.Posts, _ = handlers.GetAllPosts(db, tagID, forPage.User)
		forPage.Tags = structs.Tags
		all := structs.Tag{
			ID:   0,
			Name: "All",
		}
		forPage.Tags = append([]structs.Tag{all}, forPage.Tags...)

		handlers.RenderTemplates("homepage", forPage, w, r)
	}
}
