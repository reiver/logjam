package routers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/reiver/logjam/controllers"
	"github.com/reiver/logjam/lib/logs"
	"github.com/reiver/logjam/lib/rooms"
	"github.com/reiver/logjam/models"
	"github.com/reiver/logjam/models/contracts"
)

type roomWSRouter struct {
	roomCtrl  *controllers.RoomWSController
	roomRepo  rooms.Repository
	upgrader  websocket.Upgrader
	socketSVC contracts.ISocketService
	logger    logs.Logger
}

func newRoomWSRouter(roomCtrl *controllers.RoomWSController, roomRepo rooms.Repository, socketSVC contracts.ISocketService, logger logs.TaggedLogger) IRouteRegistrar {
	const logtag string = "ws_router"

	return &roomWSRouter{
		roomCtrl:  roomCtrl,
		roomRepo:  roomRepo,
		socketSVC: socketSVC,
		logger:    logger.Tag(logtag),
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
		r.logger.Error(err)
		return
	}
	socketId, err := r.socketSVC.OnConnect(wsConn)
	if err != nil {
		r.logger.Error(err)
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
			r.logger.Error(readErr)
			go r.roomCtrl.OnDisconnect(&models.WSContext{
				RoomId:        roomId,
				SocketID:      socketId,
				PureMessage:   nil,
				ParsedMessage: nil,
			})
			_ = wsConn.CloseHandler()(1001, readErr.Error())
			break
		}
		if messageType != websocket.TextMessage {
			r.logger.Debugf("ignoring a message of type: %d", messageType)
			continue
		}

		var msg models.MessageContract
		err := json.Unmarshal(data, &msg)
		if err != nil {
			r.logger.Error(err)
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
		r.logger.Debugf("ID[%d] event: %s", ctx.SocketID, ctx.ParsedMessage.Type)
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
	case "updateStreamId":
		{
			r.roomCtrl.UpdateStreamId(ctx)
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
				r.logger.Error(err)
			} else if room != nil {
				if room.GoldGorilla != nil {
					if ctx.ParsedMessage.Target == strconv.FormatUint((*room.GoldGorilla).ID, 10) {
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
