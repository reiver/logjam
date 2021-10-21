package roomwebsockethandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/sparkscience/logjam/backend/lib/message"
	"github.com/sparkscience/logjam/backend/lib/roomwebsocketmap"

	logger "github.com/mmcomp/go-log"
)

type httpHandler struct {
	Logger logger.Logger
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Handler(logger logger.Logger) http.Handler {
	return httpHandler{
		Logger: logger,
	}
}

func (receiver httpHandler) parseMessage(socket roomwebsocketmap.MySocket, messageJSON []byte, messageType int, userAgent string) {
	log := receiver.Logger.Begin()
	defer log.End()
	// f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	// if err == nil {
	// 	defer f.Close()
	// }

	var theMessage message.MessageContract
	{
		err := json.Unmarshal(messageJSON, &theMessage)
		if err != nil {
			log.Error("Error unmarshal message ", err)
			return
		}
		// log.Highlight("TheMessage  Type : ", theMessage.Type, " Data : ", theMessage.Data, " Target : ", theMessage.Target)
	}

	// var response message.MessageContract
	switch theMessage.Type {
	case "open-room":
		room := theMessage.Data
		var fullMessage map[string]interface{}
		err := json.Unmarshal(messageJSON, &fullMessage)
		if err != nil {
			log.Error("Error unmarshal message ", err)
			return
		}
		roomObject := roomwebsocketmap.Map.GetRoom(room)
		socketIds := ""
		for sock := range roomObject.Sockets {
			if socketIds != "" {
				socketIds += ","
			}
			socketIds += strconv.FormatUint(roomObject.Sockets[sock].ID, 10)
		}
		fullMessage["room"] = socketIds
		messageJSON, err = json.Marshal(fullMessage)
		if err != nil {
			log.Error("Error marshal message ", err)
			return
		}
		roomwebsocketmap.Map.AddToRoom(room, socket)
		log.Alert("SOKET ", socket)
		socket.Socket.WriteMessage(messageType, messageJSON)
	default:
		ID, err := strconv.ParseUint(theMessage.Target, 10, 64)
		if err != nil {
			log.Error("Inavlid Target : ", theMessage.Target)
			return
		}
		_, target, roomerr := roomwebsocketmap.Map.GetSocketByID(ID)
		log.Inform("target ", target)
		if roomerr == nil {
			// log.Inform("Default sending to ", ID, " ", string(messageJSON))
			var fullMessage map[string]interface{}
			err := json.Unmarshal(messageJSON, &fullMessage)
			if err != nil {
				log.Error("Error unmarshal message ", err)
				return
			}
			fullMessage["username"] = socket.Name
			messageJSON, err = json.Marshal(fullMessage)
			if err != nil {
				log.Error("Error marshal message ", err)
				return
			}
			target.Socket.WriteMessage(messageType, messageJSON)
		}
	}
}

func (receiver httpHandler) reader(conn *websocket.Conn, userAgent string) {
	log := receiver.Logger.Begin()
	defer log.End()
	for {
		room, socket, roomerr := roomwebsocketmap.Map.GetSocket(conn)
		log.Alert("room ", room, " roomerr ", roomerr)
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			_, closedError := err.(*websocket.CloseError)
			_, handshakeError := err.(*websocket.HandshakeError)
			if closedError || handshakeError {
				log.Warn("Socket ", socket.ID, " Closed!")
				roomwebsocketmap.Map.RemoveFromRoom(room.Name, socket)
				log.Inform("Socket closed!")
				return
			}
			log.Error("Read from socket error : ", err)
			return
		}
		if roomerr != nil {
			log.Alert("No Room Socket")
			socket = roomwebsocketmap.MySocket{
				Socket: conn,
			}
		}
		receiver.parseMessage(socket, p, messageType, userAgent)
		log.Alert("messageType ", messageType, " p ", string(p))
	}
}

func (receiver httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log := receiver.Logger.Begin()
	defer log.End()
	log.Highlight("Method", req.Method)
	log.Highlight("Path", req.URL.Path)
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Error("Upgrade Error : ", err)
		return
	}
	log.Log("Client Connected")
	userAgent := req.Header.Get("User-Agent")
	receiver.reader(ws, userAgent)
}
