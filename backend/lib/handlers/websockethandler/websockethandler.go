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

	"github.com/mmcomp/go-binarytree"
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

func (receiver httpHandler) emitUserEvent(roomName string) error {
	log := receiver.Logger.Begin()
	defer log.End()

	users, err := roommapssrv.RoomMaps.GetUsers(roomName)
	if err != nil {
		return err
	}

	usersJsonString, err := json.Marshal(users)
	if err != nil {
		return err
	}

	msg := message.MessageContract{
		Type: "user-event",
		Data: string(usersJsonString),
	}

	receiver.broadcastMessage(roomName, 1, msg)
	return nil
}

func (receiver httpHandler) findBroadcaster(roomName string) (bool, *binarytreesrv.MySocket) {
	log := receiver.Logger.Begin()
	defer log.End()

	Map, found := roommapssrv.RoomMaps.Get(roomName)
	if !found {
		log.Errorf("could not get map for room %q when trying to find broadcaster", roomName)
		return false, nil
	}
	broadcasterLevel := Map.Room.LevelNodes(1)
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
			socket.Writer.WriteMessage(messageType, responseJSON)
		} else {
			return
		}
	case "role":
		var msgData map[string]string
		err := json.Unmarshal(messageJSON, &msgData)
		if err != nil {
			return
		}
		streamId, ok := msgData["streamId"]
		if ok {
			if socket.MetaData == nil {
				socket.MetaData = make(map[string]string)
			}
			socket.MetaData["streamId"] = streamId
		}

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
			Map.Room.ToggleHead(socket.Socket)
			Map.Room.ToggleCanConnect(socket.Socket)
			{
				err := roommapssrv.RoomMaps.Set(roomName, Map.Room)
				if err != nil {
					log.Errorf("could not set room in map: %s", err)
					return
				}
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

			receiver.broadcastMessage(roomName, messageType, message.MessageContract{
				Type: "broadcasting",
				Data: strconv.FormatInt(int64(socket.ID), 10),
			})
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
				socket.Writer.WriteMessage(messageType, audianceResponseJSON)
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
			broadcaster.Writer.WriteMessage(messageType, broadResponseJSON)
			socket.Writer.WriteMessage(messageType, audianceResponseJSON)
			receiver.emitUserEvent(roomName)
			return
		} else {
			if socket.IsBroadcaster {
				response.Data = "no:broadcast"
				Map.Room.ToggleHead(socket.Socket)
				{
					err := roommapssrv.RoomMaps.Set(roomName, Map.Room)
					if err != nil {
						log.Errorf("could not set room in map: %s", err)
						return
					}
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
						targetSocketNode, e := binarytreesrv.InsertChild(socket.Socket, Map.Room)
						if e != nil {
							socket.Writer.WriteMessage(messageType, []byte("{\"type\":\"error\",\"data\":\"Insert Child Error : "+e.Error()+" \"}"))
							return
						}
						targetSocket := targetSocketNode.(*binarytreesrv.MySocket)

						msg := "Audiance " + socket.Name + " is starting to connect to " + targetSocket.Name
						fmt.Fprintln(f, msg)

						if targetSocket.Socket != broadcaster.Socket {
							broadResponse.Type = "add_broadcast_audience"
						}
						targetSocket.Writer.WriteMessage(messageType, broadResponseJSON)
					} else {
						return
					}
				}
			}
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			return
		}
		socket.Writer.WriteMessage(messageType, responseJSON)
		receiver.emitUserEvent(roomName)
		return
	case "stream":
		if !socket.CanConnect() {
			Map.Room.ToggleCanConnect(socket.Socket)
			{
				err = roommapssrv.RoomMaps.Set(roomName, Map.Room)
				if err != nil {
					log.Errorf("could not set room in map: %s", err)
					return
				}
			}

			msg := "Audiance " + socket.Name + " is receiving stream!"
			fmt.Fprintln(f, msg)
		}
	case "turn_status":
		msg := "Turn usage status of " + socket.Name + " is " + theMessage.Data
		socket.SetIsTurn(theMessage.Data == "on")
		fmt.Fprintln(f, msg)
	case "log":
		msg := "Audiance " + socket.Name + " client log :\n    " + theMessage.Data
		fmt.Fprintln(f, msg)
	case "ping":
		socket.Writer.WriteJSON(message.MessageContract{Type: "pong", Data: "pong"})
	case "tree":
		treeData := binarytreesrv.Tree(*Map.Room)
		j, e := json.Marshal(treeData)
		if e != nil {
			return
		}
		output := string(j)
		response.Type = "tree"
		response.Data = output
		responseJSON, err := json.Marshal(response)
		if err == nil {
			socket.Writer.WriteMessage(messageType, responseJSON)
		} else {
			return
		}
	case "metadata-set":
		log.Alert("Set Meta Data ", theMessage.Data)
		metaData := make(map[string]string)
		metaDataJsonError := json.Unmarshal([]byte(theMessage.Data), &metaData)
		log.Alert("Set Meta Data ", metaData)
		if metaDataJsonError == nil {
			roommapssrv.RoomMaps.SetMetData(roomName, metaData)
			response.Type = "metadata-set"
			response.Data = theMessage.Data
			responseJSON, err := json.Marshal(response)
			if err == nil {
				socket.Writer.WriteMessage(messageType, responseJSON)
			}
		}

		return
	case "metadata-get":
		// log.Alert("metadata-get")
		metaData := Map.MetaData
		metaDataJson, metaDataJsonError := json.Marshal(metaData)
		if metaDataJsonError == nil {
			response.Type = "metadata-get"
			response.Data = string(metaDataJson)
			responseJSON, err := json.Marshal(response)
			if err == nil {
				// log.Alert(responseJSON)
				socket.Writer.WriteMessage(messageType, responseJSON)
			} else {
				return
			}
		}

		return
	case "user-metadata-get":
		// log.Alert("metadata-get")
		metaData := socket.MetaData
		metaDataJson, metaDataJsonError := json.Marshal(metaData)
		if metaDataJsonError == nil {
			response.Type = "metadata-get"
			response.Data = string(metaDataJson)
			responseJSON, err := json.Marshal(response)
			if err == nil {
				// log.Alert(responseJSON)
				socket.Writer.WriteMessage(messageType, responseJSON)
			} else {
				return
			}
		}

		return
	case "user-metadata-set":
		log.Alert("Set Meta Data ", theMessage.Data)
		metaData := make(map[string]string)
		metaDataJsonError := json.Unmarshal([]byte(theMessage.Data), &metaData)
		log.Alert("Set Meta Data ", metaData)
		if metaDataJsonError == nil {
			socket.SetMetaData(Map.MetaData)
			roommapssrv.RoomMaps.SetMetData(roomName, metaData)
			response.Type = "metadata-set"
			response.Data = theMessage.Data
			responseJSON, err := json.Marshal(response)
			if err == nil {
				socket.Writer.WriteMessage(messageType, responseJSON)
			}
		}

		return
	case "user-by-stream":
		log.Alert("Get User buy Stream Id ", theMessage.Data)
		node, e := roommapssrv.RoomMaps.GetSocketByStreamId(roomName, theMessage.Data)
		log.Alert("Get User buy Stream Id res", node)
		if e == nil {
			user := node.(*binarytreesrv.MySocket)
			response.Type = "user-by-stream"
			userRole := "audience"
			if user.IsBroadcaster {
				userRole = "broadcast"
			}
			response.Data = strconv.FormatInt(int64(user.ID), 10) + "," + user.Name + "," + theMessage.Data + "," + userRole
			responseJSON, err := json.Marshal(response)
			if err == nil {
				socket.Writer.WriteMessage(messageType, responseJSON)
			}
		}
		return
	case "broadcaster-status":
		ok, _ := receiver.findBroadcaster(roomName)
		data := "available"
		if !ok {
			data = "unavailable"
		}
		log.Alert("broadcaster-status ", data)
		response.Type = "broadcaster-status"
		response.Data = data
		responseJSON, err := json.Marshal(response)
		if err == nil {
			log.Alert("broadcaster-status sending ", string(responseJSON))
			socket.Writer.WriteMessage(messageType, responseJSON)
		} else {
			return
		}

	default:
		ID, err := strconv.ParseUint(theMessage.Target, 10, 64)
		if err != nil {
			return
		}
		allSockets := Map.Room.All()
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
			target.Writer.WriteMessage(messageType, messageJSON)
		}
	}
}

