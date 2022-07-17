package roomhandler

import (
	"net/http"
	"strings"

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

	basePath := strings.Split(strings.Replace(req.URL.Path, "/rooms/", "", 1), "/")
	room := basePath[0]
	role := "audience"
	if len(basePath) == 2 && basePath[1] == "admin" {
		role = "broadcast"
	}
	log.Alert("room: " + room + " role: " + role)
	extraQuery := "&" + req.URL.Query().Encode()
	log.Alert("roomhandler.ServeHTTP ", room, " ", role)

	url := "/files/frontend/?room=" + room + "&role=" + role + extraQuery

	http.Redirect(w, req, url, http.StatusMovedPermanently)
}
