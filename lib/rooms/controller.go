package rooms

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/reiver/logjam/lib/goldgorilla"
	"github.com/reiver/logjam/lib/logs"
	"github.com/reiver/logjam/lib/msgs"
	"github.com/reiver/logjam/lib/websock"
)

type RoomWSController struct {
	logger    logs.Logger
	socketSVC websock.SocketService
	roomRepo  Repository
	ggRepo    goldgorilla.IGoldGorillaServiceRepository
}

func NewRoomWSController(socketSVC websock.SocketService, roomRepo Repository, ggRepo goldgorilla.IGoldGorillaServiceRepository, logger logs.TaggedLogger) *RoomWSController {
	const logtag string = "room_ws_ctrl"

	return &RoomWSController{
		logger:    logger.Tag(logtag),
		socketSVC: socketSVC,
		roomRepo:  roomRepo,
		ggRepo:    ggRepo,
	}
}

func (c *RoomWSController) OnConnect(ctx *WSContext) {

}

func (c *RoomWSController) OnDisconnect(ctx *WSContext) {
	defer c.emitUserList(ctx.RoomId)
	wasBroadcaster, childrenIdList, err := c.roomRepo.RemoveMember(ctx.RoomId, ctx.SocketID)
	if err != nil {
		c.error(err)
		return
	}
	if wasBroadcaster {
		err := c.roomRepo.ClearBroadcasterSeat(ctx.RoomId)
		if err != nil {
			c.error(err)
			return
		}
		membersIdList, err := c.roomRepo.GetAllMembersId(ctx.RoomId, true)
		if err != nil {
			c.error(err)
			return
		}
		brDCEvent := msgs.MessageContract{
			Type: "event-broadcaster-disconnected",
			Data: strconv.FormatUint(ctx.SocketID, 10),
		}
		_ = c.socketSVC.Send(brDCEvent, membersIdList...)
		oldggId, err := c.ggRepo.ResetRoom(ctx.RoomId)
		if err != nil {
			c.error(err)
			//return
		}
		if oldggId != nil {
			c.roomRepo.RemoveMember(ctx.RoomId, *oldggId)
		}
	} else {
		parentDCEvent := msgs.MessageContract{
			Type: "event-parent-dc",
			Data: strconv.FormatUint(ctx.SocketID, 10),
		}
		_ = c.socketSVC.Send(parentDCEvent, childrenIdList...)

		for _, id := range childrenIdList {
			if c.roomRepo.IsGGInstance(ctx.RoomId, id) {
				_, err := c.ggRepo.ResetRoom(ctx.RoomId)
				if err != nil {
					c.error(err)
					//return
				}
				c.roomRepo.RemoveMember(ctx.RoomId, id)
				go func() {
					err := c.ggRepo.Start(ctx.RoomId)
					if err != nil {
						c.error(err)
						return
					}
				}()
				break
			}
		}
	}
	membersIdList, err := c.roomRepo.GetAllMembersId(ctx.RoomId, false)
	if err != nil {
		c.error(err)
		return
	}
	if len(membersIdList) == 0 {
		err = c.roomRepo.ClearMessageHistory(ctx.RoomId)
		if err != nil {
			c.error(err)
		}
		err = c.roomRepo.SetRoomMetaData(ctx.RoomId, map[string]any{})
		if err != nil {
			c.error(err)
		}
	} else if len(membersIdList) == 1 {
		if c.roomRepo.IsGGInstance(ctx.RoomId, membersIdList[0]) {
			err = c.roomRepo.ClearMessageHistory(ctx.RoomId)
			if err != nil {
				c.error(err)
			}
			err = c.roomRepo.SetRoomMetaData(ctx.RoomId, map[string]any{})
			if err != nil {
				c.error(err)
			}
		}
	}
}

