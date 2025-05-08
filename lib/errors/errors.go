package cErrors

import (
	"net/http"
	"strconv"
)

// CError custom error
type CError struct {
	Message string `json:"message"`
	Code    int    `json:"-"`
}

func (e CError) Error() string {
	return e.Message
}

func NewError(code int) CError {
	msg := http.StatusText(code)
	if len(msg) == 0 {
		msg = "unknown error happened, contact support. code: " + strconv.Itoa(code)
	}
	return CError{
		Message: msg,
		Code:    code,
	}
}

func NewErrorWithMsg(code int, msg string) CError {
	return CError{
		Message: msg,
		Code:    code,
	}
}

func NewErrorFromErr(code int, err error) CError {
	return CError{
		Message: err.Error(),
		Code:    code,
	}
}
