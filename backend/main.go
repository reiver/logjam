package main

import (
	"github.com/sparkscience/logjam/backend/arg"
	logsrv "github.com/sparkscience/logjam/backend/srv/log"

	"fmt"
	"net/http"
)

type httpHandler struct {
}

func (receiver httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log := logsrv.Begin()
	defer log.End()
	log.Trace("Received a request!")
	w.Write([]byte("Hello World"))
	log.Trace("Finished output to client!")
}

// func httpHandler(w http.ResponseWriter, req *http.Request) {
// 	log := logsrv.Begin()
// 	defer log.End()
// 	log.Trace("Received a request!")
// 	w.Write([]byte("Hello World"))
// 	log.Trace("Finished output to client!")
// }

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

	{
		var httpAddr string = fmt.Sprintf(":%d", arg.HttpPort)
		log.Log("HTTP address:", httpAddr)

		// http.HandleFunc("/", httpHandler)
		var handler http.Handler = httpHandler{}

		log.Log("starting HTTP server")

		err := http.ListenAndServe(httpAddr, handler)
		if nil != err {
			log.Error("problem with HTTP server:", err)
			return
		}

		log.Log("HTTP stopped")
	}
}
