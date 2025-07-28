package rooms

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/reiver/go-actsock"
	"github.com/reiver/go-fediverseid"

	"github.com/reiver/logjam/lib/goldgorilla"
	"github.com/reiver/logjam/lib/logjamlink"
	"github.com/reiver/logjam/lib/logs"
	libmetadata "github.com/reiver/logjam/lib/metadata"
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

	ctrl := &RoomWSController{
		logger:    logger.Tag(logtag),
		socketSVC: socketSVC,
		roomRepo:  roomRepo,
		ggRepo:    ggRepo,
	}

	go ctrl.tick()
	return ctrl
}

func (c *RoomWSController) tick() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		i := 1 + 1
		_ = i
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
		brLeft, err := c.roomRepo.GetBroadcasterLeftState(ctx.RoomId)
		if err != nil {
			c.error(err)
			//return
		}
		onBrLeft := func() {
			membersIdList, err := c.roomRepo.GetAllMembersId(ctx.RoomId, true)
			if err != nil {
				c.error(err)
				return
			}
			brDCEvent := msgs.Message{
				Type: msgs.TypeEventBroadcasterDisconnected,
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

		}

		peopleAreOnStage := false
		onStageMembers, err := c.roomRepo.GetOnStageMembersList(ctx.RoomId)
		if err == nil && onStageMembers != nil && len(onStageMembers) > 0 {
			peopleAreOnStage = true
		}
		if brLeft || !peopleAreOnStage {
			onBrLeft()
		} else {
			// so br is just reconnecting, lets give him some time
			err = c.roomRepo.StartBroadcasterReconnectionTimer(ctx.RoomId, onBrLeft)
			if err != nil {
				c.error(err)
				//return
			}
		}

	} else {
		parentDCEvent := msgs.Message{
			Type: msgs.TypeEventParentDC,
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
		err = c.roomRepo.SetRoomMetaData(ctx.RoomId, libmetadata.MetaData{})
		if err != nil {
			c.error(err)
		}
	} else if len(membersIdList) == 1 {
		if c.roomRepo.IsGGInstance(ctx.RoomId, membersIdList[0]) {
			err = c.roomRepo.ClearMessageHistory(ctx.RoomId)
			if err != nil {
				c.error(err)
			}
			err = c.roomRepo.SetRoomMetaData(ctx.RoomId, libmetadata.MetaData{})
			if err != nil {
				c.error(err)
			}
		}
	}
	c.roomRepo.DelMemberFromStage(ctx.RoomId, ctx.SocketID)
}

func (c *RoomWSController) Leave(ctx *WSContext) {
	isBroadcaster, err := c.roomRepo.IsBroadcaster(ctx.RoomId, ctx.SocketID)
	if err != nil {
		c.error(err)
		return
	}
	if isBroadcaster {
		memberIds, err := c.roomRepo.GetAllMembersId(ctx.RoomId, true)
		if err != nil {
			c.error(err)
			return
		}
		_ = c.roomRepo.SetBroadcasterLeftState(ctx.RoomId, true)
		err = c.socketSVC.Send(msgs.Message{
			Type: msgs.TypeBroadcasterLeft,
		}, memberIds...)
		go c.socketSVC.Disconnect(ctx.SocketID)
		if err != nil {
			c.error(err)
			return
		}
	}
}

