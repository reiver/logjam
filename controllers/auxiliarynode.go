package controllers

import (
	"encoding/json"
	"github.com/sparkscience/logjam/models"
	"github.com/sparkscience/logjam/models/contracts"
	"github.com/sparkscience/logjam/models/dto"
	"io"
	"net/http"
	"strconv"
)

type AuxiliaryNodeController struct {
	roomRepo  contracts.IRoomRepository
	socketSVC contracts.ISocketService
	helper    *RestResponseHelper
	logger    contracts.ILogger
}

func NewAuxiliaryNodeController(roomRepo contracts.IRoomRepository, socketSVC contracts.ISocketService, logger contracts.ILogger) *AuxiliaryNodeController {
	return &AuxiliaryNodeController{
		roomRepo:  roomRepo,
		socketSVC: socketSVC,
		logger:    logger,
	}
}

func (ctrl *AuxiliaryNodeController) SendICECandidate(rw http.ResponseWriter, req *http.Request) {

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
	err = ctrl.roomRepo.AddMember(reqModel.RoomId, models.AuxiliaryNodeId, "auxiliary-node", "", "")
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
}
