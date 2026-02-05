package verboten

import (
	"strconv"

	"github.com/reiver/go-erorr"

	"codeberg.org/greatape/logjam/lib/msgs"
	"codeberg.org/greatape/logjam/lib/rooms"
	"codeberg.org/greatape/logjam/srv/room"
)

func handleEvent(ctx *rooms.WSContext) {
	if ctx.ParsedMessage == nil || len(ctx.PureMessage) <= 2 {
		return
	}
	if ctx.ParsedMessage.Type != msgs.TypeTree && ctx.ParsedMessage.Type != msgs.TypePing && ctx.ParsedMessage.Type != msgs.TypeMetaDataGet {
		log.Debugf("ID[%d] event: %s", ctx.SocketID, ctx.ParsedMessage.Type)
	}
	switch ctx.ParsedMessage.Type {
	case msgs.TypeStart:
		roomsrv.Controller.Start(ctx)
	case msgs.TypeRole:
		roomsrv.Controller.Role(ctx)
	case msgs.TypeStream:
		roomsrv.Controller.Stream(ctx)
	case msgs.TypeUpdateStreamID:
		roomsrv.Controller.UpdateStreamId(ctx)
	case msgs.TypePing:
		roomsrv.Controller.Ping(ctx)
	case msgs.TypeTurnStatus:
		roomsrv.Controller.TurnStatus(ctx)
	case msgs.TypeTree:
		roomsrv.Controller.Tree(ctx)
	case msgs.TypeMetaDataSet:
		roomsrv.Controller.MetaDataSet(ctx)
	case msgs.TypeMetaDataGet:
		roomsrv.Controller.MetaDataGet(ctx)
	case msgs.TypeUserByStream:
		roomsrv.Controller.UserByStream(ctx)
	case msgs.TypeMuted:
		roomsrv.Controller.Muted(ctx)
	case msgs.TypeGetLatestUserList:
		roomsrv.Controller.GetLatestUserList(ctx)
	case msgs.TypeReconnectChildren:
		roomsrv.Controller.ReconnectChildren(ctx)
	case msgs.TypeSendMessage:
		roomsrv.Controller.SendMessage(ctx)
	case msgs.TypeLeave:
		roomsrv.Controller.Leave(ctx)
	default:
		room, err := roomsrv.Repository.GetRoom(ctx.RoomId)
		if err != nil && !erorr.Is(err, rooms.ErrRoomNotFound) {
			log.Error(err)
		} else if room != nil {
			if room.GoldGorilla != nil {
				if ctx.ParsedMessage.Target == strconv.FormatUint((*room.GoldGorilla).ID, 10) {
					switch ctx.ParsedMessage.Type {
					case msgs.TypeVideoAnswer:
						roomsrv.Controller.SendAnswerToAN(ctx)
					case msgs.TypeVideoOffer:
						roomsrv.Controller.SendOfferToAN(ctx)
					case msgs.TypeNewIceCandidate:
						roomsrv.Controller.SendICECandidateToAN(ctx)
					default:
						roomsrv.Controller.DefaultHandler(ctx)
					}
				} else {
					roomsrv.Controller.DefaultHandler(ctx)
				}
			} else {
				roomsrv.Controller.DefaultHandler(ctx)
			}
		} else {
			roomsrv.Controller.DefaultHandler(ctx)
		}
	}
}
