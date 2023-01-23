package websocketsrv

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	logger "github.com/mmcomp/go-log"
	"github.com/sparkscience/logjam/backend/lib/message"
	binarytreesrv "github.com/sparkscience/logjam/backend/srv/binarytree"
	roommapssrv "github.com/sparkscience/logjam/backend/srv/roommaps"
	"net/http"
)

type WSHandlerFunc func(ws *binarytreesrv.MySocket, msg message.MessageContract, msgType int)
type WSHandlerService struct {
	handlers       map[string]WSHandlerFunc
	wsUpgrader     websocket.Upgrader
	allowedOrigins []string
	onConnected    func(wsConn *websocket.Conn, roomName string) (error, int)
	onDisconnected func(wsConn *websocket.Conn)
	logger         logger.Logger
}

func (w *WSHandlerService) shouldAllowOrigin(origin string) bool {
	// if w.allowedOrigins contains origin return true  else  return false
	///todo: check the origin from w.allowedOrigins list.
	return true
}

func (w *WSHandlerService) RegisterWSHandler(cmd string, handler WSHandlerFunc) {
	w.handlers[cmd] = handler
}

func Handler(logger logger.Logger, allowedOrigins []string, onConnected func(wsConn *websocket.Conn, roomName string) (error, int), onDisconnected func(wsConn *websocket.Conn)) *WSHandlerService {
	srv := WSHandlerService{
		handlers: map[string]WSHandlerFunc{},
		wsUpgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		allowedOrigins: allowedOrigins,
		onConnected:    onConnected,
		onDisconnected: onDisconnected,
		logger:         logger,
	}
	srv.wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return srv.shouldAllowOrigin(r.Header.Get("Origin"))
	}
	return &srv
}

func (w *WSHandlerService) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log := w.logger.Begin()
	defer log.End()
	ws, err := w.wsUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Error("Upgrade Error : ", err)
		return
	}
	defer func() {
		err := ws.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	roomName := "malmal"
	userAgent := request.Header.Get("User-Agent")

	err, errCode := w.onConnected(ws, roomName)
	if err != nil {
		http.Error(writer, err.Error(), errCode)
		return
	}

	w.startReadingMessages(ws, roomName, userAgent)
}

func (w *WSHandlerService) startReadingMessages(ws *websocket.Conn, roomName, userAgent string) {
	log := w.logger.Begin()
	defer log.End()

	defer func() {
		err := ws.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	Map, found := roommapssrv.RoomMaps.Get(roomName)
	if !found {
		log.Errorf("could not get map for room %q when trying to reader", roomName)
		return
	}
	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			_, closedError := err.(*websocket.CloseError)
			_, handshakeError := err.(*websocket.HandshakeError)
			if closedError || handshakeError {
				//receiver.deleteNode(ws, roomName, messageType)
				{
					err := roommapssrv.RoomMaps.Set(roomName, Map)
					if err != nil {
						log.Errorf("could not set room in map: %s", err)
					}
				}
			}
			return
		}

		var theMessage message.MessageContract
		err = json.Unmarshal(p, &theMessage)
		if err != nil {
			return
		}

		w.callHandler(Map.Get(ws).(*binarytreesrv.MySocket), theMessage.Type, theMessage, messageType)

		//receiver.parseMessage(Map.Get(ws).(*binarytreesrv.MySocket), p, messageType, userAgent, roomName)
	}
}

func (w *WSHandlerService) callHandler(socket *binarytreesrv.MySocket, cmd string, msg message.MessageContract, msgType int) {
	if handler, exists := w.handlers[cmd]; exists {
		handler(socket, msg, msgType)
	} else {
		w.handlers["*"](socket, msg, msgType)
	}
}
