package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/lib/logs"
	"github.com/reiver/logjam/models"
	"github.com/reiver/logjam/models/contracts"
	"github.com/reiver/logjam/models/dto"
)

const (
	logtag = "goldgorilla"
)

type GoldGorillaController struct {
	roomRepo  contracts.IRoomRepository
	ggSVCRepo contracts.IGoldGorillaServiceRepository
	socketSVC contracts.ISocketService
	conf      cfg.Configurer
	helper    *RestResponseHelper
	logger    logs.TaggedLogger
}

func NewGoldGorillaController(roomRepo contracts.IRoomRepository, ggSVCRepo contracts.IGoldGorillaServiceRepository, socketSVC contracts.ISocketService, conf cfg.Configurer, helper *RestResponseHelper, logger logs.TaggedLogger) *GoldGorillaController {
	return &GoldGorillaController{
		roomRepo:  roomRepo,
		socketSVC: socketSVC,
		logger:    logger,
		ggSVCRepo: ggSVCRepo,
		conf:      conf,
		helper:    helper,
	}
}

func (ctrl *GoldGorillaController) SendAnswer(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	var reqModel dto.SetSDPRPCModel
	err = json.Unmarshal(reqBody, &reqModel)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	ctrl.socketSVC.Send(map[string]interface{}{
		"type":   "video-answer",
		"target": strconv.FormatUint(reqModel.ID, 10),
		"name":   strconv.FormatUint(reqModel.GGID, 10),
		"sdp":    reqModel.SDP,
		"data":   strconv.FormatUint(reqModel.GGID, 10),
	}, reqModel.ID)
	_ = ctrl.helper.Write(rw, nil, 204)
}

func (ctrl *GoldGorillaController) SendOffer(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	var reqModel dto.SetSDPRPCModel
	err = json.Unmarshal(reqBody, &reqModel)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	_ = ctrl.socketSVC.Send(map[string]interface{}{
		"type":   "video-offer",
		"target": strconv.FormatUint(reqModel.ID, 10),
		"name":   strconv.FormatUint(reqModel.GGID, 10),
		"sdp":    reqModel.SDP,
		"data":   strconv.FormatUint(reqModel.GGID, 10),
	}, reqModel.ID)
	_ = ctrl.helper.Write(rw, nil, 204)
}

func (ctrl *GoldGorillaController) SendICECandidate(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	var reqModel dto.SendIceCandidateReqModel
	err = json.Unmarshal(reqBody, &reqModel)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	_ = ctrl.socketSVC.Send(map[string]interface{}{
		"Type":      "new-ice-candidate",
		"Target":    strconv.FormatUint(reqModel.ID, 10),
		"candidate": reqModel.ICECandidate,
		"data":      strconv.FormatUint(reqModel.GGID, 10),
	}, reqModel.ID)
	_ = ctrl.helper.Write(rw, nil, 204)
}

func (ctrl *GoldGorillaController) Join(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	var reqModel dto.JoinReqModel
	err = json.Unmarshal(reqBody, &reqModel)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	//ctrl.conf.GoldGorillaSVCAddr = reqModel.ServiceAddr
	newGGID := ctrl.socketSVC.GetNewID()
	err = ctrl.roomRepo.AddMember(reqModel.RoomId, newGGID, "{}", "", "", true)
	if ctrl.helper.HandleIfErr(rw, err, 500) {
		return
	}
	err = ctrl.roomRepo.UpdateCanConnect(reqModel.RoomId, newGGID, true)
	if ctrl.helper.HandleIfErr(rw, err, 500) {
		return
	}
	parentId, err := ctrl.roomRepo.InsertMemberToTree(reqModel.RoomId, newGGID, true)
	if ctrl.helper.HandleIfErr(rw, err, 500) {
		_, _, _ = ctrl.roomRepo.RemoveMember(reqModel.RoomId, newGGID)
		return
	}
	err = ctrl.ggSVCRepo.CreatePeer(reqModel.RoomId, *parentId, true, true, newGGID)
	if ctrl.helper.HandleIfErr(rw, err, 503) {
		return
	}
	_ = ctrl.socketSVC.Send(models.MessageContract{
		Type: "add_audience",
		Data: strconv.FormatUint(newGGID, 10),
	}, *parentId)

	_ = ctrl.helper.Write(rw, struct {
		ID uint64 `json:"id"`
	}{
		ID: newGGID,
	}, 200)
	memsId, err := ctrl.roomRepo.GetAllMembersId(reqModel.RoomId, false)
	if err != nil {
		ctrl.logger.Error(logtag, err)
	} else {
		_ = ctrl.socketSVC.Send(models.MessageContract{Type: "goldgorilla-joined", Data: strconv.FormatUint(newGGID, 10)}, memsId...)
	}
	go func(roomId string, svcAddr string, ggId uint64) {
		for {
			res, err := http.Get(svcAddr + "/healthcheck?roomId=" + roomId)
			if err != nil {
				break
			}
			if res.StatusCode > 204 {
				break
			}
			time.Sleep(2 * time.Second)
		}
		_, childrenIdList, err := ctrl.roomRepo.RemoveMember(roomId, newGGID)
		if err != nil {
			ctrl.logger.Error(logtag, err)
			return
		}
		parentDCEvent := models.MessageContract{
			Type: "event-parent-dc",
			Data: strconv.FormatUint(newGGID, 10),
		}
		_ = ctrl.socketSVC.Send(parentDCEvent, childrenIdList...)
		ctrl.logger.Info(logtag, "deleted a goldgorilla instance from tree")
	}(reqModel.RoomId, ctrl.conf.GoldGorillaBaseURL(), newGGID)
}

func (ctrl *GoldGorillaController) RejoinGoldGorilla(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	var reqModel struct {
		RoomId string `json:"roomId"`
		GGID   uint64 `json:"ggid"`
	}
	err = json.Unmarshal(reqBody, &reqModel)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	broadcaster, err := ctrl.roomRepo.GetBroadcaster(reqModel.RoomId)
	if ctrl.helper.HandleIfErr(rw, err, 500) {
		return
	}

	if broadcaster == nil {
		_ = ctrl.helper.Write(rw, nil, 503)
		return
	}
	_, _, err = ctrl.roomRepo.RemoveMember(reqModel.RoomId, reqModel.GGID)
	if ctrl.helper.HandleIfErr(rw, err, 500) {
		return
	}
	roomMembersIdList, err := ctrl.roomRepo.GetAllMembersId(reqModel.RoomId, true)
	if ctrl.helper.HandleIfErr(rw, err, 500) {
		return
	}

	brDCEvent := models.MessageContract{
		Type: "event-broadcaster-disconnected",
		Data: strconv.FormatUint(broadcaster.ID, 10),
	}

	_ = ctrl.socketSVC.Send(brDCEvent, roomMembersIdList...)

	go func(roomId string, membersIdList []uint64) {
		time.Sleep(500 * time.Millisecond)
		err := ctrl.ggSVCRepo.Start(roomId)
		if err != nil {
			ctrl.logger.Error(logtag, err)
			return
		}
		brIsBackEvent := models.MessageContract{
			Type: "broadcasting",
		}

		_ = ctrl.socketSVC.Send(brIsBackEvent, membersIdList...)
	}(reqModel.RoomId, roomMembersIdList)
	_ = ctrl.helper.Write(rw, nil, 204)
}
