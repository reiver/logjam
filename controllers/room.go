package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/reiver/logjam/lib/logs"
	"github.com/reiver/logjam/models"
	"github.com/reiver/logjam/models/contracts"
)

type RoomWSController struct {
	logger    logs.TaggedLogger
	socketSVC contracts.ISocketService
	roomRepo  contracts.IRoomRepository
	ggRepo    contracts.IGoldGorillaServiceRepository
}

func NewRoomWSController(socketSVC contracts.ISocketService, roomRepo contracts.IRoomRepository, ggRepo contracts.IGoldGorillaServiceRepository, logger logs.TaggedLogger) *RoomWSController {
	return &RoomWSController{
		logger:    logger,
		socketSVC: socketSVC,
		roomRepo:  roomRepo,
		ggRepo:    ggRepo,
	}
}

func (c *RoomWSController) OnConnect(ctx *models.WSContext) {

}

func (c *RoomWSController) OnDisconnect(ctx *models.WSContext) {
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
		brDCEvent := models.MessageContract{
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
		parentDCEvent := models.MessageContract{
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
		err = c.roomRepo.SetRoomMetaData(ctx.RoomId, map[string]interface{}{})
		if err != nil {
			c.error(err)
		}
	} else if len(membersIdList) == 1 {
		if c.roomRepo.IsGGInstance(ctx.RoomId, membersIdList[0]) {
			err = c.roomRepo.ClearMessageHistory(ctx.RoomId)
			if err != nil {
				c.error(err)
			}
			err = c.roomRepo.SetRoomMetaData(ctx.RoomId, map[string]interface{}{})
			if err != nil {
				c.error(err)
			}
		}
	}
}

func (c *RoomWSController) Start(ctx *models.WSContext) {
	resultEvent := models.MessageContract{
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
func (c *RoomWSController) Role(ctx *models.WSContext) {
	var eventData map[string]interface{}
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

	resultEvent := models.MessageContract{
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
				_ = c.socketSVC.Send(models.MessageContract{
					Type: "role",
					Data: "no:broadcast",
				}, ctx.SocketID)
				return
			}
			if currentUser.Name != br.Name {
				_ = c.socketSVC.Send(models.MessageContract{
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

		err = c.socketSVC.Send(models.MessageContract{
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
			_ = c.socketSVC.Send(models.MessageContract{
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
			_ = c.socketSVC.Send(models.MessageContract{
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
			_ = c.socketSVC.Send(map[string]interface{}{
				"type": "error",
				"data": "Insert Child Error : " + err.Error(),
			}, ctx.SocketID)
			return
		}
		if !c.roomRepo.IsGGInstance(ctx.RoomId, *parentId) {
			_ = c.socketSVC.Send(models.MessageContract{
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
				_ = c.socketSVC.Send(models.MessageContract{
					Type: "error",
					Data: err.Error(),
				}, ctx.SocketID)
			}
		}
		go c.emitUserList(ctx.RoomId)
	}
}

func (c *RoomWSController) Stream(ctx *models.WSContext) {
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

func (c *RoomWSController) UpdateStreamId(ctx *models.WSContext) {
	payload := make(map[string]interface{})
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

func (c *RoomWSController) Ping(ctx *models.WSContext) {
	_ = c.socketSVC.Send(models.MessageContract{Type: "pong"}, ctx.SocketID)
}

func (c *RoomWSController) TurnStatus(ctx *models.WSContext) {
	err := c.roomRepo.UpdateTurnStatus(ctx.RoomId, ctx.SocketID, ctx.ParsedMessage.Data == "on")
	if err != nil {
		c.error(err)
		return
	}
}

func (c *RoomWSController) Tree(ctx *models.WSContext) {
	room, err := c.roomRepo.GetRoom(ctx.RoomId)
	if err != nil {
		c.error(err)
		return
	}
	if room == nil {
		c.error("room doesn't exists")
		return
	}
	tree, err := room.GetTree()
	if err != nil {
		c.error(err)
		return
	}
	buffer, err := json.Marshal(tree)
	resultEvent := models.MessageContract{
		Type: "tree",
		Data: string(buffer),
	}
	_ = c.socketSVC.Send(resultEvent, ctx.SocketID)
}

func (c *RoomWSController) MetadataSet(ctx *models.WSContext) {
	metaData := make(map[string]interface{})
	err := json.Unmarshal([]byte(ctx.ParsedMessage.Data), &metaData)
	if err != nil {
		c.error(err)
		return
	}
	if _, exists := metaData[models.RoomMessagesMetaDataKey]; exists {
		c.socketSVC.Send(map[string]string{"error": "can't overwrite message history"}, ctx.SocketID)
		return
	}
	err = c.roomRepo.SetRoomMetaData(ctx.RoomId, metaData)
	if err != nil {
		c.error(err)
		return
	}
}

func (c *RoomWSController) MetadataGet(ctx *models.WSContext) {
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
	resultEvent := models.MessageContract{
		Type: "metadata-get",
		Data: string(jsonBytes),
	}
	_ = c.socketSVC.Send(resultEvent, ctx.SocketID)
}

func (c *RoomWSController) UserByStream(ctx *models.WSContext) {
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
	resultEvent := models.MessageContract{
		Type: "user-by-stream",
	}
	userRole := "audience"
	if isBroadcaster {
		userRole = "broadcast"
	}
	resultEvent.Data = strconv.FormatUint(userInfo.ID, 10) + "," + userInfo.Name + "," + ctx.ParsedMessage.Data + "," + userRole
	_ = c.socketSVC.Send(ctx.SocketID)
}

func (c *RoomWSController) GetLatestUserList(ctx *models.WSContext) {
	c.emitUserList(ctx.RoomId)
}

func (c *RoomWSController) Muted(ctx *models.WSContext) {
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
	event := models.MessageContract{
		Type: "user-event",
		Data: string(buffer),
	}
	_ = c.socketSVC.Send(event, roomMembersIdList...)
}

func (c *RoomWSController) ReconnectChildren(ctx *models.WSContext) {
	childrenIdList, err := c.roomRepo.GetChildrenIdList(ctx.RoomId, ctx.SocketID)
	if err != nil {
		c.error(err)
		return
	}
	event := models.MessageContract{
		Type: "reconnect",
		Data: strconv.FormatUint(ctx.SocketID, 10),
	}
	_ = c.socketSVC.Send(event, childrenIdList...)
}

func (c *RoomWSController) SendMessage(ctx *models.WSContext) {
	membersIdList, err := c.roomRepo.GetAllMembersId(ctx.RoomId, false)
	if err != nil {
		c.error(err)
		_ = c.socketSVC.Send(map[string]string{"error": "error getting members list"}, ctx.SocketID)
		return
	}
	_ = c.socketSVC.Send(models.MessageContract{
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

func (c *RoomWSController) SendOfferToAN(ctx *models.WSContext) {
	msg := make(map[string]interface{})
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

func (c *RoomWSController) SendAnswerToAN(ctx *models.WSContext) {
	msg := make(map[string]interface{})
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

func (c *RoomWSController) SendICECandidateToAN(ctx *models.WSContext) {
	msg := make(map[string]interface{})
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

func (c *RoomWSController) DefaultHandler(ctx *models.WSContext) {
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
	var fullMessage map[string]interface{}
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
	 c.logger.Debug("room_ws_ctrl", msg...)
}

func (c *RoomWSController) error(msg ...any) {
	c.logger.Error("room_ws_ctrl", msg...)
}

func (c *RoomWSController) info(msg ...any) {
	c.logger.Info("room_ws_ctrl", msg...)
}
