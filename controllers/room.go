package controllers

import (
	"encoding/json"
	"github.com/sparkscience/logjam/models"
	"github.com/sparkscience/logjam/models/contracts"
	"strconv"
	"time"
)

type RoomWSController struct {
	logger    contracts.ILogger
	socketSVC contracts.ISocketService
	roomRepo  contracts.IRoomRepository
}

func NewRoomWSController(socketSVC contracts.ISocketService, roomRepo contracts.IRoomRepository, logger contracts.ILogger) *RoomWSController {
	return &RoomWSController{
		logger:    logger,
		socketSVC: socketSVC,
		roomRepo:  roomRepo,
	}
}

func (c *RoomWSController) OnConnect(ctx *models.WSContext) {

}

func (c *RoomWSController) OnDisconnect(ctx *models.WSContext) {
	defer c.emitUserList(ctx.RoomId)
	wasBroadcaster, childrenIdList, err := c.roomRepo.RemoveMember(ctx.RoomId, ctx.SocketID)
	if err != nil {
		c.log(contracts.LError, err.Error())
		return
	}
	if wasBroadcaster {
		err := c.roomRepo.ClearBroadcasterSeat(ctx.RoomId)
		if err != nil {
			c.log(contracts.LError, err.Error())
			return
		}
		membersIdList, err := c.roomRepo.GetAllMembersId(ctx.RoomId, true)
		if err != nil {
			c.log(contracts.LError, err.Error())
			return
		}
		brDCEvent := models.MessageContract{
			Type: "event-broadcaster-disconnected",
			Data: strconv.FormatUint(ctx.SocketID, 10),
		}
		_ = c.socketSVC.Send(brDCEvent, membersIdList...)
	} else {
		parentDCEvent := models.MessageContract{
			Type: "event-parent-dc",
			Data: strconv.FormatUint(ctx.SocketID, 10),
		}
		_ = c.socketSVC.Send(parentDCEvent, childrenIdList...)
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
	err := c.roomRepo.AddMember(ctx.RoomId, ctx.SocketID, "", "", "")
	if err != nil {
		c.log(contracts.LError, err.Error())
		return
	}
	err = c.roomRepo.UpdateMemberName(ctx.RoomId, ctx.SocketID, ctx.ParsedMessage.Data)
	if err != nil {
		c.log(contracts.LError, err.Error())
	}
	_ = c.socketSVC.Send(resultEvent, ctx.SocketID)
}
func (c *RoomWSController) Role(ctx *models.WSContext) {
	var eventData map[string]string
	err := json.Unmarshal(ctx.PureMessage, &eventData)
	if err != nil {

		return
	}
	streamId, exists := eventData["streamId"]
	defer c.emitUserList(ctx.RoomId)
	if exists {
		err = c.roomRepo.UpdateMemberMeta(ctx.RoomId, ctx.SocketID, "streamId", streamId)
		if err != nil {
			c.log(contracts.LError, err.Error())
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
		resultEvent.Data = "yes:broadcast"
		err := c.roomRepo.UpdateCanConnect(ctx.RoomId, ctx.SocketID, true)
		if err != nil {
			c.log(contracts.LError, err.Error())
			return
		}
		err = c.roomRepo.SetBroadcaster(ctx.RoomId, ctx.SocketID)
		if err != nil {
			c.log(contracts.LError, err.Error())
			return
		}
		memberIds, err := c.roomRepo.GetAllMembersId(ctx.RoomId, true)
		if err != nil {
			c.log(contracts.LError, err.Error())
			return
		}
		err = c.socketSVC.Send(models.MessageContract{
			Type: "broadcasting",
			Data: strconv.FormatUint(ctx.SocketID, 10),
		}, memberIds...)
		if err != nil {
			c.log(contracts.LError, err.Error())
		}
		_ = c.socketSVC.Send(resultEvent, ctx.SocketID)
	} else if ctx.ParsedMessage.Data == "alt-broadcast" {
		broadcaster, err := c.roomRepo.GetBroadcaster(ctx.RoomId)
		if err != nil {
			c.log(contracts.LError, err.Error())
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
			c.log(contracts.LError, err.Error())
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
			c.log(contracts.LError, err.Error())
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
			c.log(contracts.LError, err.Error())
			_ = c.socketSVC.Send(map[string]any{
				"type": "error",
				"data": "Insert Child Error : " + err.Error(),
			}, ctx.SocketID)
			return
		}
		_ = c.socketSVC.Send(models.MessageContract{
			Type: "add_audience",
			Data: strconv.FormatUint(ctx.SocketID, 10),
		}, *parentId)
	}
}
func (c *RoomWSController) Stream(ctx *models.WSContext) {
	payload := make(map[string]string)
	err := json.Unmarshal(ctx.PureMessage, &payload)
	if err != nil {
		c.log(contracts.LError, err.Error())
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
		c.log(contracts.LError, err.Error())
		return
	}
}

func (c *RoomWSController) Ping(ctx *models.WSContext) {
	_ = c.socketSVC.Send(models.MessageContract{Type: "pong"}, ctx.SocketID)
}

func (c *RoomWSController) TurnStatus(ctx *models.WSContext) {
	err := c.roomRepo.UpdateTurnStatus(ctx.RoomId, ctx.SocketID, ctx.ParsedMessage.Data == "on")
	if err != nil {
		c.log(contracts.LError, err.Error())
		return
	}
}

func (c *RoomWSController) Tree(ctx *models.WSContext) {
	room, err := c.roomRepo.GetRoom(ctx.RoomId)
	if err != nil {
		c.log(contracts.LError, err.Error())
		return
	}
	if room == nil {
		println("room doesn't exists")
		return
	}
	tree, err := room.GetTree()
	if err != nil {
		c.log(contracts.LError, err.Error())
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
	metaData := make(map[string]any)
	err := json.Unmarshal([]byte(ctx.ParsedMessage.Data), &metaData)
	if err != nil {
		c.log(contracts.LError, err.Error())
		return
	}
	err = c.roomRepo.SetRoomMetaData(ctx.RoomId, metaData)
	if err != nil {
		c.log(contracts.LError, err.Error())
		return
	}
}

func (c *RoomWSController) MetadataGet(ctx *models.WSContext) {
	meta, err := c.roomRepo.GetRoomMetaData(ctx.RoomId)
	if err != nil {
		c.log(contracts.LError, err.Error())
		return
	}
	_ = c.socketSVC.Send(meta, ctx.SocketID)
}

func (c *RoomWSController) UserByStream(ctx *models.WSContext) {
	userInfo, err := c.roomRepo.GetUserByStreamId(ctx.RoomId, ctx.ParsedMessage.Data)
	if err != nil {
		c.log(contracts.LError, err.Error())
		return
	}
	if userInfo == nil {
		return
	}
	isBroadcaster, err := c.roomRepo.IsBroadcaster(ctx.RoomId, userInfo.ID)
	if err != nil {
		c.log(contracts.LError, err.Error())
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

func (c *RoomWSController) emitUserList(roomId string) {
	list, err := c.roomRepo.GetMembersList(roomId)
	if err != nil {
		c.log(contracts.LError, err.Error())
		return
	}
	roomMembersIdList, err := c.roomRepo.GetAllMembersId(roomId, false)
	if err != nil {
		c.log(contracts.LError, err.Error())
		return
	}
	buffer, err := json.Marshal(list)
	if err != nil {
		c.log(contracts.LError, err.Error())
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
		c.log(contracts.LError, err.Error())
		return
	}
	event := models.MessageContract{
		Type: "reconnect",
		Data: strconv.FormatUint(ctx.SocketID, 10),
	}
	_ = c.socketSVC.Send(event, childrenIdList...)
}

func (c *RoomWSController) DefaultHandler(ctx *models.WSContext) {
	id, err := strconv.Atoi(ctx.ParsedMessage.Target)
	if err != nil {
		c.log(contracts.LError, err.Error())
		return
	}
	targetMember, err := c.roomRepo.GetMember(ctx.RoomId, uint64(id))
	if err != nil {
		c.log(contracts.LError, err.Error())
		return
	}

	if targetMember == nil {
		//ignoring
		return
	}
	var fullMessage map[string]interface{}
	err = json.Unmarshal(ctx.PureMessage, &fullMessage)
	if err != nil {
		c.log(contracts.LError, err.Error())
		return
	}
	userInfo, err := c.roomRepo.GetMember(ctx.RoomId, ctx.SocketID)
	if err != nil {
		c.log(contracts.LError, err.Error())
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

func (c *RoomWSController) log(level contracts.TLogLevel, msg ...string) {
	_ = c.logger.Log("room_ctrl", level, msg...)
}
