package rest

import (
	"net/http"
)

func Error(rw http.ResponseWriter, err error, status int) {
	rw.WriteHeader(status)
	_ = Write(rw, MessageResponse{Message: err.Error()}, status)
}
