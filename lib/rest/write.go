package rest

import (
	"io"
	"net/http"

	"github.com/reiver/go-erorr"
	"github.com/reiver/go-json"
)

const (
	errNilResponseWriter = erorr.Error("nil http-response-writer")
)

const serverName  string = "Server"
const serverValue string = "logjam"

func Write(rw http.ResponseWriter, response interface{}, statusCode int) error {
	if response == nil {
		response = struct{}{}
	}

	switch casted := response.(type) {
	case string:
		writeHTML(rw, casted, statusCode)
	default:
		writeJSON(rw, response, statusCode)
	}

	return nil
}

func writeHTML(rw http.ResponseWriter, html string, statusCode int) error {
	if nil == rw {
		return errNilResponseWriter
	}

	rw.Header().Add("Content-Type", "text/html")
	rw.Header().Add(serverName, serverValue)

	rw.WriteHeader(statusCode)

	_, err := io.WriteString(rw, html)
	if nil != err {
		return err
	}

	return nil
}

func writeJSON(rw http.ResponseWriter, response any, statusCode int) error {
	if nil == rw {
		return errNilResponseWriter
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.Header().Add(serverName, serverValue)

	rw.WriteHeader(statusCode)

	bytes, err := json.Marshal(response)
	if nil != err {
		return err
	}

	_, err = rw.Write(bytes)
	if nil != err {
		return err
	}

	return nil
}
