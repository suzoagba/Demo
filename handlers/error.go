package handlers

import (
	"html/template"
	"net/http"
	"strconv"
)

func ErrorHandler(w http.ResponseWriter, status int, extra string) {
	w.WriteHeader(status)
	errorStruct := make(map[string]string)
	errorStruct["ErrorNumber"] = strconv.Itoa(status)
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		panic(err)
	}
	if status == http.StatusNotFound {
		errorStruct["ErrorText1"] = "Oops! Looks like you got lost."
		errorStruct["ErrorText2"] = "Page not found."
	} else if status == http.StatusBadRequest {
		errorStruct["ErrorText1"] = "Oops! Bad Request!"
		errorStruct["ErrorText2"] = "You can only request text/html."
	} else if status == http.StatusInternalServerError {
		errorStruct["ErrorText1"] = "Oops! This is awkward."
		errorStruct["ErrorText2"] = "Internal Server Error."
	} else {
		errorStruct["ErrorText1"] = "Are you feeling a bit naughty?"
		errorStruct["ErrorText2"] = http.StatusText(status)
	}
	errorStruct["ErrorText3"] = extra
	tmpl.Execute(w, errorStruct)
}
