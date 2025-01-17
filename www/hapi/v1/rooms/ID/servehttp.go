package verboten

import (
	"net/http"

	"github.com/reiver/go-erorr"
	"github.com/reiver/go-jsonld"
	libpath "github.com/reiver/go-path"

	"github.com/reiver/logjam/lib/rooms"
	"github.com/reiver/logjam/srv/http"
	"github.com/reiver/logjam/srv/room"
)

const roomIDname  string = "room-id"
const pathpattern string = "/hapi/v1/rooms/{"+roomIDname+"}"

func init() {
	httpsrv.Router.HandleFunc(pathpattern, serveHTTP).Methods(http.MethodGet, http.MethodOptions)
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

	var roomID string
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
	log.Debugf("room-id = %q", roomID)

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
				log.Errorf("problem getting room %q: %s", roomID, err)
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

	var bytes []byte
	{
		var path string = libpath.Parent(pathpattern)

		var name string = roomID
		var id string = rooms.RoomURL(roomID, path, request.Host)

		var response = struct{
			NameSpace jsonld.NameSpace `jsonld:"https://www.w3.org/ns/activitystreams"`

			Type string `json:"type"`
			ID   string `json:"id"`
			Name string `json:"name"`
		}{
			Type: "Object",
			ID:   id,
			Name: name,
		}

		var err error
		bytes, err = jsonld.Marshal(response)
		if nil != err {
			log.Errorf("problem encoding response for room %q as JSON: %s", roomID, err)
			return
		}
	}

	{
		responsewriter.Header().Add("Access-Control-Allow-Origin", "*")
		responsewriter.Header().Add("Content-Type", "application/activity+json")

		_, err := responsewriter.Write(bytes)
		if nil != err {
			log.Errorf("problem sending response: %s", err)
			return
		}
	}
}
