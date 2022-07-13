package room

import (
	"github.com/sparkscience/logjam/backend/lib/handlers/roomhandler"
	httproutersrv "github.com/sparkscience/logjam/backend/srv/http/router"
	logsrv "github.com/sparkscience/logjam/backend/srv/log"
)

const Path = "/rooms"
const Method = "GET"

func init() {
	log := logsrv.Begin()
	defer log.End()

	err := httproutersrv.Router.DelegatePath(roomhandler.Handler(logsrv.Logger), Path, Method)
	if nil != err {
		log.Errorf("Could not register http.Handler for %q %q: %v", Method, Path, err)
		panic(err)
	}
}
