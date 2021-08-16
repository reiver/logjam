package commandhandler

import (
	"net/http"

	logger "github.com/mmcomp/go-log"
	"github.com/sparkscience/logjam/backend/lib/websocketmap"
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
	websocketmap.Map.Reset()
	w.Write([]byte("Command '" + req.URL.Path + "' executed"))
}
