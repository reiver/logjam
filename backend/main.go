package main

import (
	"github.com/sparkscience/logjam/backend/arg"
	// "github.com/sparkscience/logjam/backend/public/statics"
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
	// log.Log("Static listening on ", statics.Path)

	{
		var httpAddr string = fmt.Sprintf(":%d", arg.HttpPort)
		log.Log("HTTP address:", httpAddr)

		log.Log("starting HTTP server")

		// Question : Charles this is not workoing :(
		fs := http.FileServer(http.Dir("./public/statics/files"))
  		http.Handle("/", fs)
		// ENd of Question
		
		err := http.ListenAndServe(httpAddr, httproutersrv.Router)
		if nil != err {
			log.Error("problem with HTTP server:", err)
			return
		}

		log.Log("HTTP stopped")
	}
}
