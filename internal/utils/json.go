package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Items interface{} `json:"items"`
	Total int `json:"total"`
	Page int `json:"page"`
	PageSize int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

// Send JSON as message
func WriteJSON(w http.ResponseWriter, status int, message string, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := APIResponse{
		Message: message,
		Data: data,
	}

	return json.NewEncoder(w).Encode(resp)
}

// Read JSON from request body
func ReadJSON(r *http.Request, data any) error {
	return json.NewDecoder(r.Body).Decode(data)
}

// Handling error
func WriteError(w http.ResponseWriter, err error, customMsg ...string) error {
	var message string

	var status int
	var errorMsg string

	// Check if its from Error we defined
	if appErr, ok := err.(*Error); ok {
		status = appErr.StatusCode
		errorMsg = appErr.Message
		message = appErr.Message
	} else {
		status = http.StatusInternalServerError
		errorMsg = "Server internal error"
		message = "Server error"
	}

	// If customMsg is specifically given and not empty, override
    if len(customMsg) > 0 && customMsg[0] != "" {
        message = customMsg[0]
    }

	resp := map[string]any {
		"error": errorMsg,
	}

	return WriteJSON(w, status, message, resp)
}