func (receiver httpHandler) broadcastMessage(roomName string, messageType int, message message.MessageContract) {
	log := receiver.Logger.Begin()
	defer log.End()
	Map, found := roommapssrv.RoomMaps.Get(roomName)
	if !found {
		log.Errorf("could not get map for room %q when trying to reader", roomName)
		return
	}
	messageTxt, err := json.Marshal(message)
	if err != nil {
		return
	}

	for _, s := range Map.Room.All() {
		s.(*binarytreesrv.MySocket).Writer.WriteMessage(messageType, messageTxt)
		log.Inform("[broadcastMessage] socket ", s.(*binarytreesrv.MySocket).Name, " SENT ", messageTxt)
	}
}

func (receiver httpHandler) deleteNode(conn *websocket.Conn, roomName string, messageType int) (socketId int64) {
	log := receiver.Logger.Begin()
	defer log.End()
	Map, found := roommapssrv.RoomMaps.Get(roomName)
	if !found {
		log.Errorf("could not get map for room %q when trying to reader", roomName)
		return
	}
	socket := Map.Room.Get(conn).(*binarytreesrv.MySocket)
	socketId = int64(socket.ID)
	var response message.MessageContract
	response.Type = "event-broadcaster-disconnected"
	response.Data = strconv.FormatInt(int64(socket.ID), 10)
	messageTxt, _ := json.Marshal(response)
	response.Type = "event-reconnect"
	otherMessageTxt, _ := json.Marshal(response)
	response.Type = "event-parent-dc"
	theMessageTxt, _ := json.Marshal(response)
	var chosenOne binarytree.SingleNode
	if socket.IsBroadcaster {
		for _, s := range Map.Room.All() {
			log.Inform("[deleteNode] checking socket ", s.(*binarytreesrv.MySocket).Name)
			if chosenOne == nil {
				chosenOne = s
				s.(*binarytreesrv.MySocket).Writer.WriteMessage(1, messageTxt)
				log.Inform("[deleteNode] checking socket ", s.(*binarytreesrv.MySocket).Name, " SENT")
			} else {
				s.(*binarytreesrv.MySocket).Writer.WriteMessage(1, otherMessageTxt)
				log.Inform("[deleteNode] checking other socket ", s.(*binarytreesrv.MySocket).Name, " SENT")
			}
		}
	} else {
		for _, s := range socket.ConnectedSockets {
			s.Writer.WriteMessage(1, theMessageTxt)
		}
	}
	Map.Room.Delete(conn)
	receiver.emitUserEvent(roomName)
	return
}

