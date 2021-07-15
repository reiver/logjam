package main

import (
	"github.com/sparkscience/logjam/backend/arg"
	"github.com/sparkscience/logjam/backend/srv/log"

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

	fmt.Println("Hello world! — I am logjam! ✨")

//@TODO
/*
	{
		var httpAddr string = fmt.Sprintf(":%d", arg.HttpPort)
		log.Log("HTTP address:", httpAddr)

		var handler http.Handler = ????

		log.Log("starting HTTP server")

		err := http.ListenAndServe(httpAddr, handler)
		if nil != err {
			log.Error("problem with HTTP server:", err)
/////////////////////// RETURN
			return
		}

		log.Log("HTTP stopped")
	}
*/
}
