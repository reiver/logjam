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
	"github.com/sparkscience/logjam/backend/lib/websocketmap"

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

func (receiver httpHandler) Reset() {
	websocketmap.Map.Reset()
}

func (receiver httpHandler) levelSockets(level uint) []websocketmap.MySocket {

	var output []websocketmap.MySocket = []websocketmap.MySocket{}
	broadCaster, ok := websocketmap.Map.GetBroadcaster()
	if ok {
		output = append(output, broadCaster)
	}
	if level == 1 || !ok {
		return output
	}
	var index uint = 1
	var currentLevelSockets []websocketmap.MySocket = output
	for {
		if len(currentLevelSockets) == 0 {
			break
		}
		output = []websocketmap.MySocket{}
		for _, socks := range currentLevelSockets {
			for _, child := range socks.ConnectedSockets {
				if child.HasStream {
					output = append(output, child)
				}
			}
		}
		if len(output) == 0 {
			break
		}
		if index == level-1 {
			return output
		}
		currentLevelSockets = output
		output = []websocketmap.MySocket{}
		index++
	}
	return output
}

func (receiver httpHandler) decideWhomToConnect(broadcaster websocketmap.MySocket) websocketmap.MySocket {
	var level uint = 2
	var levelSockets []websocketmap.MySocket = []websocketmap.MySocket{}
	for {
		levelSockets = receiver.levelSockets(level)
		parentSockets := receiver.levelSockets(level - 1)
		if len(levelSockets) == 0 {
			return parentSockets[0]
		}
		for i := 0; i < len(parentSockets); i++ {
			if len(parentSockets[i].ConnectedSockets) < 2 {
				return parentSockets[i]
			}
		}

		level++
	}
}

var filename string

func (receiver httpHandler) parseMessage(socket websocketmap.MySocket, messageJSON []byte, messageType int) {
	log := receiver.Logger.Begin()
	defer log.End()
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err == nil {
		defer f.Close()
	}

	var theMessage message.MessageContract
	{
		err := json.Unmarshal(messageJSON, &theMessage)
		if err != nil {
			log.Error("Error unmarshal message ", err)
			return
		}
		// log.Highlight("TheMessage  Type : ", theMessage.Type, " Data : ", theMessage.Data, " Target : ", theMessage.Target)
	}

	var response message.MessageContract
	switch theMessage.Type {
	case "start":
		response.Type = "start"
		websocketmap.Map.SetName(socket.Socket, theMessage.Data)
		response.Data = strconv.FormatInt(int64(socket.ID), 10)
		responseJSON, err := json.Marshal(response)
		if err == nil {
			socket.Socket.WriteMessage(messageType, responseJSON)
		} else {
			log.Error("Marshal Error of `start` response", err)
			return
		}
	case "role":
		response.Type = "role"
		if theMessage.Data == "broadcast" {
			response.Data = "yes:broadcast"
			websocketmap.Map.RemoveBroadcasters()
			websocketmap.Map.SetBroadcaster(socket.Socket)
			if _, err := os.Stat("./logs/" + socket.Name); os.IsNotExist(err) {
				os.Mkdir("./logs/"+socket.Name, os.ModeDir)
			}
			currentTime := time.Now()
			logName := currentTime.Format("2006-01-02_15-04-05")
			filename = "./logs/" + socket.Name + "/" + logName + ".log"
			f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err == nil {
				log.Alert("Set output to file")
				defer f.Close()
			}
			fmt.Fprintln(f, "Broadcast "+socket.Name+" Started")
			f.Sync()
		} else {
			log.Highlight("Filename : ", filename)
			if socket.IsBroadcaster {
				response.Data = "no:broadcast"
				websocketmap.Map.RemoveBroadcaster(socket.Socket)

				msg := "Broadcaster " + socket.Name + " removed broadcasting role from himself"
				fmt.Fprintln(f, msg)
				f.Sync()
			} else {
				response.Data = "yes:audience"
				broadcaster, ok := websocketmap.Map.GetBroadcaster()

				msg := "Audiance " + socket.Name + " trying to receive stream"
				fmt.Fprintln(f, msg)

				if !ok {
					response.Data = "no:audience"

					msg := "Audiance " + socket.Name + " could not receive stream because there is no broadcaster now!"
					fmt.Fprintln(f, msg)
				} else {
					var broadResponse message.MessageContract
					broadResponse.Type = "add_audience"
					broadResponse.Data = strconv.FormatInt(int64(socket.ID), 10) // + "," + socket.Name
					broadResponseJSON, err := json.Marshal(broadResponse)

					msg := "Audiance " + socket.Name + " is going to get connected"
					fmt.Fprintln(f, msg)

					if err == nil {
						log.Highlight("Deciding to connect ...")
						targetSocket := receiver.decideWhomToConnect(broadcaster)
						log.Highlight("target ", targetSocket.ID)

						msg := "Audiance " + socket.Name + " is starting to connect to " + targetSocket.Name
						fmt.Fprintln(f, msg)

						if targetSocket.Socket != broadcaster.Socket {
							broadResponse.Type = "add_broadcast_audience"
						}
						websocketmap.Map.InsertConnected(targetSocket.Socket, socket.Socket)
						// log.Informf("Target Socket has %d sockets connected!", len(websocketmap.Map.Get(targetSocket.Socket).ConnectedSockets))
						targetSocket.Socket.WriteMessage(messageType, broadResponseJSON)
						level := 1
						for {
							if len(receiver.levelSockets(uint(level))) == 0 {
								break
							}
							// log.Highlightf("Level %d: %d", level, len(receiver.levelSockets(uint(level))))
							level++
						}
					} else {
						log.Error("Marshal Error of `add_audience` broadResponse", err)
						return
					}
				}
			}
		}
		responseJSON, err := json.Marshal(response)
		if err == nil {
			socket.Socket.WriteMessage(messageType, responseJSON)
		} else {
			log.Error("Marshal Error of `role` response", err)
			return
		}
	case "stream":
		log.Alert("Stream Received ", socket.Name)
		websocketmap.Map.SetStreamState(socket.Socket, true)

		msg := "Audiance " + socket.Name + " is receiving stream!"
		fmt.Fprintln(f, msg)
	case "log":
		log.Alert("Log Received ", socket.Name)
		msg := "Audiance " + socket.Name + " client log :\n    " + theMessage.Data
		fmt.Fprintln(f, msg)
	default:
		ID, err := strconv.ParseUint(theMessage.Target, 10, 64)
		if err != nil {
			log.Error("Inavlid Target : ", theMessage.Target)
			return
		}
		target, ok := websocketmap.Map.GetSocketByID(ID)
		if ok {
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

func (receiver httpHandler) reader(conn *websocket.Conn) {
	log := receiver.Logger.Begin()
	defer log.End()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			_, closedError := err.(*websocket.CloseError)
			_, handshakeError := err.(*websocket.HandshakeError)
			if closedError || handshakeError {
				log.Warn("Socket ", websocketmap.Map.Get(conn).ID, " Closed!")
				websocketmap.Map.Delete(conn)
				log.Inform("Socket closed!")
				return
			}
			log.Error("Read from socket error : ", err)
			return
		}
		receiver.parseMessage(websocketmap.Map.Get(conn), p, messageType)
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
	websocketmap.Map.Insert(ws)
	receiver.reader(ws)
}
