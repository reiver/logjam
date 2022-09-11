package websockethandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"github.com/sparkscience/logjam/backend/lib/message"
	binarytreesrv "github.com/sparkscience/logjam/backend/srv/binarytree"
	roommapssrv "github.com/sparkscience/logjam/backend/srv/roommaps"

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

var filename string

func (receiver httpHandler) findBroadcaster(roomName string) (bool, *binarytreesrv.MySocket) {
	log := receiver.Logger.Begin()
	defer log.End()

	Map, found := roommapssrv.RoomMaps.Get(roomName)
	if !found {
		log.Errorf("could not get map for room %q when trying to find broadcaster", roomName)
		return false, nil
	}
	broadcasterLevel := Map.LevelNodes(1)
	var broadcaster *binarytreesrv.MySocket
	ok := len(broadcasterLevel) == 1
	if ok {
		broadcaster = broadcasterLevel[0].(*binarytreesrv.MySocket)
	}
	return ok, broadcaster
}

func (receiver httpHandler) parseMessage(socket *binarytreesrv.MySocket, messageJSON []byte, messageType int, userAgent string, roomName string) {
	log := receiver.Logger.Begin()
	defer log.End()

	Map, found := roommapssrv.RoomMaps.Get(roomName)
	if !found {
		log.Errorf("could not get map for room %q when trying to parse message", roomName)
		return
	}
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err == nil {
		defer f.Close()
	}

	var theMessage message.MessageContract
	{
		err := json.Unmarshal(messageJSON, &theMessage)
		if err != nil {
			return
		}
	}

	var response message.MessageContract
	switch theMessage.Type {
	case "start":
		response.Type = "start"
		socket.SetName(theMessage.Data)
		response.Data = strconv.FormatInt(int64(socket.ID), 10)
		responseJSON, err := json.Marshal(response)
		if err == nil {
			socket.Socket.WriteMessage(messageType, responseJSON)
		} else {
			return
		}
	case "role":
		log.Alert("role message received ", theMessage.Data)
		response.Type = "role"
		if theMessage.Data == "unknown" {
			ok, _ := receiver.findBroadcaster(roomName)
			if ok {
				theMessage.Data = "audience"
			} else {
				theMessage.Data = "broadcast"
			}
		}
		log.Alert("role message processed ", theMessage.Data)
		if theMessage.Data == "broadcast" {
			response.Data = "yes:broadcast"
			Map.ToggleHead(socket.Socket)
			Map.ToggleCanConnect(socket.Socket)
			{
				err := roommapssrv.RoomMaps.Set(roomName, Map)
				log.Error("could not set room in map: %s", err)
			}
			if _, err := os.Stat("./logs/" + socket.Name); os.IsNotExist(err) {
				os.Mkdir("./logs/"+socket.Name, 0755)
			}
			currentTime := time.Now()
			logName := currentTime.Format("2006-01-02_15-04-05")
			filename = "./logs/" + socket.Name + "/" + logName + ".log"
			f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err == nil {
				defer f.Close()
			}
		} else if theMessage.Data == "alt-broadcast" {
			response.Data = "yes:broadcast"
			ok, broadcaster := receiver.findBroadcaster(roomName)

			var audianceResponse message.MessageContract
			audianceResponse.Type = "alt-broadcast"
			if !ok {
				audianceResponse.Data = "no-broadcaster"
				audianceResponseJSON, aerr := json.Marshal(audianceResponse)
				if aerr != nil {
					return
				}
				socket.Socket.WriteMessage(messageType, audianceResponseJSON)
				return
			}
			audianceResponse.Data = strconv.FormatInt(int64(broadcaster.ID), 10)
			audianceResponseJSON, aerr := json.Marshal(audianceResponse)
			if aerr != nil {
				return
			}
			var broadResponse map[string]interface{} = make(map[string]interface{})
			broadResponse["Type"] = "alt-broadcast"
			broadResponse["Data"] = strconv.FormatInt(int64(socket.ID), 10)
			broadResponse["name"] = socket.Name
			broadResponseJSON, err := json.Marshal(broadResponse)
			if err != nil {
				return
			}
			broadcaster.Socket.WriteMessage(messageType, broadResponseJSON)
			socket.Socket.WriteMessage(messageType, audianceResponseJSON)
			return
		} else {
			if socket.IsBroadcaster {
				response.Data = "no:broadcast"
				Map.ToggleHead(socket.Socket)
				{
					err := roommapssrv.RoomMaps.Set(roomName, Map)
					log.Error("could not set room in map: %s", err)
				}

				msg := "Broadcaster " + socket.Name + " removed broadcasting role from himself"
				fmt.Fprintln(f, msg)
			} else {
				response.Data = "yes:audience"
				ok, broadcaster := receiver.findBroadcaster(roomName)

				msg := "Audiance " + socket.Name + " trying to receive stream\n" + "    " + userAgent
				fmt.Fprintln(f, msg)

				if !ok {
					response.Data = "no:audience"

					msg := "Audiance " + socket.Name + " could not receive stream because there is no broadcaster now!"
					fmt.Fprintln(f, msg)
				} else {
					var broadResponse message.MessageContract
					broadResponse.Type = "add_audience"
					broadResponse.Data = strconv.FormatInt(int64(socket.ID), 10)
					broadResponseJSON, err := json.Marshal(broadResponse)

					msg := "Audiance " + socket.Name + " is going to get connected"
					fmt.Fprintln(f, msg)

					if err == nil {
						targetSocketNode, e := binarytreesrv.InsertChild(socket.Socket, *Map)
						if e != nil {
							socket.Socket.WriteMessage(messageType, []byte("{\"type\":\"error\",\"data\":\"Insert Child Error : "+e.Error()+" \"}"))
							return
						}
						targetSocket := targetSocketNode.(*binarytreesrv.MySocket)

						msg := "Audiance " + socket.Name + " is starting to connect to " + targetSocket.Name
						fmt.Fprintln(f, msg)

						if targetSocket.Socket != broadcaster.Socket {
							broadResponse.Type = "add_broadcast_audience"
						}
						targetSocket.Socket.WriteMessage(messageType, broadResponseJSON)
					} else {
						return
					}
				}
			}
		}
		responseJSON, err := json.Marshal(response)
		if err == nil {
			socket.Socket.WriteMessage(messageType, responseJSON)
			return
		} else {
			return
		}
	case "stream":
		Map.ToggleCanConnect(socket.Socket)
		{
			err = roommapssrv.RoomMaps.Set(roomName, Map)
			log.Error("could not set room in map: %s", err)
		}

		msg := "Audiance " + socket.Name + " is receiving stream!"
		fmt.Fprintln(f, msg)
	case "turn_status":
		msg := "Turn usage status of " + socket.Name + " is " + theMessage.Data
		socket.SetIsTurn(theMessage.Data == "on")
		fmt.Fprintln(f, msg)
	case "log":
		msg := "Audiance " + socket.Name + " client log :\n    " + theMessage.Data
		fmt.Fprintln(f, msg)
	case "ping":
		socket.Socket.WriteJSON(message.MessageContract{Type: "pong", Data: "pong"})
	case "tree":
		treeData := binarytreesrv.Tree(*Map)
		j, e := json.Marshal(treeData)
		if e != nil {
			return
		}
		output := string(j)
		response.Type = "tree"
		response.Data = output
		responseJSON, err := json.Marshal(response)
		if err == nil {
			socket.Socket.WriteMessage(messageType, responseJSON)
		} else {
			return
		}
	default:
		ID, err := strconv.ParseUint(theMessage.Target, 10, 64)
		if err != nil {
			return
		}
		allSockets := Map.All()
		var ok = false
		var target *binarytreesrv.MySocket
		for _, node := range allSockets {
			if node.(*binarytreesrv.MySocket).ID == ID {
				ok = true
				target = node.(*binarytreesrv.MySocket)
				break
			} else {
				for _, child := range node.All() {
					if child.(*binarytreesrv.MySocket).ID == ID {
						ok = true
						target = child.(*binarytreesrv.MySocket)
						break
					}
				}
			}
		}
		if ok {
			var fullMessage map[string]interface{}
			err := json.Unmarshal(messageJSON, &fullMessage)
			if err != nil {
				return
			}
			fullMessage["username"] = socket.Name
			fullMessage["data"] = strconv.FormatInt(int64(socket.ID), 10)
			messageJSON, err = json.Marshal(fullMessage)
			if err != nil {
				return
			}
			target.Socket.WriteMessage(messageType, messageJSON)
		}
	}
}

