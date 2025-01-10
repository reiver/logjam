package rest

import (
	"encoding/json"
	"net/http"
)

func Write(rw http.ResponseWriter, response interface{}, statusCode int) error {
	if response == nil {
		response = struct{}{}
	}
	if strVersion, isStr := response.(string); isStr {
		rw.Header().Add("Content-Type", "text/html")
		rw.WriteHeader(statusCode)
		_, _ = rw.Write([]byte(strVersion))
	} else {
		rw.WriteHeader(statusCode)
		bytes, err := json.Marshal(response)
		if err != nil {
			return err
		}
		rw.Header().Add("Content-Type", "application/json")
		_, _ = rw.Write(bytes)
	}

	return nil
}
