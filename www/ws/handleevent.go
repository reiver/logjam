package verboten

import (
	"strconv"

	"github.com/reiver/logjam/models"
	"github.com/reiver/logjam/srv/room"
)

func handleEvent(ctx *models.WSContext) {
	if ctx.ParsedMessage == nil || len(ctx.PureMessage) <= 2 {
		return
	}
	if ctx.ParsedMessage.Type != "tree" && ctx.ParsedMessage.Type != "ping" && ctx.ParsedMessage.Type != "metadata-get" {
		log.Debugf("ID[%d] event: %s", ctx.SocketID, ctx.ParsedMessage.Type)
	}
	switch ctx.ParsedMessage.Type {
	case "start":
		{
			roomsrv.Controller.Start(ctx)
			break
		}
	case "role":
		{
			roomsrv.Controller.Role(ctx)
			break
		}
	case "stream":
		{
			roomsrv.Controller.Stream(ctx)
			break
		}
	case "updateStreamId":
		{
			roomsrv.Controller.UpdateStreamId(ctx)
			break
		}
	case "ping":
		{
			roomsrv.Controller.Ping(ctx)
			break
		}

	case "turn_status":
		{
			roomsrv.Controller.TurnStatus(ctx)
			break
		}
	case "tree":
		{
			roomsrv.Controller.Tree(ctx)
			break
		}
	case "metadata-set":
		{
			roomsrv.Controller.MetadataSet(ctx)
			break
		}
	case "metadata-get":
		{
			roomsrv.Controller.MetadataGet(ctx)
			break
		}
	case "user-by-stream":
		{
			roomsrv.Controller.UserByStream(ctx)
			break
		}

	case "muted":
		{
			roomsrv.Controller.Muted(ctx)
			break
		}
	case "get-latest-user-list":
		{
			roomsrv.Controller.GetLatestUserList(ctx)
			break
		}

	case "reconnect-children":
		{
			roomsrv.Controller.ReconnectChildren(ctx)
			break
		}
	case "send-message":
		{
			roomsrv.Controller.SendMessage(ctx)
			break
		}
	default:
		{
			room, err := roomsrv.Repository.GetRoom(ctx.RoomId)
			if err != nil {
				log.Error(err)
			} else if room != nil {
				if room.GoldGorilla != nil {
					if ctx.ParsedMessage.Target == strconv.FormatUint((*room.GoldGorilla).ID, 10) {
						switch ctx.ParsedMessage.Type {
						case "video-answer":
							{
								roomsrv.Controller.SendAnswerToAN(ctx)
								break
							}
						case "video-offer":
							{
								roomsrv.Controller.SendOfferToAN(ctx)
								break
							}
						case "new-ice-candidate":
							{
								roomsrv.Controller.SendICECandidateToAN(ctx)
								break
							}
						default:
							{
								roomsrv.Controller.DefaultHandler(ctx)
								break
							}
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
			break
		}
	}
}
