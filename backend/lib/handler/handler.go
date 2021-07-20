package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/sparkscience/logjam/backend/lib/message"
	logsrv "github.com/sparkscience/logjam/backend/srv/log"
)

type httpHandler struct {
}

type mySocket struct {
	Socket        *websocket.Conn
	Id            uint64
	IsBroadcaster bool
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
var Handler httpHandler
var ConnectedSockets map[*websocket.Conn]mySocket = make(map[*websocket.Conn]mySocket)
var ConnectedSocketsIndex uint64 = 0

// func parseMessage(socket mySocket, theMessage message.MessageContract) {
// 	// if theMessage.Type ==
// }

func reader(conn *websocket.Conn) {
	log := logsrv.Begin()
	defer log.End()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if err.Error() == "websocket: close 1005 (no status)" {
				delete(ConnectedSockets, conn)
				log.Inform("Socket closed!")
				return
			}
			log.Error("Read from socket error : ", err.Error())
			continue
		}
		log.Inform("Read from socket : ", string(p))
		var dat message.MessageContract
		if err := json.Unmarshal(p, &dat); err != nil {
			log.Error("JSON Unmarshal Error : ", err)
			continue
		}
		log.Inform("Data from socket : ", dat)
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Error("Error writing to socket : ", err)
		}
	}
}

func (receiver httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log := logsrv.Begin()
	defer log.End()
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Error("Upgrade Error : ", err)
	}
	log.Trace("Client Connected")
	ConnectedSockets[ws] = mySocket{
		Socket:        ws,
		Id:            ConnectedSocketsIndex,
		IsBroadcaster: false,
	}
	ConnectedSocketsIndex++
	reader(ws)
}
