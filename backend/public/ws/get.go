package websockets

import (
	"github.com/sparkscience/logjam/backend/lib/handlers/websockethandler"
	httproutersrv "github.com/sparkscience/logjam/backend/srv/http/router"
)

const Path = "/ws"

func init() {
	httproutersrv.Router.Register(websockethandler.Handler, Path, "GET")
}
