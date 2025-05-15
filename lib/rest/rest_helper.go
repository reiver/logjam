package rest

import (
	"net/http"
)

func HandleIfErr(rw http.ResponseWriter, err error, status int) bool {
	if err == nil {
		return false
	}

	Error(rw, err, status)
	return true
}
