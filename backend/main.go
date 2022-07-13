package main

import (
	"github.com/sparkscience/logjam/backend/arg"
	"github.com/sparkscience/logjam/backend/public/graphapi"
	"github.com/sparkscience/logjam/backend/public/logfiles"
	"github.com/sparkscience/logjam/backend/public/room"
	roomwebsockets "github.com/sparkscience/logjam/backend/public/roomws"
	"github.com/sparkscience/logjam/backend/public/statics"
	websockets "github.com/sparkscience/logjam/backend/public/ws"
	httproutersrv "github.com/sparkscience/logjam/backend/srv/http/router"
	logsrv "github.com/sparkscience/logjam/backend/srv/log"

	"fmt"
	"net/http"
)

func main() {
	log := logsrv.Begin()
	defer log.End()

	// If --test-logs flag was used, then do a logs test and then exit.
	if arg.TestLogs {
		log.Trace("trace test")
		log.Log("log test")
		log.Inform("inform test")
		log.Highlight("highlight test")
		log.Warn("warn test")
		log.Alert("alert test")
		log.Error("error test")
		/////////////// RETURN
		return
	}

	log.Log("Hello world! — I am logjam! ✨")
	log.Log("Websocket listening on ", websockets.Path)
	log.Log("Static listening on ", statics.Path)
	log.Log("Log File listening on ", logfiles.Path)
	log.Log("Graphapi listening on ", graphapi.Path)
	log.Log("RoomWebsocket listening on ", roomwebsockets.Path)
	log.Log("Room listening on ", room.Path)

	{
		var httpAddr string = fmt.Sprintf(":%d", arg.HttpPort)
		log.Log("HTTP address:", httpAddr)

		log.Log("starting HTTP server")

		err := http.ListenAndServe(httpAddr, httproutersrv.Router)
		if nil != err {
			log.Error("problem with HTTP server:", err)
			return
		}

		log.Log("HTTP stopped")
	}
}
