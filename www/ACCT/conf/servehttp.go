package verboten

import (
	"net/http"

	"github.com/reiver/go-http400"

	"github.com/reiver/logjam/srv/http"
	"github.com/reiver/logjam/srv/websock"
)

const acct string = "acct"

const path string = "/{"+acct+"}/conf"

func init() {
	httpsrv.Router.HandleFunc(path, ServeHTTP).Methods(http.MethodGet, http.MethodOptions)
}

func ServeHTTP(responsewriter http.ResponseWriter, request *http.Request) {
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

	responsewriter.Header().Add("Access-Control-Allow-Origin", "*")

	wsConn, err := upgrader.Upgrade(responsewriter, request, nil)
	if nil != err {
		log.Errorf("problem upgrading to websocket: %s", err)
		return
	}
	socketID, err := websocksrv.WebSockSrv.OnConnect(wsConn)
	if nil != err {
		log.Errorf("problem on-connecting websocket: %s", err)
		_ = wsConn.Close()
		return
	}

	var account string
	{
		vars := httpsrv.Vars(request)
		if 0 < len(vars) {
			account = vars[acct]
		}

		// For backwards compatibility reasons.
		// Can remove later.
		if "" == account {
			account = request.URL.Query().Get("room")
		}
	}

	if "" == account {
		http400.BadRequest(responsewriter, request)
		log.Debugf("HTTP Bad Request — acct == %q", account)
		return
	}

	log.Debugf("account (fediverse-id): %q", account)
	roomID := account
	go serveWS(wsConn, socketID, roomID)
}