func (c *RoomWSController) Start(ctx *WSContext) {
	resultEvent := msgs.MessageContract{
		Type:   "start",
		Data:   strconv.FormatInt(int64(ctx.SocketID), 10),
		Target: "",
		Name:   "",
	}
	_ = c.roomRepo.CreateRoom(ctx.RoomId)
	err := c.roomRepo.AddMember(ctx.RoomId, ctx.SocketID, "", "", "", false)
	if err != nil {
		c.error(err)
		return
	}
	err = c.roomRepo.UpdateMemberName(ctx.RoomId, ctx.SocketID, ctx.ParsedMessage.Data)
	if err != nil {
		c.error(err)
	}
	_ = c.socketSVC.Send(resultEvent, ctx.SocketID)
}
func (c *RoomWSController) Role(ctx *WSContext) {
	var eventData map[string]any
	err := json.Unmarshal(ctx.PureMessage, &eventData)
	if err != nil {

		return
	}
	streamId, exists := eventData["streamId"]
	defer c.emitUserList(ctx.RoomId)
	if exists {
		err = c.roomRepo.UpdateMemberMeta(ctx.RoomId, ctx.SocketID, "streamId", streamId.(string))
		if err != nil {
			c.error(err)
			return
		}
	}

	resultEvent := msgs.MessageContract{
		Type:   "role",
		Data:   "",
		Target: "",
		Name:   "",
	}
	if ctx.ParsedMessage.Data == "broadcast" {
		br, err := c.roomRepo.GetBroadcaster(ctx.RoomId)
		if err != nil {
			c.error(err)
			return
		}
		if br != nil {
			currentUser, err := c.roomRepo.GetMember(ctx.RoomId, ctx.SocketID)
			if err != nil {
				c.error(err)
				_ = c.socketSVC.Send(msgs.MessageContract{
					Type: "role",
					Data: "no:broadcast",
				}, ctx.SocketID)
				return
			}
			if currentUser.Name != br.Name {
				_ = c.socketSVC.Send(msgs.MessageContract{
					Type: "role",
					Data: "no:broadcast",
				}, ctx.SocketID)
				return
			}
			go c.socketSVC.Disconnect(br.ID)
		}
		resultEvent.Data = "yes:broadcast"
		err = c.roomRepo.UpdateCanConnect(ctx.RoomId, ctx.SocketID, true)
		if err != nil {
			c.error(err)
			return
		}
		err = c.roomRepo.SetBroadcaster(ctx.RoomId, ctx.SocketID)
		if err != nil {
			c.error(err)
			return
		}
		memberIds, err := c.roomRepo.GetAllMembersId(ctx.RoomId, true)
		if err != nil {
			c.error(err)
			return
		}

		err = c.socketSVC.Send(msgs.MessageContract{
			Type: "broadcasting",
			Data: strconv.FormatUint(ctx.SocketID, 10),
		}, memberIds...)
		if err != nil {
			c.error(err)
		}
		_ = c.socketSVC.Send(resultEvent, ctx.SocketID)
		ggEnabled := true
		ggEnabledInReqBody, exists := eventData["ggEnabled"]
		if exists && ggEnabledInReqBody == false {
			ggEnabled = false
		}
		if ggEnabled == true {
			go func() {
				err := c.ggRepo.Start(ctx.RoomId)
				if err != nil {
					c.error(err)
					return
				}
			}()
		}
	} else if ctx.ParsedMessage.Data == "alt-broadcast" {
		broadcaster, err := c.roomRepo.GetBroadcaster(ctx.RoomId)
		if err != nil {
			c.error(err)
			return
		}
		if broadcaster == nil {
			_ = c.socketSVC.Send(msgs.MessageContract{
				Type: "alt-broadcast",
				Data: "no-broadcaster",
			}, ctx.SocketID)
			return
		}

		resultEvent.Type = "alt-broadcast"
		resultEvent.Data = strconv.FormatUint(broadcaster.ID, 10)
		_ = c.socketSVC.Send(resultEvent, ctx.SocketID)

		userInfo, err := c.roomRepo.GetMember(ctx.RoomId, ctx.SocketID)
		if err != nil {
			c.error(err)
			return
		}
		broadcasterReceivingEvent := map[string]string{
			"Type": `alt-broadcast`,
			"Data": strconv.FormatUint(ctx.SocketID, 10),
			"name": userInfo.Name,
		}
		_ = c.socketSVC.Send(broadcasterReceivingEvent, broadcaster.ID)
	} else if ctx.ParsedMessage.Data == "audience" {
		broadcaster, err := c.roomRepo.GetBroadcaster(ctx.RoomId)
		if err != nil {
			c.error(err)
			return
		}
		if broadcaster == nil {
			_ = c.socketSVC.Send(msgs.MessageContract{
				Type: "role",
				Data: "no:audience",
			}, ctx.SocketID)
			return
		}
		tryCount := 0
	start:
		parentId, err := c.roomRepo.InsertMemberToTree(ctx.RoomId, ctx.SocketID, false)
		if err != nil && tryCount <= 20 {
			time.Sleep(500 * time.Millisecond)
			tryCount++
			goto start
		}
		if err != nil {
			c.error(err)
			_ = c.socketSVC.Send(map[string]any{
				"type": "error",
				"data": "Insert Child Error : " + err.Error(),
			}, ctx.SocketID)
			return
		}
		if !c.roomRepo.IsGGInstance(ctx.RoomId, *parentId) {
			_ = c.socketSVC.Send(msgs.MessageContract{
				Type: "add_audience",
				Data: strconv.FormatUint(ctx.SocketID, 10),
			}, *parentId)
		} else {
			ggid, err := c.roomRepo.GetRoomGoldGorillaId(ctx.RoomId)
			if err != nil {
				c.error(err)
				return
			}
			if ggid == nil {
				c.socketSVC.Disconnect(ctx.SocketID) // so aud will reconnect
				go c.roomRepo.RemoveMember(ctx.RoomId, *parentId)
				return
			}
			err = c.ggRepo.CreatePeer(ctx.RoomId, ctx.SocketID, true, false, *ggid)
			if err != nil {
				_ = c.socketSVC.Send(msgs.MessageContract{
					Type: "error",
					Data: err.Error(),
				}, ctx.SocketID)
			}
		}
		go c.emitUserList(ctx.RoomId)
	}
}

