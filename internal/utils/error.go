package utils

import (
    "errors"
    "fmt"
    "net/http"
)

type Error struct {
    Err        error
    Message    string
    StatusCode int
}

func (e *Error) Error() string {
    if e.Err != nil {
        return e.Err.Error()
    }
    return e.Message
}

func (e *Error) Unwrap() error {
    return e.Err 
}

func NewNotFoundError(resource string) *Error {
    return &Error{
        Message:    fmt.Sprintf("%s not found", resource),
        StatusCode: http.StatusNotFound, // 404
    }
}

func NewBadRequestError(message string) *Error {
    return &Error{
        Message:    message,
        StatusCode: http.StatusBadRequest, // 400
    }
}

func NewInternalServerError(err error) *Error {
    return &Error{
        Err:        err,
        Message:    "Internal server error",
        StatusCode: http.StatusInternalServerError, // 500
    }
}

func NewConflictError(message string) *Error {
    return &Error{
        Message:    message,
        StatusCode: http.StatusConflict, // 409
    }
}

func IsNotFound(err error) bool {
    var e *Error
    return errors.As(err, &e) && e.StatusCode == http.StatusNotFound
}

func IsBadRequest(err error) bool {
    var e *Error
    return errors.As(err, &e) && e.StatusCode == http.StatusBadRequest
}

func IsConflict(err error) bool {
    var e *Error
    return errors.As(err, &e) && e.StatusCode == http.StatusConflict
}
