package verboten

import (
	"net/http"

	"github.com/reiver/go-jsonld"
	libpath "github.com/reiver/go-path"

	"github.com/reiver/logjam/lib/rooms"
	"github.com/reiver/logjam/srv/http"
	"github.com/reiver/logjam/srv/room"
)

const path string = "/hapi/v1/rooms"

func init() {
	httpsrv.Router.HandleFunc(path, serveGET).Methods(http.MethodGet, http.MethodOptions)
}

func serveGET(responsewriter http.ResponseWriter, request *http.Request) {
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

	type Item struct {
		Type string `json:"type"`
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	var items []Item
	{
		err := roomsrv.Repository.ForEachRoom(func(roomModel *rooms.RoomModel)error{
			var id   string = rooms.RoomURL(roomModel.ID, path, request.Host)
			var name string = roomModel.ID

			var item = Item{
				Type:"Link",
				ID:   id,
				Name: name,
			}

			items = append(items, item)
			return nil
		})
		if nil != err {
			const code int = http.StatusInternalServerError
			http.Error(responsewriter, http.StatusText(code), code)
			log.Errorf("problem iterating through rooms: %s", err)
			return
		}
	}
	log.Debugf("number-of-items: %d", len(items))

	var numRooms int = roomsrv.Repository.NumRooms()
	log.Debugf("number-of-rooms: %d", numRooms)

	var bytes []byte
	{

		var response = struct{
			NameSpace jsonld.NameSpace `jsonld:"https://www.w3.org/ns/activitystreams"`

			Type       string `json:"type"`
			TotalItems int    `json:"totalItems"`
			Items    []Item   `json:"items"`
		}{
			Type:    "Collection",
			TotalItems: numRooms,
			Items: items,
		}

		var err error
		bytes, err = jsonld.Marshal(response)
		if nil != err {
			log.Errorf("problem encoding response as JSON-LD: %s", err)
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
