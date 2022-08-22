package roomwebsockethandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/sparkscience/logjam/backend/lib/message"

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

type RoomSocket struct {
	WS   *websocket.Conn
	Name string
}

// ------ Room ------
type Room struct {
	mutex   sync.Mutex
	Sockets map[*websocket.Conn]RoomSocket
	Name    string
}

func (receiver *Room) Insert(roomSocket RoomSocket) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	_, ok := receiver.Sockets[roomSocket.WS]
	if !ok {
		receiver.Sockets[roomSocket.WS] = roomSocket
	}
}

func (receiver *Room) Delete(ws *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	delete(receiver.Sockets, ws)
}

//\------ Room ------

// ------ Sockets ------
type Sockets struct {
	mutex      sync.Mutex
	AllSockets map[*websocket.Conn]RoomSocket
}

func (receiver *Sockets) Insert(ws *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if receiver.AllSockets == nil {
		receiver.AllSockets = make(map[*websocket.Conn]RoomSocket)
	}

	receiver.AllSockets[ws] = RoomSocket{
		WS:   ws,
		Name: "NoName",
	}
}

func (receiver *Sockets) Delete(ws *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	delete(receiver.AllSockets, ws)
}

func (receiver *Sockets) SetName(ws *websocket.Conn, name string) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	sock, ok := receiver.AllSockets[ws]
	if ok {
		sock.Name = name
		receiver.AllSockets[ws] = sock
	}
}

func (receiver *Sockets) FindByName(name string) (RoomSocket, bool) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for _, s := range receiver.AllSockets {
		if s.Name == name {
			return s, true
		}
	}

	return RoomSocket{}, false
}

//\------ Sockets ------

// ------ RoomValues ------
type RoomValues struct {
	mutex  sync.Mutex
	Values map[string]Room
}

func (receiver *RoomValues) Insert(roomName string, room Room) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if receiver.Values == nil {
		receiver.Values = make(map[string]Room)
	}

	receiver.Values[roomName] = room
}

func (receiver *RoomValues) RemoveSocket(ws *websocket.Conn) {
	for _, room := range receiver.Values {
		room.Delete(ws)
	}
}

//\------ RoomValues ------

var Rooms RoomValues = RoomValues{}
var AllSocks Sockets = Sockets{}

func Handler(logger logger.Logger) http.Handler {
	return httpHandler{
		Logger: logger,
	}
}

var filename string

func (receiver httpHandler) parseMessage(socket *websocket.Conn, messageJSON []byte, messageType int, userAgent string) {
	log := receiver.Logger.Begin()
	defer log.End()
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err == nil {
		defer f.Close()
	}
	theSocket, ok := AllSocks.AllSockets[socket]
	log.Alert("Socket Exists ? ", ok)
	var theMessage message.MessageContract
	{
		err := json.Unmarshal(messageJSON, &theMessage)
		if err != nil {
			log.Error("Error unmarshal message ", err)
			return
		}
	}

	var response message.MessageContract
	switch theMessage.Type {
	case "start":
		response.Type = "start"
		AllSocks.SetName(socket, theMessage.Data)
		response.Data = theMessage.Data
		responseJSON, err := json.Marshal(response)
		if err == nil {
			socket.WriteMessage(messageType, responseJSON)
		} else {
			log.Error("Marshal Error of `start` response", err)
			return
		}
	case "open-room":
		roomName := theMessage.Data
		if _, err := os.Stat("./logs/" + roomName); os.IsNotExist(err) {
			os.Mkdir("./logs/"+roomName, 0755)
		}
		currentTime := time.Now()
		logName := currentTime.Format("2006-01-02_15-04-05")
		filename = "./logs/" + roomName + "/" + logName + ".log"
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err == nil {
			defer f.Close()
		}
		room, ok := Rooms.Values[roomName]
		if !ok {
			room = Room{
				Sockets: make(map[*websocket.Conn]RoomSocket),
				Name:    roomName,
			}
			Rooms.Insert(roomName, room)
			fmt.Fprintln(f, "Room opened by "+theSocket.Name)
		}

		fmt.Fprintln(f, theSocket.Name+" joined the room")
		room.Insert(theSocket)
		response.Type = "open-room"
		response.Data = ""
		for _, so := range room.Sockets {
			if response.Data != "" {
				response.Data += ","
			}

			response.Data += so.Name
		}
		responseJSON, err := json.Marshal(response)
		if err == nil {
			socket.WriteMessage(messageType, responseJSON)
		} else {
			log.Error("Marshal Error of `start` response", err)
			return
		}
	default:
		log.Alert("Default Message Target ", theMessage.Target)
		target, ok := AllSocks.FindByName(theMessage.Target)
		log.Alert("Find target ? ", ok)
		if ok {
			var fullMessage map[string]interface{}
			err := json.Unmarshal(messageJSON, &fullMessage)
			if err != nil {
				log.Error("Error unmarshal message ", err)
				return
			}
			me := AllSocks.AllSockets[socket]
			fullMessage["name"] = me.Name
			messageJSON, err = json.Marshal(fullMessage)
			if err != nil {
				log.Error("Error marshal message ", err)
				return
			}
			target.WS.WriteMessage(messageType, messageJSON)
		}
	}
}

func (receiver httpHandler) reader(conn *websocket.Conn, userAgent string) {
	log := receiver.Logger.Begin()
	defer log.End()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			_, closedError := err.(*websocket.CloseError)
			_, handshakeError := err.(*websocket.HandshakeError)
			if closedError || handshakeError {
				// log.Warn("Socket ", Map.Get(conn).(*binarytreesrv.MySocket).ID, " Closed!")
				Rooms.RemoveSocket(conn)
				AllSocks.Delete(conn)
				log.Inform("Socket closed!")
				return
			}
			log.Error("Read from socket error : ", err)
			return
		}
		// log.Informf("Message from socket ID %d name %s ", Map.Get(conn).(*binarytreesrv.MySocket).ID, Map.Get(conn).(*binarytreesrv.MySocket).Name)
		log.Inform("Message", string(p))
		receiver.parseMessage(conn, p, messageType, userAgent)
		// log.Alert(messageType, p, Map.Get(conn))
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
	AllSocks.Insert(ws)
	userAgent := req.Header.Get("User-Agent")
	receiver.reader(ws, userAgent)
}
