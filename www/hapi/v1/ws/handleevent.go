package verboten

import (
	"strconv"

	"github.com/reiver/logjam/lib/rooms"
	"github.com/reiver/logjam/srv/room"
)

func handleEvent(ctx *rooms.WSContext) {
	if ctx.ParsedMessage == nil || len(ctx.PureMessage) <= 2 {
		return
	}
	if ctx.ParsedMessage.Type != "tree" && ctx.ParsedMessage.Type != "ping" && ctx.ParsedMessage.Type != "metadata-get" {
		log.Debugf("ID[%d] event: %s", ctx.SocketID, ctx.ParsedMessage.Type)
	}
	switch ctx.ParsedMessage.Type {
	case "start":
		roomsrv.Controller.Start(ctx)
	case "role":
		roomsrv.Controller.Role(ctx)
	case "stream":
		roomsrv.Controller.Stream(ctx)
	case "updateStreamId":
		roomsrv.Controller.UpdateStreamId(ctx)
	case "ping":
		roomsrv.Controller.Ping(ctx)
	case "turn_status":
		roomsrv.Controller.TurnStatus(ctx)
	case "tree":
		roomsrv.Controller.Tree(ctx)
	case "metadata-set":
		roomsrv.Controller.MetadataSet(ctx)
	case "metadata-get":
		roomsrv.Controller.MetadataGet(ctx)
	case "user-by-stream":
		roomsrv.Controller.UserByStream(ctx)
	case "muted":
		roomsrv.Controller.Muted(ctx)
	case "get-latest-user-list":
		roomsrv.Controller.GetLatestUserList(ctx)
	case "reconnect-children":
		roomsrv.Controller.ReconnectChildren(ctx)
	case "send-message":
		roomsrv.Controller.SendMessage(ctx)
	default:
		room, err := roomsrv.Repository.GetRoom(ctx.RoomId)
		if err != nil {
			log.Error(err)
		} else if room != nil {
			if room.GoldGorilla != nil {
				if ctx.ParsedMessage.Target == strconv.FormatUint((*room.GoldGorilla).ID, 10) {
					switch ctx.ParsedMessage.Type {
					case "video-answer":
						roomsrv.Controller.SendAnswerToAN(ctx)
					case "video-offer":
						roomsrv.Controller.SendOfferToAN(ctx)
					case "new-ice-candidate":
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
