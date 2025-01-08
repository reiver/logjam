package verboten

import (
	"net/http"

	"github.com/reiver/logjam/srv/http"
	"github.com/reiver/logjam/srv/websock"
)

const path string = "/ws"

func init() {
	httpsrv.Router.HandleFunc(path, serveHTTP)
}

func serveHTTP(writer http.ResponseWriter, request *http.Request) {
	if nil == writer {
		return
	}
	if nil == request {
		const code int = http.StatusInternalServerError
		http.Error(writer, http.StatusText(code), code)
		return
	}

	wsConn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Error(err)
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
