package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func RenderTemplates(page string, data interface{}, w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("PANIC:", err)
			ErrorHandler(w, http.StatusInternalServerError, "")
		}
	}()
	if !validateRequest(r.Header) {
		ErrorHandler(w, http.StatusBadRequest, "")
		return
	}
	fmt.Println("#PAGE: " + page)

	link := "./templates/"
	switch page {
	case "homepage":
		link += "home.html"
	case "register":
		link += "authentication/register.html"
	case "login":
		link += "authentication/login.html"
	case "createPost":
		link += "createPost.html"
	case "viewPost":
		link += "viewPost.html"
	default:
		ErrorHandler(w, http.StatusNotFound, "")
		return
	}

	templates := template.Must(template.ParseFiles("./templates/base.html", link))
	err := templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
