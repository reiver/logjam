package verboten

import (
	"encoding/json"

	"github.com/gorilla/websocket"

	"codeberg.org/greatape/logjam/lib/msgs"
	"codeberg.org/greatape/logjam/lib/rooms"
	"codeberg.org/greatape/logjam/srv/room"
)


func serveWS(wsConn *websocket.Conn, socketID uint64, roomID string) {
	for {
		messageType, data, err := wsConn.ReadMessage()
		if nil != err {
			log.Error(err)
			go roomsrv.Controller.OnDisconnect(&rooms.WSContext{
				RoomId:        roomID,
				SocketID:      socketID,
				PureMessage:   nil,
				ParsedMessage: nil,
			})
			_ = wsConn.CloseHandler()(1001, err.Error())
			break
		}
		if messageType != websocket.TextMessage {
			log.Debugf("ignoring websocket message of type: %d", messageType)
			continue
		}

		var msg msgs.Message
		err = json.Unmarshal(data, &msg)
		if nil != err {
			log.Errorf("problem json-unmarshaling data from websocket as %T: %s", msg, err)
			continue
		}

		ctx := &rooms.WSContext{
			RoomId:        roomID,
			SocketID:      socketID,
			PureMessage:   data,
			ParsedMessage: &msg,
		}
		go handleEvent(ctx)
	}
}