func (receiver httpHandler) deleteNode(conn *websocket.Conn, roomName string, messageType int) {
	log := receiver.Logger.Begin()
	defer log.End()
	Map, found := roommapssrv.RoomMaps.Get(roomName)
	if !found {
		log.Errorf("could not get map for room %q when trying to reader", roomName)
	}
	socket := Map.Get(conn).(*binarytreesrv.MySocket)
	var response message.MessageContract
	response.Type = "event-broadcaster-disconnected"
	response.Data = strconv.FormatInt(int64(socket.ID), 10)
	messageTxt, _ := json.Marshal(response)
	if socket.IsBroadcaster {
		for _, s := range Map.All() {
			log.Inform("[deleteNode] checking socket ", s.(*binarytreesrv.MySocket).Name)
			if !s.(*binarytreesrv.MySocket).IsBroadcaster {
				s.(*binarytreesrv.MySocket).Socket.WriteMessage(1, messageTxt)
				log.Inform("[deleteNode] checking socket ", s.(*binarytreesrv.MySocket).Name, " SENT")
				// time.Sleep(2000 * time.Millisecond)
				break
			}
		}
	}
	Map.Delete(conn)
}

func (receiver httpHandler) reader(conn *websocket.Conn, userAgent string, roomName string) {
	log := receiver.Logger.Begin()
	defer log.End()

	defer conn.Close()
	Map, found := roommapssrv.RoomMaps.Get(roomName)
	if !found {
		log.Errorf("could not get map for room %q when trying to reader", roomName)
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			_, closedError := err.(*websocket.CloseError)
			_, handshakeError := err.(*websocket.HandshakeError)
			if closedError || handshakeError {
				receiver.deleteNode(conn, roomName, messageType)
				{
					err := roommapssrv.RoomMaps.Set(roomName, Map)
					if err != nil {
						log.Errorf("could not set room in map: %s", err)
					}
				}
				return
			}
			return
		}

		receiver.parseMessage(Map.Get(conn).(*binarytreesrv.MySocket), p, messageType, userAgent, roomName)
	}
}

func (receiver httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log := receiver.Logger.Begin()
	defer log.End()

	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Error("Upgrade Error : ", err)
		return
	}
	defer ws.Close()
	roomName := req.URL.Query().Get("room")
	var mapFound bool
	_, mapFound = roommapssrv.RoomMaps.Get(roomName)
	if !mapFound {
		AMap := binarytreesrv.GetMap()
		{
			err := roommapssrv.RoomMaps.Set(roomName, &AMap)
			log.Error("could not set room in map: %s", err)
		}
	}
	Map, found := roommapssrv.RoomMaps.Get(roomName)
	if !found {
		receiver.Logger.Errorf("could not get map for room %q when trying to serve-http", roomName)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	Map.Insert(ws)
	{
		err := roommapssrv.RoomMaps.Set(roomName, Map)
		log.Error("could not set room in map: %s", err)
	}
	userAgent := req.Header.Get("User-Agent")
	receiver.reader(ws, userAgent, roomName)
}
