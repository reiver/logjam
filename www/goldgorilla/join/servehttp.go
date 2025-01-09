package verboten

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/lib/goldgorilla"
	"github.com/reiver/logjam/lib/msgs"
	"github.com/reiver/logjam/lib/rest"
	"github.com/reiver/logjam/srv/goldgorilla"
	"github.com/reiver/logjam/srv/http"
	"github.com/reiver/logjam/srv/room"
	"github.com/reiver/logjam/srv/websock"
)

const path string = "/goldgorilla/join"

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
	var reqModel goldgorilla.JoinReqModel
	err = json.Unmarshal(reqBody, &reqModel)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}
	//ctrl.conf.GoldGorillaSVCAddr = reqModel.ServiceAddr
	newGGID := websocksrv.WebSockSrv.GetNewID()
	err = roomsrv.Repository.AddMember(reqModel.RoomId, newGGID, "{}", "", "", true)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}
	err = roomsrv.Repository.UpdateCanConnect(reqModel.RoomId, newGGID, true)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}
	parentId, err := roomsrv.Repository.InsertMemberToTree(reqModel.RoomId, newGGID, true)
	if rest.HandleIfErr(responsewriter, err, 500) {
		_, _, _ = roomsrv.Repository.RemoveMember(reqModel.RoomId, newGGID)
		return
	}
	err = goldgorillasrv.Repository.CreatePeer(reqModel.RoomId, *parentId, true, true, newGGID)
	if rest.HandleIfErr(responsewriter, err, 503) {
		return
	}
	_ = websocksrv.WebSockSrv.Send(msgs.MessageContract{
		Type: "add_audience",
		Data: strconv.FormatUint(newGGID, 10),
	}, *parentId)

	_ = rest.Write(responsewriter, struct {
		ID uint64 `json:"id"`
	}{
		ID: newGGID,
	}, 200)
	memsId, err := roomsrv.Repository.GetAllMembersId(reqModel.RoomId, false)
	if err != nil {
		log.Error(err)
	} else {
		_ = websocksrv.WebSockSrv.Send(msgs.MessageContract{Type: "goldgorilla-joined", Data: strconv.FormatUint(newGGID, 10)}, memsId...)
	}
	go func(roomId string, svcAddr string, ggId uint64) {
		for {
			res, err := http.Get(svcAddr + "/healthcheck?roomId=" + roomId)
			if err != nil {
				break
			}
			if res.StatusCode > 204 {
				break
			}
			time.Sleep(2 * time.Second)
		}
		_, childrenIdList, err := roomsrv.Repository.RemoveMember(roomId, newGGID)
		if err != nil {
			log.Error(err)
			return
		}
		parentDCEvent := msgs.MessageContract{
			Type: "event-parent-dc",
			Data: strconv.FormatUint(newGGID, 10),
		}
		_ = websocksrv.WebSockSrv.Send(parentDCEvent, childrenIdList...)
		log.Info("deleted a goldgorilla instance from tree")
	}(reqModel.RoomId, cfg.Config.GoldGorillaBaseURL(), newGGID)
}
