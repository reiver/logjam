package verboten

import (
	"encoding/json"
	"net/http"

	"github.com/reiver/go-erorr"

	"github.com/reiver/logjam/lib/rooms"
	"github.com/reiver/logjam/srv/http"
	"github.com/reiver/logjam/srv/room"
)

const roomIDname  string = "room-id"
const pathpattern string = "/hapi/v1/rooms/{"+roomIDname+"}"

func init() {
	httpsrv.Router.HandleFunc(pathpattern, serveHTTP)
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

	if nil == roomsrv.Repository {
		const code int = http.StatusInternalServerError
		http.Error(responsewriter, http.StatusText(code), code)
		log.Error("nil room-service")
		return
	}

	var roomID string = "----TODO----"
	{
		var m map[string]string = httpsrv.Vars(request)
		if len(m) <= 0 {
			const code int = http.StatusInternalServerError
			http.Error(responsewriter, http.StatusText(code), code)
			log.Error("nil request-vars")
			return
		}

		var found bool
		roomID, found = m[roomIDname]
		if !found {
			const code int = http.StatusInternalServerError
			http.Error(responsewriter, http.StatusText(code), code)
			log.Error("room-id missing")
			return
		}
	}

	var room *rooms.RoomModel
	{
		var err error

		room, err = roomsrv.Repository.GetRoom(roomID)
		if nil != err {
			switch {
			case erorr.Is(err, rooms.ErrRoomNotFound):
				const code int = http.StatusNotFound
				http.Error(responsewriter, http.StatusText(code), code)
				log.Debugf("room %q not found (ErrRoomNotFound)", roomID)
				return
			default:
				const code int = http.StatusInternalServerError
				http.Error(responsewriter, http.StatusText(code), code)
				log.Errorf("problem getting room: %s", err)
				return
			}
		}
		if nil == room {
			const code int = http.StatusNotFound
			http.Error(responsewriter, http.StatusText(code), code)
			log.Debugf("room %q not found (nil)", roomID)
			return
		}
	}

	{
		var response = struct{
			Type string `json:"type"`
			ID   string `json:"id"`
			Name string `json:"name"`
		}{
			Type: "Room",
			ID:   roomID,
			Name: room.Title,
		}

		err := json.NewEncoder(responsewriter).Encode(response)
		if nil != err {
			log.Errorf("problem encoding rsponse as JSON: %s", err)
			return
		}
	}
}
