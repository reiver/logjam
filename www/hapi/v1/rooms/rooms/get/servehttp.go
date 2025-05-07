package verboten

import (
	"github.com/reiver/logjam/lib/rest"
	httpsrv "github.com/reiver/logjam/srv/http"
	roomsrv "github.com/reiver/logjam/srv/room"
	"net/http"
)

const path string = "/hapi/v1/rooms"

func init() {
	httpsrv.Router.HandleFunc(path, serveHTTP).Methods(http.MethodGet, http.MethodOptions)
}

func serveHTTP(responsewriter http.ResponseWriter, request *http.Request) {
	if nil == responsewriter {
		return
	}
	if nil == request {
		const code int = http.StatusInternalServerError
		http.Error(responsewriter, http.StatusText(code), code)
		return
	}
	uid := request.URL.Query().Get("roomUID")
	room, err := roomsrv.Repository.GetRoom(uid)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	if room == nil {
		_ = rest.Write(responsewriter, nil, http.StatusNotFound)
	} else {
		_ = rest.Write(responsewriter, room, http.StatusOK)
	}
}