func (c *RoomWSController) Stream(ctx *WSContext) {
	payload := make(map[string]string)
	err := json.Unmarshal(ctx.PureMessage, &payload)
	if err != nil {
		c.error(err)
		return
	}
	newState := true
	if data, exists := payload["data"]; exists {
		if data == "true" {
			newState = true
		} else {
			newState = false
		}
	}
	err = c.roomRepo.UpdateCanConnect(ctx.RoomId, ctx.SocketID, newState)
	if err != nil {
		c.error(err)
		return
	}
}

func (c *RoomWSController) UpdateStreamId(ctx *WSContext) {
	payload := make(map[string]any)
	err := json.Unmarshal(ctx.PureMessage, &payload)
	if err != nil {
		c.error(err)
		return
	}
	streamId, exists := payload["streamId"]
	defer c.emitUserList(ctx.RoomId)
	if exists {
		err = c.roomRepo.UpdateMemberMeta(ctx.RoomId, ctx.SocketID, "streamId", streamId.(string))
		if err != nil {
			c.error(err)
			return
		}
	}
}

func (c *RoomWSController) Ping(ctx *WSContext) {
	_ = c.socketSVC.Send(msgs.MessageContract{Type: "pong"}, ctx.SocketID)
}

func (c *RoomWSController) TurnStatus(ctx *WSContext) {
	err := c.roomRepo.UpdateTurnStatus(ctx.RoomId, ctx.SocketID, ctx.ParsedMessage.Data == "on")
	if err != nil {
		c.error(err)
		return
	}
}

