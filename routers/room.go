package routers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sparkscience/logjam/controllers"
	"github.com/sparkscience/logjam/models"
	"github.com/sparkscience/logjam/models/contracts"
	"net/http"
)

type RoomRouter struct {
	roomCtrl  *controllers.RoomController
	socketSVC contracts.ISocketService
	upgrader  websocket.Upgrader
	logger    contracts.ILogger
}

func NewRoomRouter(roomCtrl *controllers.RoomController, socketSVC contracts.ISocketService, logger contracts.ILogger) *RoomRouter {
	return &RoomRouter{
		roomCtrl:  roomCtrl,
		socketSVC: socketSVC,
		logger:    logger,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 4096,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (r *RoomRouter) Serve(listenHost string, listenPort int) error {
	httpMux := http.NewServeMux()
	addr := fmt.Sprintf(`%s:%d`, listenHost, listenPort)
	httpMux.HandleFunc("/ws", r.wsHandler)
	httpMux.Handle("/", http.FileServer(http.Dir("./views/")))
	println(fmt.Sprintf(`[HTTP] Listening on %s ..`, addr))
	err := http.ListenAndServe(addr, httpMux)
	if err != nil {
		return err
	}

	return nil
}

func (r *RoomRouter) wsHandler(writer http.ResponseWriter, request *http.Request) {
	wsConn, err := r.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		_ = r.logger.Log("ws_router", contracts.LError, err.Error())
		return
	}
	socketId, err := r.socketSVC.OnConnect(wsConn)
	if err != nil {
		_ = r.logger.Log("ws_router", contracts.LError, err.Error())
		_ = wsConn.Close()
		return
	}

	roomId := request.URL.Query().Get("room")
	go r.startReadingFromWS(wsConn, socketId, roomId)
}

func (r *RoomRouter) startReadingFromWS(wsConn *websocket.Conn, socketId uint64, roomId string) {
	for {
		messageType, data, readErr := wsConn.ReadMessage()
		if readErr != nil {
			_ = r.logger.Log("ws_router", contracts.LError, readErr.Error())
			go r.roomCtrl.OnDisconnect(&models.WSContext{
				RoomId:        roomId,
				SocketID:      socketId,
				PureMessage:   nil,
				ParsedMessage: nil,
			})
			break
		}
		if messageType != websocket.TextMessage {
			continue
		}

		var msg models.MessageContract
		err := json.Unmarshal(data, &msg)
		if err != nil {
			_ = r.logger.Log("ws_router", contracts.LError, err.Error())
			continue
		}

		ctx := &models.WSContext{
			RoomId:        roomId,
			SocketID:      socketId,
			PureMessage:   data,
			ParsedMessage: &msg,
		}
		go r.handleEvent(ctx)
	}
}

func (r *RoomRouter) handleEvent(ctx *models.WSContext) {
	if ctx.ParsedMessage == nil || len(ctx.PureMessage) <= 2 {
		return
	}
	_ = r.logger.Log("ws_router", contracts.LDebug, "event: "+ctx.ParsedMessage.Type)
	switch ctx.ParsedMessage.Type {
	case "start":
		{
			r.roomCtrl.Start(ctx)
			break
		}
	case "role":
		{
			r.roomCtrl.Role(ctx)
			break
		}
	case "stream":
		{
			r.roomCtrl.Stream(ctx)
			break
		}

	case "ping":
		{
			r.roomCtrl.Ping(ctx)
			break
		}

	case "turn_status":
		{
			r.roomCtrl.TurnStatus(ctx)
			break
		}
	case "tree":
		{
			r.roomCtrl.Tree(ctx)
			break
		}
	case "metadata-set":
		{
			r.roomCtrl.MetadataSet(ctx)
			break
		}
	case "metadata-get":
		{
			r.roomCtrl.MetadataGet(ctx)
			break
		}
	case "user-by-stream":
		{
			r.roomCtrl.UserByStream(ctx)
			break
		}

	case "get-latest-user-list":
		{
			r.roomCtrl.GetLatestUserList(ctx)
			break
		}

	case "reconnect-children":
		{
			r.roomCtrl.ReconnectChildren(ctx)
			break
		}
	default:
		{
			r.roomCtrl.DefaultHandler(ctx)
			break
		}
	}
}
