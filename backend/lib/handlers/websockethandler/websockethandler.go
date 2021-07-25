package websockethandler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"

	"github.com/sparkscience/logjam/backend/lib/message"
	"github.com/sparkscience/logjam/backend/lib/websocketmap"
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
var webSocketMaps websocketmap.WebSocketMapType = websocketmap.WebSocketMapType{
	Connections: make(map[*websocket.Conn]websocketmap.MySocket),
}

func decideWhomToConnect(broadcaster websocketmap.MySocket) websocketmap.MySocket {
	if len(broadcaster.ConnectedSockets) < 2 {
		return broadcaster
	}

	var toCheckSocket websocketmap.MySocket = broadcaster
	for subSocket := range toCheckSocket.ConnectedSockets {
		theSubSocket := webSocketMaps.Connections[subSocket]
		if len(theSubSocket.ConnectedSockets) < 2 {
			return theSubSocket
		}
		toCheckSocket = theSubSocket
	}

	return toCheckSocket
}

func parseMessage(socket websocketmap.MySocket, messageJSON []byte, messageType int) {
	log := logsrv.Begin()
	defer log.End()

	var theMessage message.MessageContract
	{
		err := json.Unmarshal(messageJSON, &theMessage)
		if err != nil {
			log.Error("Error unmarshal message ", err)
		}
		log.Highlight("TheMessage  Type : ", theMessage.Type, " Data : ", theMessage.Data, " Target : ", theMessage.Target)
	}

	var response message.MessageContract
	switch theMessage.Type {
	case "start":
		response.Type = "start"
		response.Data = strconv.FormatInt(int64(socket.ID), 10)
		responseJSON, err := json.Marshal(response)
		if err == nil {
			socket.Socket.WriteMessage(messageType, responseJSON)
		} else {
			log.Error("Marshal Error of `start` response", err)
		}
	case "role":
		response.Type = "role"
		if theMessage.Data == "broadcast" {
			log.Highlight("New Broadcaster : ", socket.ID)
			response.Data = "yes:broadcast"
			webSocketMaps.RemoveBroadcasters()
			webSocketMaps.SetBroadcaster(socket.Socket)
		} else {
			log.Highlight("New Audiance : ", socket.ID)
			if socket.IsBroadcaster {
				response.Data = "no:broadcast"
				webSocketMaps.RemoveBroadcaster(socket.Socket)
			} else {
				response.Data = "yes:audience"
				broadcaster, ok := webSocketMaps.GetBroadcaster()
				if !ok {
					response.Data = "no:audience"
				} else {
					var broadResponse message.MessageContract
					broadResponse.Type = "add_audience"
					broadResponse.Data = strconv.FormatInt(int64(socket.ID), 10)
					broadResponseJSON, err := json.Marshal(broadResponse)
					if err == nil {
						// broadcaster.Socket.WriteMessage(messageType, broadResponseJSON)
						targetSocket := decideWhomToConnect(broadcaster)
						if targetSocket.Socket != broadcaster.Socket {
							broadResponse.Type = "add_broadcast_audience"
						}
						webSocketMaps.InsertConnected(targetSocket.Socket, socket.Socket)
						targetSocket.Socket.WriteMessage(messageType, broadResponseJSON)
					} else {
						log.Error("Marshal Error of `add_audience` broadResponse", err)
					}
				}
			}
		}
		responseJSON, err := json.Marshal(response)
		if err == nil {
			socket.Socket.WriteMessage(messageType, responseJSON)
		} else {
			log.Error("Marshal Error of `role` response", err)
		}
	default:
		ID, err := strconv.ParseUint(theMessage.Target, 10, 64)
		if err != nil {
			log.Error("Inavlid Target : ", theMessage.Target)
			return
		}
		target, ok := webSocketMaps.GetSocketByID(ID)
		if ok {
			log.Inform("Default sending to ", ID, " ", string(messageJSON))
			target.Socket.WriteMessage(messageType, messageJSON)
		}
	}
}

func reader(conn *websocket.Conn) {
	log := logsrv.Begin()
	defer log.End()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				webSocketMaps.Delete(conn)
				log.Inform("Socket closed!")
				return
			}
			log.Error("Read from socket error : ", err.Error())
			continue
		}
		log.Inform("Read from socket : ", string(p))
		parseMessage(webSocketMaps.Connections[conn], p, messageType)
	}
}

func (receiver httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log := logsrv.Begin()
	defer log.End()
	log.Highlight("Method", req.Method)
	log.Highlight("Path", req.URL.Path)
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Error("Upgrade Error : ", err)
	}
	log.Trace("Client Connected")
	webSocketMaps.Insert(ws)
	reader(ws)
}
