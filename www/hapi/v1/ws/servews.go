package verboten

import (
	"encoding/json"

	"github.com/gorilla/websocket"

	"github.com/reiver/logjam/models"
	"github.com/reiver/logjam/lib/rooms"
	"github.com/reiver/logjam/srv/room"
)


func serveWS(wsConn *websocket.Conn, socketId uint64, roomId string) {
	for {
		messageType, data, readErr := wsConn.ReadMessage()
		if readErr != nil {
			log.Error(readErr)
			go roomsrv.Controller.OnDisconnect(&rooms.WSContext{
				RoomId:        roomId,
				SocketID:      socketId,
				PureMessage:   nil,
				ParsedMessage: nil,
			})
			_ = wsConn.CloseHandler()(1001, readErr.Error())
			break
		}
		if messageType != websocket.TextMessage {
			log.Debugf("ignoring a message of type: %d", messageType)
			continue
		}

		var msg models.MessageContract
		err := json.Unmarshal(data, &msg)
		if err != nil {
			log.Error(err)
			continue
		}

		ctx := &rooms.WSContext{
			RoomId:        roomId,
			SocketID:      socketId,
			PureMessage:   data,
			ParsedMessage: &msg,
		}
		go handleEvent(ctx)
	}
}
