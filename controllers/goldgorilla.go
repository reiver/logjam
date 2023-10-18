package controllers

import (
	"encoding/json"
	"github.com/sparkscience/logjam/models"
	"github.com/sparkscience/logjam/models/contracts"
	"github.com/sparkscience/logjam/models/dto"
	"io"
	"net/http"
	"strconv"
	"time"
)

type GoldGorillaController struct {
	roomRepo  contracts.IRoomRepository
	ggSVCRepo contracts.IGoldGorillaServiceRepository
	socketSVC contracts.ISocketService
	conf      *models.ConfigModel
	helper    *RestResponseHelper
	logger    contracts.ILogger
}

func NewGoldGorillaController(roomRepo contracts.IRoomRepository, ggSVCRepo contracts.IGoldGorillaServiceRepository, socketSVC contracts.ISocketService, conf *models.ConfigModel, helper *RestResponseHelper, logger contracts.ILogger) *GoldGorillaController {
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
		"name":   strconv.FormatUint(models.GetGoldGorillaId(), 10),
		"sdp":    reqModel.SDP,
		"data":   strconv.FormatUint(models.GetGoldGorillaId(), 10),
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
		"name":   strconv.FormatUint(models.GetGoldGorillaId(), 10),
		"sdp":    reqModel.SDP,
		"data":   strconv.FormatUint(models.GetGoldGorillaId(), 10),
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
		"data":      strconv.FormatUint(models.GetGoldGorillaId(), 10),
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
	ctrl.conf.GoldGorillaSVCAddr = reqModel.ServiceAddr
	_ = ctrl.ggSVCRepo.Init(reqModel.ServiceAddr)
	models.DecreaseGoldGorillaId()
	err = ctrl.roomRepo.AddMember(reqModel.RoomId, models.GetGoldGorillaId(), "{}", "", "")
	if ctrl.helper.HandleIfErr(rw, err, 500) {
		return
	}
	err = ctrl.roomRepo.UpdateCanConnect(reqModel.RoomId, models.GetGoldGorillaId(), true)
	if ctrl.helper.HandleIfErr(rw, err, 500) {
		return
	}
	parentId, err := ctrl.roomRepo.InsertMemberToTree(reqModel.RoomId, models.GetGoldGorillaId(), true)
	if ctrl.helper.HandleIfErr(rw, err, 500) {
		_, _, _ = ctrl.roomRepo.RemoveMember(reqModel.RoomId, models.GetGoldGorillaId())
		return
	}
	err = ctrl.ggSVCRepo.CreatePeer(reqModel.RoomId, *parentId, true, true)
	if ctrl.helper.HandleIfErr(rw, err, 503) {
		return
	}
	_ = ctrl.socketSVC.Send(models.MessageContract{
		Type: "add_audience",
		Data: strconv.FormatUint(models.GetGoldGorillaId(), 10),
	}, *parentId)

	_ = ctrl.helper.Write(rw, nil, 204)
	go func(roomId string, svcAddr string) {
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
		_, childrenIdList, err := ctrl.roomRepo.RemoveMember(roomId, models.GetGoldGorillaId())
		if err != nil {
			println(err.Error())
			return
		}
		parentDCEvent := models.MessageContract{
			Type: "event-parent-dc",
			Data: strconv.FormatUint(models.GetGoldGorillaId(), 10),
		}
		_ = ctrl.socketSVC.Send(parentDCEvent, childrenIdList...)
		println("deleted goldgorilla from tree")
	}(reqModel.RoomId, reqModel.ServiceAddr)
}

func (ctrl *GoldGorillaController) RejoinGoldGorilla(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	var reqModel struct {
		RoomId string `json:"roomId"`
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
	_, _, err = ctrl.roomRepo.RemoveMember(reqModel.RoomId, models.GetGoldGorillaId())
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

	go func(membersIdList []uint64) {
		time.Sleep(500 * time.Millisecond)
		err := ctrl.ggSVCRepo.Start()
		if err != nil {
			println(err.Error())
			return
		}
		brIsBackEvent := models.MessageContract{
			Type: "broadcasting",
		}

		_ = ctrl.socketSVC.Send(brIsBackEvent, membersIdList...)
	}(roomMembersIdList)
	_ = ctrl.helper.Write(rw, nil, 204)
}
