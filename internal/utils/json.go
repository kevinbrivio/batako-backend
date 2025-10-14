package utils

import (
	"encoding/json"
	"net/http"
)

// Send JSON as message
func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// Read JSON from request body
func ReadJSON(r *http.Request, data any) error {
	return json.NewDecoder(r.Body).Decode(data)
}

// Handling error
func WriteError(w http.ResponseWriter, err error) error {
	// Check if its from Error we defined
	if appErr, ok := err.(*Error); ok {
		return WriteJSON(w, appErr.StatusCode, map[string]string{
			"error": appErr.Message,
		})
	}

	// Default to 500
	return WriteJSON(w, http.StatusInternalServerError, map[string]string{
		"error": "Internal server error",
	})
}