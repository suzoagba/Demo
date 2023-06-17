package handlers

import (
	"net/http"
	"regexp"
)

// Check if the header "Accept" contains text/html
func validateRequest(h http.Header) bool {
	reg := regexp.MustCompile(`(?m)text\/html`)
	for _, header := range h["Accept"] {
		match := reg.Match([]byte(header))
		if match {
			return true
		}
	}
	return false
}