func (c *RoomWSController) Start(ctx *WSContext) {
	resultEvent := msgs.Message{
		Type:   msgs.TypeStart,
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

	{
		var acctURI string
		var id string
		{
			// The room-id is the fediverse-id of the actor who created the room.
			var fediverseIDString string = ctx.RoomId

			fediverseID, err := fediverseid.ParseFediverseIDString(fediverseIDString)
			if nil == err {
				acctURI = fediverseID.AcctURI()
				id = logjamlink.LogJamLink(fediverseID.HostElse(""), fediverseID.NameElse(""))
			}
		}

		var activity = actsock.Create{
			Actor: acctURI,
			Object: actsock.Conference{
				Actor: acctURI,
				ID:    id,
				Origin: []string{
					acctURI,
				},
				To: []string{
					acctURI,
				},
			},
		}
		_ = c.socketSVC.Send(activity, ctx.SocketID)
		c.debugf("[actsock] %T\n%s", activity, activity)
	}
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
	ggid, err := c.roomRepo.GetRoomGoldGorillaId(ctx.RoomId)
	if err != nil {
		c.error(err)
		return
	}
	peopleAreOnStage := false
	onStageMembers, err := c.roomRepo.GetOnStageMembersList(ctx.RoomId)
	if err == nil && onStageMembers != nil && len(onStageMembers) > 0 {
		peopleAreOnStage = true
	}

	resultEvent := msgs.Message{
		Type:   msgs.TypeRole,
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
				_ = c.socketSVC.Send(msgs.Message{
					Type: msgs.TypeRole,
					Data: "no:broadcast",
				}, ctx.SocketID)
				return
			}
			if currentUser.Name != br.Name {
				_ = c.socketSVC.Send(msgs.Message{
					Type: msgs.TypeRole,
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
		ggEnabled := true
		ggEnabledInReqBody, exists := eventData["ggEnabled"]
		if exists && ggEnabledInReqBody == false {
			ggEnabled = false
		}

		if ggid == nil && ggEnabled {
			err := c.ggRepo.Start(ctx.RoomId)
			if err != nil {
				c.error(err)
				ggEnabled = false
				//return
			}
		}
		if ggEnabled {
			brParentId, err := c.roomRepo.InsertMemberToTree(ctx.RoomId, ctx.SocketID, false, true)
			if brParentId != nil {
				println(fmt.Sprintf(`brParentId: %d`, *brParentId))
			}
			if err != nil {
				c.error(err)
				return
			}
			ggid, err = c.roomRepo.GetRoomGoldGorillaId(ctx.RoomId)
			if err != nil {
				c.error(err)
				return
			}
			err = c.ggRepo.CreatePeer(ctx.RoomId, ctx.SocketID, true, true, *ggid, goldgorilla.CDSend)
			if err != nil {
				c.error(err)
				return
			}
			err = c.ggRepo.CreatePeer(ctx.RoomId, ctx.SocketID, true, true, *ggid, goldgorilla.CDRecv)
			if err != nil {
				c.error(err)
				return
			}

			type customMsg struct {
				msgs.Message
				IsGG bool `json:"isGoldGorilla"`
			}
			_ = c.socketSVC.Send(customMsg{
				Message: msgs.Message{
					Type: msgs.TypeAddAudience,
					Data: strconv.FormatUint(*ggid, 10),
				},
				IsGG: true,
			}, ctx.SocketID)
		} else {
			err = c.roomRepo.SetBroadcaster(ctx.RoomId, ctx.SocketID)
			if err != nil {
				c.error(err)
				return
			}
		}

		if !peopleAreOnStage {
			memberIds, err := c.roomRepo.GetAllMembersId(ctx.RoomId, true)
			if err != nil {
				c.error(err)
				return
			}

			err = c.socketSVC.Send(msgs.Message{
				Type: msgs.TypeBroadcasting,
				Data: strconv.FormatUint(ctx.SocketID, 10),
			}, memberIds...)
			if err != nil {
				c.error(err)
			}
		} else {
			// then this is the br reconnection
			err = c.roomRepo.OnBroadcasterConnectedBack(ctx.RoomId)
			if err != nil {
				c.error(err)
			}
		}
		_ = c.socketSVC.Send(resultEvent, ctx.SocketID)
	} else if ctx.ParsedMessage.Data == msgs.TypeAltBroadcast {
		broadcaster, err := c.roomRepo.GetBroadcaster(ctx.RoomId)
		if err != nil {
			c.error(err)
			return
		}
		if broadcaster == nil {
			_ = c.socketSVC.Send(msgs.Message{
				Type: msgs.TypeAltBroadcast,
				Data: "no-broadcaster",
			}, ctx.SocketID)
			return
		}

		err = c.ggRepo.CreatePeer(ctx.RoomId, ctx.SocketID, true, true, *ggid, goldgorilla.CDSend)
		if err != nil {
			c.error(err)
			return
		}

		resultEvent.Type = msgs.TypeAltBroadcast
		resultEvent.Data = strconv.FormatUint(broadcaster.ID, 10)
		_ = c.socketSVC.Send(resultEvent, ctx.SocketID)

		userInfo, err := c.roomRepo.GetMember(ctx.RoomId, ctx.SocketID)
		if err != nil {
			c.error(err)
			return
		}
		broadcasterReceivingEvent := map[string]string{
			"Type": msgs.TypeAltBroadcast,
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
			if !peopleAreOnStage {
				_ = c.socketSVC.Send(msgs.Message{
					Type: msgs.TypeRole,
					Data: "no:audience",
				}, ctx.SocketID)
				return
			}
		}
		tryCount := 0
	start:
		parentId, err := c.roomRepo.InsertMemberToTree(ctx.RoomId, ctx.SocketID, false, false)
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
			println("parent id", *parentId)
			_ = c.socketSVC.Send(msgs.Message{
				Type: msgs.TypeAddAudience,
				Data: strconv.FormatUint(ctx.SocketID, 10),
			}, *parentId)
		} else {
			if ggid == nil {
				c.socketSVC.Disconnect(ctx.SocketID) // so aud will reconnect
				go c.roomRepo.RemoveMember(ctx.RoomId, *parentId)
				return
			}
			err = c.ggRepo.CreatePeer(ctx.RoomId, ctx.SocketID, true, false, *ggid, goldgorilla.CDRecv)
			if err != nil {
				_ = c.socketSVC.Send(msgs.Message{
					Type: msgs.TypeError,
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
	_ = c.socketSVC.Send(msgs.Message{Type: msgs.TypePong, Data: ctx.ParsedMessage.Data}, ctx.SocketID)
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
	buffer, _ := json.Marshal(tree)
	resultEvent := msgs.Message{
		Type: msgs.TypeTree,
		Data: string(buffer),
	}
	_ = c.socketSVC.Send(resultEvent, ctx.SocketID)
}

func (c *RoomWSController) MetaDataSet(ctx *WSContext) {
	var metaData libmetadata.MetaData
	err := json.Unmarshal([]byte(ctx.ParsedMessage.Data), &metaData)
	if nil != err {
		c.error(err)
		return
	}
	err = c.roomRepo.SetRoomMetaData(ctx.RoomId, metaData)
	if err != nil {
		c.error(err)
		return
	}
}

func (c *RoomWSController) MetaDataGet(ctx *WSContext) {
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
	resultEvent := msgs.Message{
		Type: msgs.TypeMetaDataGet,
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
	resultEvent := msgs.Message{
		Type: msgs.TypeUserByStream,
	}
	userRole := "audience"
	if isBroadcaster {
		userRole = "broadcast"
	}
	resultEvent.Data = strconv.FormatUint(userInfo.ID, 10) + "," + userInfo.Name + "," + ctx.ParsedMessage.Data + "," + userRole
	_ = c.socketSVC.Send(resultEvent, ctx.SocketID)
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
	_ = c.socketSVC.Send(ctx.PureMessage, list...)
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
		if c.roomRepo.IsGGInstance(roomId, v.ID) {
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
	event := msgs.Message{
		Type: msgs.TypeUserEvent,
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
	event := msgs.Message{
		Type: msgs.TypeReconnect,
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
	_ = c.socketSVC.Send(msgs.Message{
		Type:   msgs.TypeNewMessage,
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
	err = c.ggRepo.SendOffer(ctx.RoomId, ctx.SocketID, msg["sdp"], goldgorilla.CDSend)
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
	err = c.ggRepo.SendAnswer(ctx.RoomId, ctx.SocketID, msg["sdp"], goldgorilla.CDRecv)
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
	var dir goldgorilla.ConnectionDirection
	dir = goldgorilla.CDRecv
	if cdir, exists := msg["connectionDirection"]; exists {
		if cdir == goldgorilla.CDSend {
			dir = goldgorilla.CDSend
		}
	}
	err = c.ggRepo.SendICECandidate(ctx.RoomId, ctx.SocketID, msg["candidate"], dir)
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
	roomGGID, err := c.roomRepo.GetRoomGoldGorillaId(ctx.RoomId)
	if err != nil {
		c.error(err)
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
	var payload map[string]any
	err = json.Unmarshal(ctx.PureMessage, &payload)
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
	payload["username"] = userInfo.Name
	payload["data"] = strconv.FormatUint(ctx.SocketID, 10)
	if slices.Index([]string{
		"invite-to-stage",
		"alt-broadcast-approve",
		"-",
	}, ctx.ParsedMessage.Type) >= 0 {
		if roomGGID != nil {
			payload["goldgorillaID"] = strconv.FormatUint(*roomGGID, 10)
		}
	}
	if ctx.ParsedMessage.Type == "audience-broadcasting" {
		if joinedStage, exists := payload["joinedStage"]; exists {
			if joinedStageState, isBoolean := joinedStage.(bool); isBoolean {
				if joinedStageState {
					err = c.roomRepo.AddOnStageMember(ctx.RoomId, ctx.SocketID)
					if err != nil {
						c.error(err)
					}
				} else {
					err = c.roomRepo.DelMemberFromStage(ctx.RoomId, ctx.SocketID)
					if err != nil {
						c.error(err)
					}
				}

			}
		}
	}
	_ = c.socketSVC.Send(payload, targetMember.ID)
}

func (c *RoomWSController) debug(msg ...any) {
	c.logger.Debug(msg...)
}

func (c *RoomWSController) debugf(format string, msg ...any) {
	c.logger.Debugf(format, msg...)
}

func (c *RoomWSController) error(msg ...any) {
	c.logger.Error(msg...)
}

func (c *RoomWSController) info(msg ...any) {
	c.logger.Info(msg...)
}
