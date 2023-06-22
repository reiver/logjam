package controllers

import (
	"encoding/json"
	"github.com/sparkscience/logjam/models"
	"github.com/sparkscience/logjam/models/contracts"
	"strconv"
)

type RoomController struct {
	logger    contracts.ILogger
	socketSVC contracts.ISocketService
	roomRepo  contracts.IRoomRepository
}

func NewRoomController(socketSVC contracts.ISocketService, roomRepo contracts.IRoomRepository, logger contracts.ILogger) *RoomController {
	return &RoomController{
		logger:    logger,
		socketSVC: socketSVC,
		roomRepo:  roomRepo,
	}
}

func (c *RoomController) OnConnect(ctx *models.WSContext) {

}

func (c *RoomController) OnDisconnect(ctx *models.WSContext) {

}

func (c *RoomController) Start(ctx *models.WSContext) {
	resultEvent := models.MessageContract{
		Type:   "start",
		Data:   strconv.FormatInt(int64(ctx.SocketID), 10),
		Target: "",
		Name:   "",
	}
	_ = c.roomRepo.CreateRoom(ctx.RoomId)
	err := c.socketSVC.Send(resultEvent, ctx.SocketID)
	println("sent")
	if err != nil {
		println(err.Error())
	}
}
func (c *RoomController) Role(ctx *models.WSContext) {
	var eventData map[string]string
	err := json.Unmarshal(ctx.PureMessage, &eventData)
	if err != nil {
		return
	}
	streamId, exists := eventData["streamId"]
	err = c.roomRepo.AddMember(ctx.RoomId, ctx.SocketID, "", "", streamId)
	if err != nil {
		println(err.Error())
		return
	}
	if exists {
		err = c.roomRepo.UpdateMemberMeta(ctx.RoomId, ctx.SocketID, "streamId", streamId)
		if err != nil {
			println(err.Error())
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
			println(err.Error())
			return
		}
		err = c.roomRepo.SetBroadcaster(ctx.RoomId, ctx.SocketID, "", streamId)
		if err != nil {
			println(err.Error())
			return
		}
		memberIds, err := c.roomRepo.GetAllMembersId(ctx.RoomId, true)
		if err != nil {
			println(err.Error())
			return
		}
		err = c.socketSVC.Send(models.MessageContract{
			Type: "broadcasting",
			Data: strconv.FormatUint(ctx.SocketID, 10),
		}, memberIds...)
		if err != nil {
			println(err.Error())
		}
		_ = c.socketSVC.Send(resultEvent, ctx.SocketID)
	} else if ctx.ParsedMessage.Data == "alt-broadcast" {
		broadcaster, err := c.roomRepo.GetBroadcaster(ctx.RoomId)
		if err != nil {
			println(err.Error())
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

		//todo: set name
		broadcasterReceivingEvent := models.MessageContract{
			Type: `alt-broadcast`,
			Data: strconv.FormatUint(ctx.SocketID, 10),
			Name: "",
		}
		_ = c.socketSVC.Send(broadcasterReceivingEvent, broadcaster.ID)
	} else if ctx.ParsedMessage.Data == "audience" {
		broadcaster, err := c.roomRepo.GetBroadcaster(ctx.RoomId)
		if err != nil {
			println(err.Error())
			return
		}
		if broadcaster == nil {
			_ = c.socketSVC.Send(models.MessageContract{
				Type: "role",
				Data: "no:audience",
			}, ctx.SocketID)
			return
		}
		parentId, err := c.roomRepo.InsertMemberToTree(ctx.RoomId, ctx.SocketID)
		if err != nil {
			println(err.Error())
			return
		}
		c.socketSVC.Send(models.MessageContract{
			Type: "add_audience",
			Data: strconv.FormatUint(ctx.SocketID, 10),
		}, *parentId)
	}
}
func (c *RoomController) Stream(ctx *models.WSContext) {
	payload := make(map[string]string)
	err := json.Unmarshal(ctx.PureMessage, &payload)
	if err != nil {
		println(err.Error())
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
		println(err.Error())
		return
	}
}

func (c *RoomController) Ping(ctx *models.WSContext) {
	_ = c.socketSVC.Send(models.MessageContract{Type: "pong"}, ctx.SocketID)
}

func (c *RoomController) TurnStatus(ctx *models.WSContext) {

}

func (c *RoomController) Tree(ctx *models.WSContext) {

}

func (c *RoomController) MetadataSet(ctx *models.WSContext) {

}

func (c *RoomController) MetadataGet(ctx *models.WSContext) {

}

func (c *RoomController) UserByStream(ctx *models.WSContext) {

}

func (c *RoomController) GetLatestUserList(ctx *models.WSContext) {

}

func (c *RoomController) DefaultHandler(ctx *models.WSContext) {
	id, err := strconv.Atoi(ctx.ParsedMessage.Target)
	if err != nil {
		println(err.Error())
		return
	}
	targetMember, err := c.roomRepo.GetMember(ctx.RoomId, uint64(id))
	if err != nil {
		println(err.Error())
		return
	}

	var fullMessage map[string]interface{}
	err = json.Unmarshal(ctx.PureMessage, &fullMessage)
	if err != nil {
		println(err.Error())
		return
	}
	//todo: get name
	fullMessage["username"] = ""
	fullMessage["data"] = strconv.FormatUint(ctx.SocketID, 10)
	c.socketSVC.Send(fullMessage, targetMember.ID)
}
