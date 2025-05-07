package verboten

import (
	"github.com/reiver/logjam/lib/rest"
	httpsrv "github.com/reiver/logjam/srv/http"
	roomsrv "github.com/reiver/logjam/srv/room"
	"net/http"
)

const path string = "/hapi/v1/rooms"

func init() {
	httpsrv.RouterWithAuth.HandleFunc(path, serveHTTP).Methods(http.MethodDelete, http.MethodOptions)
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
	user := rest.GetUser(request)
	if user == nil {
		http.Error(responsewriter, "not authenticated", http.StatusInternalServerError)
		return
	}
	err := roomsrv.Repository.DeleteRoom(request.URL.Query().Get("roomId"), user.ID)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	_ = rest.Write(responsewriter, nil, http.StatusNoContent)
}
