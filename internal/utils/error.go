package utils

import (
	"fmt"
	"net/http"
)

type Error struct {
	Err error
	Message string
	StatusCode int
}

func (e * Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func NewNotFoundError(resource string) *Error {
	return &Error{
		Message: fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound, // 404
	}
}

func NewBadRequestError(message string) *Error {
	return &Error{
		Message: message,
		StatusCode: http.StatusBadRequest, // 400
	}
}

func NewInternalServerError(err error) *Error {
	return &Error{
		Err: err,
		Message: "Internal server error",
		StatusCode: http.StatusInternalServerError, // 500
	}
}

func NewConflictError(message string) *Error {
    return &Error{
        Message:    message,
        StatusCode: http.StatusConflict, // 409
    }
}