func (receiver httpHandler) reader(conn *websocket.Conn, userAgent string, roomName string) {
	log := receiver.Logger.Begin()
	defer log.End()

	defer conn.Close()
	Map, found := roommapssrv.RoomMaps.Get(roomName)
	if !found {
		log.Errorf("could not get map for room %q when trying to reader", roomName)
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			_, closedError := err.(*websocket.CloseError)
			_, handshakeError := err.(*websocket.HandshakeError)
			if closedError || handshakeError {
				socketId := receiver.deleteNode(conn, roomName, messageType)
				{
					if raiseHandsStr, _ := roommapssrv.RoomMaps.GetFromMetaData(roomName, "raiseHands"); raiseHandsStr != nil {
						var usernameList []string
						err = json.Unmarshal([]byte(*raiseHandsStr), &usernameList)
						if err != nil {
							log.Errorf("could not read room metaData: %s", err)
						} else {
							if len(usernameList) > 0 {
								i := indexOf(usernameList, strconv.FormatInt(socketId, 10))
								if i >= 0 {
									usernameList = append(usernameList[:i], usernameList[i+1:]...)
									if len(usernameList) == 0 {
										roommapssrv.RoomMaps.DelFromMetaData(roomName, "raiseHands")
									} else {
										jsonBytes, err := json.Marshal(usernameList)
										if err != nil {
											log.Errorf("%s", err)
										}
										roommapssrv.RoomMaps.SetToMetaData(roomName, "raiseHands", string(jsonBytes))
									}
								}
							}
						}
					}
				}
			}
			return
		}

		receiver.parseMessage(Map.Room.Get(conn).(*binarytreesrv.MySocket), p, messageType, userAgent, roomName)
	}
}

func indexOf(array []string, element string) int {
	for i, v := range array {
		if v == element {
			return i
		}
	}
	return -1
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
			if err != nil {
				log.Errorf("could not set room in map: %s", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	}
	Map, found := roommapssrv.RoomMaps.Get(roomName)
	if !found {
		log.Errorf("could not get map for room %q when trying to serve-http", roomName)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	Map.Room.Insert(ws)
	{
		err := roommapssrv.RoomMaps.Set(roomName, Map.Room)
		if err != nil {
			log.Errorf("could not set room in map: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	userAgent := req.Header.Get("User-Agent")
	receiver.reader(ws, userAgent, roomName)
}
