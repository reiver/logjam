package websockets

import (
	"github.com/sparkscience/logjam/backend/lib/handlers/websockethandler"
	httproutersrv "github.com/sparkscience/logjam/backend/srv/http/router"
	logsrv "github.com/sparkscience/logjam/backend/srv/log"
	"github.com/sparkscience/logjam/backend/srv/websocket"
)

const Path = "/ws"
const Method = "GET"

func init() {
	log := logsrv.Begin()
	defer log.End()

	wsHandlers := websockethandler.NewWSHandlers()

	wsHttpHandler := websocketsrv.Handler(logsrv.Logger, []string{"example.com"}, wsHandlers.OnConnected, wsHandlers.OnDisconnected)
	wsHttpHandler.RegisterWSHandler("start", wsHandlers.Start)
	wsHttpHandler.RegisterWSHandler("role", wsHandlers.Role)
	wsHttpHandler.RegisterWSHandler("stream", wsHandlers.Stream)
	wsHttpHandler.RegisterWSHandler("turn_status", wsHandlers.TurnStatus)
	wsHttpHandler.RegisterWSHandler("log", wsHandlers.Log)
	wsHttpHandler.RegisterWSHandler("ping", wsHandlers.Ping)
	wsHttpHandler.RegisterWSHandler("tree", wsHandlers.Tree)
	wsHttpHandler.RegisterWSHandler("broadcaster-status", wsHandlers.BroadcasterStatus)
	wsHttpHandler.RegisterWSHandler("*", wsHandlers.DefaultHandler)

	err := httproutersrv.Router.Register(wsHttpHandler, Path, Method)
	if nil != err {
		log.Errorf("Could not register http.Handler for %q %q: %v", Method, Path, err)
		panic(err)
	}
}
