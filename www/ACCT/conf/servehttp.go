package verboten

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/reiver/go-actsock"
	"github.com/reiver/go-fediverseid"
	"github.com/reiver/go-http400"
	"github.com/reiver/go-http500"
	libpath "github.com/reiver/go-path"

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
		http500.InternalServerError(responsewriter, request)
		log.Error("nil request")
		return
	}

	responsewriter.Header().Add("Access-Control-Allow-Origin", "*")

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

		if "" == account {
			http400.BadRequest(responsewriter, request)
			log.Debugf("HTTP Bad Request — acct == %q", account)
			return
		}

		log.Debugf("account (fediverse-id): %q", account)
	}

	switch {
	case websocket.IsWebSocketUpgrade(request):
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

		roomID := account
		go serveWS(wsConn, socketID, roomID)

	default:
		serveHTTP(responsewriter, request, account)
		return
	}
}

func serveHTTP(responsewriter http.ResponseWriter, request *http.Request, account string) {
	if nil == responsewriter {
		log.Error("nil response-writer")
		return
	}
	if nil == request {
		http500.InternalServerError(responsewriter, request)
		log.Error("nil request")
		return
	}
	if nil == request.URL {
		http500.InternalServerError(responsewriter, request)
		log.Error("nil request-url")
		return
	}

	var acctURI string
	{
		fediverseID, err := fediverseid.ParseFediverseIDString(account)
		if nil != err {
			http400.BadRequest(responsewriter, request)
			log.Debugf("HTTP Bad Request — bad account — acct == %q: %s", account, err)
			return
		}

		acctURI = fediverseID.AcctURI()
	}

	var id string
	var inoutbox string
	{
		var uri = *request.URL
		uri.User = nil
		uri.Scheme = "https"
		uri.Host = request.Host

		id = uri.String()

		uri.Scheme  = "wss"
		uri.Path = libpath.Canonical(uri.Path)
		inoutbox = uri.String()
	}

	var object = actsock.Conference{
		Actor: acctURI,
		EndPoints: map[string]string{
			"inoutbox":inoutbox,
		},
		ID: id,
		Name: fmt.Sprintf("%s — GreatApe", account),
	}

	responsewriter.Header().Set("Content-Type", "application/activity+json")
	io.WriteString(responsewriter, object.String())
}
