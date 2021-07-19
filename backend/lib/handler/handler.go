package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"

	logsrv "github.com/sparkscience/logjam/backend/srv/log"
)

type httpHandler struct {
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
var Handler httpHandler
var ConnectedSockets map[int64]*websocket.Conn = make(map[int64]*websocket.Conn)
var ConnectedSocketsIndex int64 = 0

func reader(conn *websocket.Conn) {
	log := logsrv.Begin()
	defer log.End()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if err.Error() == "websocket: close 1005 (no status)" {
				for key, value := range ConnectedSockets {
					if value == conn {
						log.Inform("Removed socket index : ", key)
						delete(ConnectedSockets, key)
					}
				}
				log.Inform("Socket closed!")
				return
			}
			log.Error("Read from socket error : ", err.Error())
			continue
		}
		log.Inform("Read from socket : ", string(p))
		var dat map[string]interface{}
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
	ConnectedSockets[ConnectedSocketsIndex] = ws
	ConnectedSocketsIndex++
	reader(ws)
}
