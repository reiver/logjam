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
	// "github.com/sparkscience/logjam/backend/lib/websocketmap"
	binarytreesrv "github.com/sparkscience/logjam/backend/srv/binarytree"

	"github.com/mmcomp/go-log"
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

var Map = binarytreesrv.GetMap()

func Handler(logger logger.Logger) http.Handler {
	return httpHandler{
		Logger: logger,
	}
}

var filename string

func (receiver httpHandler) findBroadcaster() (bool, *binarytreesrv.MySocket) {
	broadcasterLevel := Map.LevelNodes(1)
	var broadcaster *binarytreesrv.MySocket
	ok := len(broadcasterLevel) == 1
	if ok {
		broadcaster = broadcasterLevel[0].(*binarytreesrv.MySocket)
	}
	return ok, broadcaster
}

func (receiver httpHandler) parseMessage(socket *binarytreesrv.MySocket, messageJSON []byte, messageType int, userAgent string) {
	// log := receiver.Logger.Begin()
	// defer log.End()
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err == nil {
		defer f.Close()
	}

	var theMessage message.MessageContract
	{
		err := json.Unmarshal(messageJSON, &theMessage)
		if err != nil {
			// log.Error("Error unmarshal message ", err)
			return
		}
		// log.Highlight("TheMessage  Type : ", theMessage.Type, " Data : ", theMessage.Data, " Target : ", theMessage.Target)
	}

	var response message.MessageContract
	switch theMessage.Type {
	case "start":
		response.Type = "start"
		// websocketmap.Map.SetName(socket.Socket, theMessage.Data)
		socket.SetName(theMessage.Data)
		response.Data = strconv.FormatInt(int64(socket.ID), 10)
		responseJSON, err := json.Marshal(response)
		if err == nil {
			socket.Socket.WriteMessage(messageType, responseJSON)
		} else {
			// log.Error("Marshal Error of `start` response", err)
			return
		}
	case "role":
		response.Type = "role"
		if theMessage.Data == "broadcast" {
			response.Data = "yes:broadcast"
			// websocketmap.Map.RemoveBroadcasters()
			// websocketmap.Map.SetBroadcaster(socket.Socket)
			Map.ToggleHead(socket.Socket)
			Map.ToggleCanConnect(socket.Socket)
			if _, err := os.Stat("./logs/" + socket.Name); os.IsNotExist(err) {
				os.Mkdir("./logs/"+socket.Name, 0755)
			}
			currentTime := time.Now()
			logName := currentTime.Format("2006-01-02_15-04-05")
			filename = "./logs/" + socket.Name + "/" + logName + ".log"
			f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err == nil {
				// log.Alert("Set output to file")
				defer f.Close()
			}
			fmt.Fprintln(f, "Broadcast "+socket.Name+" Started")
		} else if theMessage.Data == "alt-broadcast" {
			response.Data = "yes:broadcast"
			fmt.Fprintln(f, "Alt Broadcast "+socket.Name+" Started")
			ok, broadcaster := receiver.findBroadcaster()

			var audianceResponse message.MessageContract
			audianceResponse.Type = "alt-broadcast"
			if !ok {
				audianceResponse.Data = "no-broadcaster"
				audianceResponseJSON, aerr := json.Marshal(audianceResponse)
				if aerr != nil {
					// log.Error("Marshal Error of `alt-broadcast` audianceResponse1", aerr)
					return
				}
				socket.Socket.WriteMessage(messageType, audianceResponseJSON)
				return
			}
			audianceResponse.Data = strconv.FormatInt(int64(broadcaster.ID), 10)
			audianceResponseJSON, aerr := json.Marshal(audianceResponse)
			if aerr != nil {
				// log.Error("Marshal Error of `alt-broadcast` audianceResponse2", aerr)
				return
			}
			var broadResponse map[string]interface{} = make(map[string]interface{})
			broadResponse["Type"] = "alt-broadcast"
			broadResponse["Data"] = strconv.FormatInt(int64(socket.ID), 10)
			broadResponse["name"] = socket.Name
			broadResponseJSON, err := json.Marshal(broadResponse)
			if err != nil {
				// log.Error("Marshal Error of `alt-broadcast` broadResponse", err)
				return
			}
			broadcaster.Socket.WriteMessage(messageType, broadResponseJSON)
			socket.Socket.WriteMessage(messageType, audianceResponseJSON)
			return
		} else {
			// log.Highlight("userAgent : ", userAgent)
			if socket.IsBroadcaster {
				response.Data = "no:broadcast"
				// websocketmap.Map.RemoveBroadcaster(socket.Socket)
				Map.ToggleHead(socket.Socket)

				msg := "Broadcaster " + socket.Name + " removed broadcasting role from himself"
				fmt.Fprintln(f, msg)
			} else {
				response.Data = "yes:audience"
				ok, broadcaster := receiver.findBroadcaster()
				// broadcasterLevel := Map.LevelNodes(1)
				// broadcaster, ok := websocketmap.Map.GetBroadcaster()
				// var broadcaster *binarytreesrv.MySocket
				// ok := len(broadcasterLevel) == 1
				// log.Alert("broadcasterLevel ", broadcasterLevel)

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
						// log.Highlight("Deciding to connect ...")
						targetSocketNode, e := binarytreesrv.InsertChild(socket.Socket, Map)
						// targetSocketNode, e := Map.InsertChild(socket.Socket, false)
						if e != nil {
							// log.Error("Insert Child Error ", e)
							socket.Socket.WriteMessage(messageType, []byte("{\"type\":\"error\",\"data\":\"Insert Child Error : "+e.Error()+" \"}"))
							return
						}
						// log.Highlight("Deciding to connect end")
						targetSocket := targetSocketNode.(*binarytreesrv.MySocket)
						// log.Highlight("target ", targetSocket.ID)

						msg := "Audiance " + socket.Name + " is starting to connect to " + targetSocket.Name
						fmt.Fprintln(f, msg)

						if targetSocket.Socket != broadcaster.Socket {
							broadResponse.Type = "add_broadcast_audience"
						}
						targetSocket.Socket.WriteMessage(messageType, broadResponseJSON)
					} else {
						// log.Error("Marshal Error of `add_audience` broadResponse", err)
						return
					}
				}
			}
		}
		responseJSON, err := json.Marshal(response)
		if err == nil {
			socket.Socket.WriteMessage(messageType, responseJSON)
		} else {
			// log.Error("Marshal Error of `role` response", err)
			return
		}
	case "stream":
		// log.Alert("Stream Received ", socket.Name)
		Map.ToggleCanConnect(socket.Socket)

		msg := "Audiance " + socket.Name + " is receiving stream!"
		fmt.Fprintln(f, msg)
	case "turn_status":
		msg := "Turn usage status of " + socket.Name + " is " + theMessage.Data
		socket.SetIsTurn(theMessage.Data == "on")
		fmt.Fprintln(f, msg)
	case "log":
		// log.Alert("Log Received ", socket.Name)
		msg := "Audiance " + socket.Name + " client log :\n    " + theMessage.Data
		fmt.Fprintln(f, msg)
	case "tree":
		// log.Alert("Tree Received ", socket.Name)
		treeData := binarytreesrv.Tree(Map)
		// log.Inform("treeData ", treeData)
		j, e := json.Marshal(treeData)
		if e != nil {
			log.Error(e)
			return
		}
		output := string(j)
		// log.Inform("treeData ", output)
		response.Type = "tree"
		response.Data = output
		responseJSON, err := json.Marshal(response)
		if err == nil {
			socket.Socket.WriteMessage(messageType, responseJSON)
		} else {
			// log.Error("Marshal Error of `tree` response", err)
			return
		}
	default:
		// log.Alert("Default Message")
		ID, err := strconv.ParseUint(theMessage.Target, 10, 64)
		if err != nil {
			// log.Error("Inavlid Target : ", theMessage.Target)
			return
		}
		// log.Alert("Target ", ID)
		allockets := Map.All()
		var ok = false
		var target binarytreesrv.MySocket
		// target, ok := Map. .GetSocketByID(ID)
		for _, node := range allockets {
			// log.Alert(node.(*binarytreesrv.MySocket).ID)
			if node.(*binarytreesrv.MySocket).ID == ID {
				ok = true
				target = *node.(*binarytreesrv.MySocket)
				break
			} else {
				for _, child := range node.All() {
					// log.Alert(child.(*binarytreesrv.MySocket).ID)
					if child.(*binarytreesrv.MySocket).ID == ID {
						ok = true
						target = *child.(*binarytreesrv.MySocket)
						break
					}
				}
			}
		}
		// log.Alert("Find target ? ", ok)
		if ok {
			// log.Inform("Default sending to ", ID, " ", string(messageJSON))
			var fullMessage map[string]interface{}
			err := json.Unmarshal(messageJSON, &fullMessage)
			if err != nil {
				// log.Error("Error unmarshal message ", err)
				return
			}
			fullMessage["username"] = socket.Name
			fullMessage["data"] = strconv.FormatInt(int64(socket.ID), 10)
			messageJSON, err = json.Marshal(fullMessage)
			if err != nil {
				// log.Error("Error marshal message ", err)
				return
			}
			target.Socket.WriteMessage(messageType, messageJSON)
		}
	}
}

