package verboten

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/reiver/logjam/lib/rest"
	"github.com/reiver/logjam/models"
	"github.com/reiver/logjam/srv/goldgorilla"
	"github.com/reiver/logjam/srv/http"
	"github.com/reiver/logjam/srv/room"
	"github.com/reiver/logjam/srv/websock"
)

const path string = "/goldgorilla/rejoin"

func init() {
        httpsrv.Router.HandleFunc(path, serveHTTP)
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
	var reqModel struct {
		RoomId string `json:"roomId"`
		GGID   uint64 `json:"ggid"`
	}
	err = json.Unmarshal(reqBody, &reqModel)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}
	broadcaster, err := roomsrv.Repository.GetBroadcaster(reqModel.RoomId)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	if broadcaster == nil {
		_ = rest.Write(responsewriter, nil, 503)
		return
	}
	_, _, err = roomsrv.Repository.RemoveMember(reqModel.RoomId, reqModel.GGID)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}
	roomMembersIdList, err := roomsrv.Repository.GetAllMembersId(reqModel.RoomId, true)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	brDCEvent := models.MessageContract{
		Type: "event-broadcaster-disconnected",
		Data: strconv.FormatUint(broadcaster.ID, 10),
	}

	_ = websocksrv.WebSockSrv.Send(brDCEvent, roomMembersIdList...)

	go func(roomId string, membersIdList []uint64) {
		time.Sleep(500 * time.Millisecond)
		err := goldgorillasrv.Repository.Start(roomId)
		if err != nil {
			log.Error(err)
			return
		}
		brIsBackEvent := models.MessageContract{
			Type: "broadcasting",
		}

		_ = websocksrv.WebSockSrv.Send(brIsBackEvent, membersIdList...)
	}(reqModel.RoomId, roomMembersIdList)
	_ = rest.Write(responsewriter, nil, 204)
}
