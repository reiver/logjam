package verboten

import (
	"net/http"

	"github.com/reiver/logjam/srv/http"
	"github.com/reiver/logjam/srv/websock"
)

const path string = "/hapi/v1/ws"

func init() {
	httpsrv.Router.HandleFunc(path, serveHTTP)
}

func serveHTTP(responsewriter http.ResponseWriter, request *http.Request) {
	if nil == responsewriter {
		log.Error("nil response-writer")
		return
	}
	if nil == request {
		const code int = http.StatusInternalServerError
		http.Error(responsewriter, http.StatusText(code), code)
		log.Error("nil request")
		return
	}

	wsConn, err := upgrader.Upgrade(responsewriter, request, nil)
	if err != nil {
		log.Errorf("problem upgrading to websocket: %s", err)
		return
	}
	socketId, err := websocksrv.WebSockSrv.OnConnect(wsConn)
	if err != nil {
		log.Error(err)
		_ = wsConn.Close()
		return
	}

	roomId := request.URL.Query().Get("room")
	go serveWS(wsConn, socketId, roomId)
}
