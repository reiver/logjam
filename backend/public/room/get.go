package websockets

import (
	"github.com/sparkscience/logjam/backend/lib/handlers/roomwebsockethandler"
	httproutersrv "github.com/sparkscience/logjam/backend/srv/http/router"
	logsrv "github.com/sparkscience/logjam/backend/srv/log"
)

const Path = "/room"
const Method = "GET"

func init() {
	log := logsrv.Begin()
	defer log.End()

	err := httproutersrv.Router.Register(roomwebsockethandler.Handler(logsrv.Logger), Path, Method)
	if nil != err {
		log.Errorf("Could not register http.Handler for %q %q: %v", Method, Path, err)
		panic(err)
	}
}
