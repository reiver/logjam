package websockethandler

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/mmcomp/go-log"
	"github.com/sparkscience/logjam/backend/lib/message"
	binarytreesrv "github.com/sparkscience/logjam/backend/srv/binarytree"
	roommapssrv "github.com/sparkscience/logjam/backend/srv/roommaps"
	"net/http"
)

type handlers struct {
}

func NewWSHandlers() *handlers {
	return &handlers{}
}

func (h *handlers) OnConnected(ws *websocket.Conn, roomName string) (err error, statusCode int) {
	var mapFound bool
	_, mapFound = roommapssrv.RoomMaps.Get(roomName)
	if !mapFound {
		AMap := binarytreesrv.GetMap()
		{
			err := roommapssrv.RoomMaps.Set(roomName, &AMap)
			if err != nil {
				log.Errorf("could not set room in map: %s", err)
				return errors.New("internal Server Error"), http.StatusInternalServerError
			}
		}
	}
	Map, found := roommapssrv.RoomMaps.Get(roomName)
	if !found {
		log.Errorf("could not get map for room %q when trying to serve-http", roomName)
		return errors.New("internal Server Error"), http.StatusInternalServerError
	}
	Map.Insert(ws)
	{
		err := roommapssrv.RoomMaps.Set(roomName, Map)
		if err != nil {
			log.Errorf("could not set room in map: %s", err)
			return errors.New("internal Server Error"), http.StatusInternalServerError
		}
	}

	return nil, 101
}

func (h *handlers) OnDisconnected(ws *websocket.Conn) {

}

func (h *handlers) Start(ws *binarytreesrv.MySocket, msg message.MessageContract, msgType int) {

}

func (h *handlers) Role(ws *binarytreesrv.MySocket, msg message.MessageContract, msgType int) {

}

func (h *handlers) Stream(ws *binarytreesrv.MySocket, msg message.MessageContract, msgType int) {

}

func (h *handlers) TurnStatus(ws *binarytreesrv.MySocket, msg message.MessageContract, msgType int) {

}

func (h *handlers) Log(ws *binarytreesrv.MySocket, msg message.MessageContract, msgType int) {

}

func (h *handlers) Ping(ws *binarytreesrv.MySocket, msg message.MessageContract, msgType int) {

}

func (h *handlers) Tree(ws *binarytreesrv.MySocket, msg message.MessageContract, msgType int) {

}

func (h *handlers) BroadcasterStatus(ws *binarytreesrv.MySocket, msg message.MessageContract, msgType int) {

}

func (h *handlers) DefaultHandler(ws *binarytreesrv.MySocket, msg message.MessageContract, msgType int) {

}
