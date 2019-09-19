package handlers

import "net/http"

// HTTPError sends a HTTP Response with the given statusCode to `w`, and describes the status code.
func HTTPError(w http.ResponseWriter, statusCode int) {
	// Clear the header to avoid leaking any unintended info to the client
	h := w.Header()
	for k := range h {
		delete(h, k)
	}

	// Ensure the client always checks to see if the error was fixed
	w.Header().Set("Cache-Control", "No-Cache")

	// Send the error
	http.Error(w, http.StatusText(statusCode), statusCode)
}
