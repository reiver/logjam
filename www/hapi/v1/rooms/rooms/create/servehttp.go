package verboten

import (
	"encoding/json"
	"github.com/reiver/logjam/lib/rest"
	"github.com/reiver/logjam/lib/room"
	httpsrv "github.com/reiver/logjam/srv/http"
	roomsrv "github.com/reiver/logjam/srv/room"
	"io"
	"net/http"
)

const path string = "/hapi/v1/rooms"

func init() {
	httpsrv.Router.HandleFunc(path, serveHTTP).Methods(http.MethodPost, http.MethodOptions)
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
	reqBody, err := io.ReadAll(request.Body)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}
	var req room.CreateRoomDTO
	err = json.Unmarshal(reqBody, &req)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}

	req.OwnerID = "" //read from token
	err = roomsrv.Repository.CreateRoom(req)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	_ = rest.Write(responsewriter, nil, http.StatusCreated)
}
