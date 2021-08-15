package commandhandler

import (
	"net/http"

	logger "github.com/mmcomp/go-log"
)

type httpHandler struct {
	Logger logger.Logger
}

func Handler(logger logger.Logger) http.Handler {
	return httpHandler{
		Logger: logger,
	}
}

func (receiver httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log := receiver.Logger.Begin()
	defer log.End()

	log.Inform("Command : ", req.URL.Path)
	if req.URL.Path == "/command/reset" {
		w.Write([]byte("Command '" + req.URL.Path + "' executed"))
		return
	}
	w.Write([]byte("Command '" + req.URL.Path + "' not found!"))
}
