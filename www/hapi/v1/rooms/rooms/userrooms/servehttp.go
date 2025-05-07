package verboten

import (
	"github.com/reiver/logjam/lib/rest"
	httpsrv "github.com/reiver/logjam/srv/http"
	roomsrv "github.com/reiver/logjam/srv/room"
	"net/http"
)

const path string = "/hapi/v1/rooms/me"

func init() {
	httpsrv.RouterWithAuth.HandleFunc(path, serveHTTP).Methods(http.MethodGet, http.MethodOptions)
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
	rooms, err := roomsrv.Repository.GetUserRooms(user.ID)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	_ = rest.Write(responsewriter, rooms, http.StatusOK)
}
