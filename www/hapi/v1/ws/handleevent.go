package verboten

import (
	"strconv"

	"github.com/reiver/go-erorr"

	"github.com/reiver/logjam/lib/rtc-rooms"
	"github.com/reiver/logjam/srv/rtc-rooms"
)

func handleEvent(ctx *rtc_rooms.WSContext) {
	if ctx.ParsedMessage == nil || len(ctx.PureMessage) <= 2 {
		return
	}
	if ctx.ParsedMessage.Type != "tree" && ctx.ParsedMessage.Type != "ping" && ctx.ParsedMessage.Type != "metadata-get" {
		log.Debugf("ID[%d] event: %s", ctx.SocketID, ctx.ParsedMessage.Type)
	}
	switch ctx.ParsedMessage.Type {
	case "start":
		rtcroomsrv.Controller.Start(ctx)
	case "role":
		rtcroomsrv.Controller.Role(ctx)
	case "stream":
		rtcroomsrv.Controller.Stream(ctx)
	case "updateStreamId":
		rtcroomsrv.Controller.UpdateStreamId(ctx)
	case "ping":
		rtcroomsrv.Controller.Ping(ctx)
	case "turn_status":
		rtcroomsrv.Controller.TurnStatus(ctx)
	case "tree":
		rtcroomsrv.Controller.Tree(ctx)
	case "metadata-set":
		rtcroomsrv.Controller.MetadataSet(ctx)
	case "metadata-get":
		rtcroomsrv.Controller.MetadataGet(ctx)
	case "user-by-stream":
		rtcroomsrv.Controller.UserByStream(ctx)
	case "muted":
		rtcroomsrv.Controller.Muted(ctx)
	case "get-latest-user-list":
		rtcroomsrv.Controller.GetLatestUserList(ctx)
	case "reconnect-children":
		rtcroomsrv.Controller.ReconnectChildren(ctx)
	case "send-message":
		rtcroomsrv.Controller.SendMessage(ctx)
	default:
		room, err := rtcroomsrv.Repository.GetRoom(ctx.RoomId)
		if err != nil && !erorr.Is(err, rtc_rooms.ErrRoomNotFound) {
			log.Error(err)
		} else if room != nil {
			if room.GoldGorilla != nil {
				if ctx.ParsedMessage.Target == strconv.FormatUint((*room.GoldGorilla).ID, 10) {
					switch ctx.ParsedMessage.Type {
					case "video-answer":
						rtcroomsrv.Controller.SendAnswerToAN(ctx)
					case "video-offer":
						rtcroomsrv.Controller.SendOfferToAN(ctx)
					case "new-ice-candidate":
						rtcroomsrv.Controller.SendICECandidateToAN(ctx)
					default:
						rtcroomsrv.Controller.DefaultHandler(ctx)
					}
				} else {
					rtcroomsrv.Controller.DefaultHandler(ctx)
				}
			} else {
				rtcroomsrv.Controller.DefaultHandler(ctx)
			}
		} else {
			rtcroomsrv.Controller.DefaultHandler(ctx)
		}
	}
}