func (c *RoomWSController) Tree(ctx *WSContext) {
	room, err := c.roomRepo.GetRoom(ctx.RoomId)
	if err != nil {
		c.error(err)
		return
	}
	if room == nil {
		c.error(ErrRoomNotFound)
		return
	}
	tree, err := room.GetTree()
	if err != nil {
		c.error(err)
		return
	}
	buffer, err := json.Marshal(tree)
	resultEvent := msgs.MessageContract{
		Type: "tree",
		Data: string(buffer),
	}
	_ = c.socketSVC.Send(resultEvent, ctx.SocketID)
}

func (c *RoomWSController) MetadataSet(ctx *WSContext) {
	metaData := make(map[string]any)
	err := json.Unmarshal([]byte(ctx.ParsedMessage.Data), &metaData)
	if err != nil {
		c.error(err)
		return
	}
	if _, exists := metaData[RoomMessagesMetaDataKey]; exists {
		c.socketSVC.Send(map[string]string{"error": "can't overwrite message history"}, ctx.SocketID)
		return
	}
	err = c.roomRepo.SetRoomMetaData(ctx.RoomId, metaData)
	if err != nil {
		c.error(err)
		return
	}
}

func (c *RoomWSController) MetadataGet(ctx *WSContext) {
	meta, err := c.roomRepo.GetRoomMetaData(ctx.RoomId)
	if err != nil {
		c.error(err)
		return
	}
	jsonBytes, err := json.Marshal(meta)
	if err != nil {
		c.error(err)
		return
	}
	resultEvent := msgs.MessageContract{
		Type: "metadata-get",
		Data: string(jsonBytes),
	}
	_ = c.socketSVC.Send(resultEvent, ctx.SocketID)
}

func (c *RoomWSController) UserByStream(ctx *WSContext) {
	userInfo, err := c.roomRepo.GetUserByStreamId(ctx.RoomId, ctx.ParsedMessage.Data)
	if err != nil {
		c.error(err)
		return
	}
	if userInfo == nil {
		return
	}
	isBroadcaster, err := c.roomRepo.IsBroadcaster(ctx.RoomId, userInfo.ID)
	if err != nil {
		c.error(err)
	}
	resultEvent := msgs.MessageContract{
		Type: "user-by-stream",
	}
	userRole := "audience"
	if isBroadcaster {
		userRole = "broadcast"
	}
	resultEvent.Data = strconv.FormatUint(userInfo.ID, 10) + "," + userInfo.Name + "," + ctx.ParsedMessage.Data + "," + userRole
	_ = c.socketSVC.Send(ctx.SocketID)
}

func (c *RoomWSController) GetLatestUserList(ctx *WSContext) {
	c.emitUserList(ctx.RoomId)
}

func (c *RoomWSController) Muted(ctx *WSContext) {
	list, err := c.roomRepo.GetAllMembersId(ctx.RoomId, false)
	if err != nil {
		c.error(err)
		return
	}
	err = c.socketSVC.Send(ctx.PureMessage, list...)
}

func (c *RoomWSController) emitUserList(roomId string) {
	list, err := c.roomRepo.GetMembersList(roomId)
	if err != nil {
		c.error(err)
		return
	}
	roomMembersIdList, err := c.roomRepo.GetAllMembersId(roomId, false)
	if err != nil {
		c.error(err)
		return
	}
	index := -1
	for i, v := range list {
		if c.roomRepo.IsGGInstance(roomId, v.Id) {
			index = i
			break
		}
	}
	if index > -1 {
		list = append(list[:index], list[index+1:]...)
	}
	buffer, err := json.Marshal(list)
	if err != nil {
		c.error(err)
		return
	}
	event := msgs.MessageContract{
		Type: "user-event",
		Data: string(buffer),
	}
	_ = c.socketSVC.Send(event, roomMembersIdList...)
}