func (receiver httpHandler) reader(conn *websocket.Conn, userAgent string) {
	// log := receiver.Logger.Begin()
	// defer log.End()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			_, closedError := err.(*websocket.CloseError)
			_, handshakeError := err.(*websocket.HandshakeError)
			if closedError || handshakeError {
				// log.Warn("Socket ", Map.Get(conn).(*binarytreesrv.MySocket).ID, " Closed!")
				Map.Delete(conn)
				// log.Inform("Socket closed!")
				return
			}
			// log.Error("Read from socket error : ", err)
			return
		}
		// log.Informf("Message from socket ID %d name %s ", Map.Get(conn).(*binarytreesrv.MySocket).ID, Map.Get(conn).(*binarytreesrv.MySocket).Name)
		// log.Inform("Message", string(p))
		receiver.parseMessage(Map.Get(conn).(*binarytreesrv.MySocket), p, messageType, userAgent)
		// log.Alert(messageType, p, Map.Get(conn))
	}

}

func (receiver httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log := receiver.Logger.Begin()
	defer log.End()
	// log.Highlight("Method", req.Method)
	// log.Highlight("Path", req.URL.Path)
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Error("Upgrade Error : ", err)
		return
	}
	// log.Log("Client Connected")
	Map.Insert(ws)
	userAgent := req.Header.Get("User-Agent")
	receiver.reader(ws, userAgent)
}
