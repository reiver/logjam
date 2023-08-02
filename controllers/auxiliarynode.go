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

type AuxiliaryNodeController struct {
	roomRepo  contracts.IRoomRepository
	anSVCRepo contracts.IAuxiliaryNodeServiceRepository
	socketSVC contracts.ISocketService
	conf      *models.ConfigModel
	helper    *RestResponseHelper
	logger    contracts.ILogger
}

func NewAuxiliaryNodeController(roomRepo contracts.IRoomRepository, anSVCRepo contracts.IAuxiliaryNodeServiceRepository, socketSVC contracts.ISocketService, conf *models.ConfigModel, helper *RestResponseHelper, logger contracts.ILogger) *AuxiliaryNodeController {
	return &AuxiliaryNodeController{
		roomRepo:  roomRepo,
		socketSVC: socketSVC,
		logger:    logger,
		anSVCRepo: anSVCRepo,
		conf:      conf,
		helper:    helper,
	}
}

func (ctrl *AuxiliaryNodeController) SendAnswer(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	var reqModel dto.SetSDPRPCModel
	err = json.Unmarshal(reqBody, &reqModel)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	ctrl.socketSVC.Send(map[string]any{
		"type":   "video-answer",
		"target": strconv.FormatUint(reqModel.ID, 10),
		"name":   strconv.FormatUint(models.AuxiliaryNodeId, 10),
		"sdp":    reqModel.SDP,
		"data":   strconv.FormatUint(models.AuxiliaryNodeId, 10),
	}, reqModel.ID)
	_ = ctrl.helper.Write(rw, nil, 204)
}

func (ctrl *AuxiliaryNodeController) SendOffer(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	var reqModel dto.SetSDPRPCModel
	err = json.Unmarshal(reqBody, &reqModel)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	_ = ctrl.socketSVC.Send(map[string]any{
		"type":   "video-offer",
		"target": strconv.FormatUint(reqModel.ID, 10),
		"name":   strconv.FormatUint(models.AuxiliaryNodeId, 10),
		"sdp":    reqModel.SDP,
		"data":   strconv.FormatUint(models.AuxiliaryNodeId, 10),
	}, reqModel.ID)
	_ = ctrl.helper.Write(rw, nil, 204)
}

func (ctrl *AuxiliaryNodeController) SendICECandidate(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	var reqModel dto.SendIceCandidateReqModel
	err = json.Unmarshal(reqBody, &reqModel)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	_ = ctrl.socketSVC.Send(map[string]any{
		"Type":      "new-ice-candidate",
		"Target":    strconv.FormatUint(reqModel.ID, 10),
		"candidate": reqModel.ICECandidate,
		"data":      strconv.FormatUint(models.AuxiliaryNodeId, 10),
	}, reqModel.ID)
	_ = ctrl.helper.Write(rw, nil, 204)
}

func (ctrl *AuxiliaryNodeController) Join(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	var reqModel dto.JoinReqModel
	err = json.Unmarshal(reqBody, &reqModel)
	if ctrl.helper.HandleIfErr(rw, err, 400) {
		return
	}
	ctrl.conf.AuxiliaryNodeSVCAddr = reqModel.ServiceAddr
	_ = ctrl.anSVCRepo.Init(reqModel.ServiceAddr)
	err = ctrl.roomRepo.AddMember(reqModel.RoomId, models.AuxiliaryNodeId, "{}", "", "")
	if ctrl.helper.HandleIfErr(rw, err, 500) {
		return
	}
	err = ctrl.roomRepo.UpdateCanConnect(reqModel.RoomId, models.AuxiliaryNodeId, true)
	if ctrl.helper.HandleIfErr(rw, err, 500) {
		return
	}
	parentId, err := ctrl.roomRepo.InsertMemberToTree(reqModel.RoomId, models.AuxiliaryNodeId, true)
	if ctrl.helper.HandleIfErr(rw, err, 500) {
		return
	}
	_ = ctrl.socketSVC.Send(models.MessageContract{
		Type: "add_audience",
		Data: strconv.FormatUint(models.AuxiliaryNodeId, 10),
	}, *parentId)
	err = ctrl.anSVCRepo.CreatePeer(reqModel.RoomId, *parentId, true, true)
	if ctrl.helper.HandleIfErr(rw, err, 503) {
		return
	}
	_ = ctrl.helper.Write(rw, nil, 204)
	go func(roomId string, svcAddr string) {
		for {
			res, err := http.Get(svcAddr + "/healthcheck")
			if err != nil {
				break
			}
			if res.StatusCode > 204 {
				break
			}
			time.Sleep(2 * time.Second)
		}
		_, childrenIdList, err := ctrl.roomRepo.RemoveMember(roomId, models.AuxiliaryNodeId)
		if err != nil {
			println(err.Error())
			return
		}
		parentDCEvent := models.MessageContract{
			Type: "event-parent-dc",
			Data: strconv.FormatUint(models.AuxiliaryNodeId, 10),
		}
		_ = ctrl.socketSVC.Send(parentDCEvent, childrenIdList...)
		println("deleted an from tree")
	}(reqModel.RoomId, reqModel.ServiceAddr)
}