func (c *RoomWSController) ReconnectChildren(ctx *WSContext) {
	childrenIdList, err := c.roomRepo.GetChildrenIdList(ctx.RoomId, ctx.SocketID)
	if err != nil {
		c.error(err)
		return
	}
	event := msgs.MessageContract{
		Type: "reconnect",
		Data: strconv.FormatUint(ctx.SocketID, 10),
	}
	_ = c.socketSVC.Send(event, childrenIdList...)
}

func (c *RoomWSController) SendMessage(ctx *WSContext) {
	membersIdList, err := c.roomRepo.GetAllMembersId(ctx.RoomId, false)
	if err != nil {
		c.error(err)
		_ = c.socketSVC.Send(map[string]string{"error": "error getting members list"}, ctx.SocketID)
		return
	}
	_ = c.socketSVC.Send(msgs.MessageContract{
		Type:   "new-message",
		Data:   ctx.ParsedMessage.Data,
		Target: "",
		Name:   strconv.FormatUint(ctx.SocketID, 10),
	}, membersIdList...)
	err = c.roomRepo.AddMessageToHistory(ctx.RoomId, ctx.SocketID, ctx.ParsedMessage.Data)
	if err != nil {
		c.error(err)
		return
	}

}

func (c *RoomWSController) SendOfferToAN(ctx *WSContext) {
	msg := make(map[string]any)
	err := json.Unmarshal(ctx.PureMessage, &msg)
	if err != nil {
		c.error(err)
		return
	}
	err = c.ggRepo.SendOffer(ctx.RoomId, ctx.SocketID, msg["sdp"])
	if err != nil {
		c.error(err)
		return
	}
}

func (c *RoomWSController) SendAnswerToAN(ctx *WSContext) {
	msg := make(map[string]any)
	err := json.Unmarshal(ctx.PureMessage, &msg)
	if err != nil {
		c.error(err)
		return
	}
	err = c.ggRepo.SendAnswer(ctx.RoomId, ctx.SocketID, msg["sdp"])
	if err != nil {
		c.error(err)
		return
	}
}

func (c *RoomWSController) SendICECandidateToAN(ctx *WSContext) {
	msg := make(map[string]any)
	err := json.Unmarshal(ctx.PureMessage, &msg)
	if err != nil {
		c.error(err)
		return
	}
	err = c.ggRepo.SendICECandidate(ctx.RoomId, ctx.SocketID, msg["candidate"])
	if err != nil {
		c.error(err)
		return
	}
}

func (c *RoomWSController) DefaultHandler(ctx *WSContext) {
	id, err := strconv.ParseUint(ctx.ParsedMessage.Target, 10, 64)
	if err != nil {
		c.error(err)
		return
	}
	if c.roomRepo.IsGGInstance(ctx.RoomId, id) {
		return // as there is no GoldGorilla in tree(as a browser user!!), we ignore messages that targets it
	}
	targetMember, err := c.roomRepo.GetMember(ctx.RoomId, id)
	if err != nil {
		c.error(err)
		return
	}

	if targetMember == nil {
		//ignoring
		return
	}
	var fullMessage map[string]any
	err = json.Unmarshal(ctx.PureMessage, &fullMessage)
	if err != nil {
		c.error(err)
		return
	}
	userInfo, err := c.roomRepo.GetMember(ctx.RoomId, ctx.SocketID)
	if err != nil {
		c.error(err)
		return
	}
	if userInfo == nil {
		//ignoring
		return
	}
	fullMessage["username"] = userInfo.Name
	fullMessage["data"] = strconv.FormatUint(ctx.SocketID, 10)
	_ = c.socketSVC.Send(fullMessage, targetMember.ID)
}

func (c *RoomWSController) debug(msg ...any) {
	 c.logger.Debug(msg...)
}

func (c *RoomWSController) error(msg ...any) {
	c.logger.Error(msg...)
}

func (c *RoomWSController) info(msg ...any) {
	c.logger.Info(msg...)
}
