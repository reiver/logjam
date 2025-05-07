package verboten

import (
	"encoding/json"
	rtcroomsrv "github.com/reiver/logjam/srv/rtc-rooms"

	"github.com/gorilla/websocket"

	"github.com/reiver/logjam/lib/msgs"
	"github.com/reiver/logjam/lib/rtc-rooms"
)

func serveWS(wsConn *websocket.Conn, socketId uint64, roomId string) {
	for {
		messageType, data, readErr := wsConn.ReadMessage()
		if readErr != nil {
			log.Error(readErr)
			go rtcroomsrv.Controller.OnDisconnect(&rtc_rooms.WSContext{
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

		var msg msgs.MessageContract
		err := json.Unmarshal(data, &msg)
		if err != nil {
			log.Error(err)
			continue
		}

		ctx := &rtc_rooms.WSContext{
			RoomId:        roomId,
			SocketID:      socketId,
			PureMessage:   data,
			ParsedMessage: &msg,
		}
		go handleEvent(ctx)
	}
}
