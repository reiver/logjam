package main

import (
	"github.com/sparkscience/logjam/backend/arg"
	"github.com/sparkscience/logjam/backend/srv/log"

	"fmt"
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
}
