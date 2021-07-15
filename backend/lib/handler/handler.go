package handler

import (
	"net/http"

	logsrv "github.com/sparkscience/logjam/backend/srv/log"
)

type httpHandler struct {
}

var Handler httpHandler

func (receiver httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log := logsrv.Begin()
	defer log.End()
	log.Trace("Receive a new Request!")
	w.Write([]byte("Hello World"))
	log.Trace("Finished sending `Hello World`!")
}
