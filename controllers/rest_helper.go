package controllers

import (
	"encoding/json"
	"github.com/sparkscience/logjam/models"
	"net/http"
)

type RestResponseHelper struct {
}

func (r *RestResponseHelper) Write(rw http.ResponseWriter, response any, statusCode int) error {
	if strVersion, isStr := response.(string); isStr {

		rw.Header().Add("Content-Type", "text/html")
		rw.WriteHeader(statusCode)
		_, _ = rw.Write([]byte(strVersion))
	} else {
		bytes, err := json.Marshal(response)
		if err != nil {
			return err
		}
		rw.Header().Add("Content-Type", "application/json")
		_, _ = rw.Write(bytes)
	}

	return nil
}

func (r *RestResponseHelper) HandleIfErr(rw http.ResponseWriter, err error, status int) bool {
	if err == nil {
		return false
	}
	_ = r.Write(rw, models.MessageResponse{Message: err.Error()}, status)
	return true
}
