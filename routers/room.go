package routers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sparkscience/logjam/controllers"
	"github.com/sparkscience/logjam/models"
	"github.com/sparkscience/logjam/models/contracts"
	"net/http"
	"strconv"
)

type roomWSRouter struct {
	roomCtrl  *controllers.RoomWSController
	roomRepo  contracts.IRoomRepository
	upgrader  websocket.Upgrader
	socketSVC contracts.ISocketService
	logger    contracts.ILogger
}

func newRoomWSRouter(roomCtrl *controllers.RoomWSController, roomRepo contracts.IRoomRepository, socketSVC contracts.ISocketService, logger contracts.ILogger) IRouteRegistrar {
	return &roomWSRouter{
		roomCtrl:  roomCtrl,
		roomRepo:  roomRepo,
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

func (r *roomWSRouter) registerRoutes(router *mux.Router) {
	router.HandleFunc("/ws", r.wsHandler)
}

func (r *roomWSRouter) wsHandler(writer http.ResponseWriter, request *http.Request) {
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

func (r *roomWSRouter) startReadingFromWS(wsConn *websocket.Conn, socketId uint64, roomId string) {
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

func (r *roomWSRouter) handleEvent(ctx *models.WSContext) {
	if ctx.ParsedMessage == nil || len(ctx.PureMessage) <= 2 {
		return
	}
	if ctx.ParsedMessage.Type != "tree" && ctx.ParsedMessage.Type != "ping" && ctx.ParsedMessage.Type != "metadata-get" {
		_ = r.logger.Log("ws_router", contracts.LDebug, "ID["+strconv.FormatUint(ctx.SocketID, 10)+"] event: "+ctx.ParsedMessage.Type)
	}
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

	case "muted":
		{
			r.roomCtrl.Muted(ctx)
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
	case "send-message":
		{
			r.roomCtrl.SendMessage(ctx)
			break
		}
	default:
		{
			room, err := r.roomRepo.GetRoom(ctx.RoomId)
			if err != nil {
				println(err.Error())
			} else if room != nil {
				if room.GoldGorilla != nil {
					if ctx.ParsedMessage.Target == strconv.FormatUint(models.GetGoldGorillaId(), 10) {
						switch ctx.ParsedMessage.Type {
						case "video-answer":
							{
								r.roomCtrl.SendAnswerToAN(ctx)
								break
							}
						case "video-offer":
							{
								r.roomCtrl.SendOfferToAN(ctx)
								break
							}
						case "new-ice-candidate":
							{
								r.roomCtrl.SendICECandidateToAN(ctx)
								break
							}
						default:
							{
								r.roomCtrl.DefaultHandler(ctx)
								break
							}
						}
					} else {
						r.roomCtrl.DefaultHandler(ctx)
					}
				} else {
					r.roomCtrl.DefaultHandler(ctx)
				}
			} else {
				r.roomCtrl.DefaultHandler(ctx)
			}
			break
		}
	}
}
