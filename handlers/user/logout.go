package user

import (
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	deleteCookie(w, "forum-session")
	http.Redirect(w, r, "/", http.StatusFound)
}

// Delete session cookie
func deleteCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:    name,
		Value:   "",
		Expires: time.Unix(0, 0), // Set the expiration time to a past time to delete the cookie
	}

	http.SetCookie(w, cookie)
}
