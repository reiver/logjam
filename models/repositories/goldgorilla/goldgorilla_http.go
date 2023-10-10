package GoldGorillaRepository

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sparkscience/logjam/models/contracts"
	"github.com/sparkscience/logjam/models/dto"
	"net/http"
	"time"
)

type HTTPRepository struct {
	client  *http.Client
	svcAddr string
}

func NewHTTPRepository() contracts.IGoldGorillaServiceRepository {
	return &HTTPRepository{
		client: &http.Client{
			Timeout: 8 * time.Second,
		},
		svcAddr: "",
	}
}

func (a *HTTPRepository) Init(svcAddr string) error {
	a.svcAddr = svcAddr
	return nil
}

func (a *HTTPRepository) isConfigured() bool {
	return len(a.svcAddr) > 0
}

func (a *HTTPRepository) CreatePeer(roomId string, id uint64, canPublish bool, isCaller bool) error {
	if !a.isConfigured() {
		return errors.New("gg repository not initialized yet")
	}
	body, err := getReader(
		dto.CreatePeerRPCModel{
			RoomPeerDTO: dto.RoomPeerDTO{
				RoomId: roomId,
				ID:     id,
			},
			CanPublish: canPublish,
			IsCaller:   isCaller,
		})
	if err != nil {
		return err
	}
	resp, err := a.client.Post(a.svcAddr+"/room/peer", "application/json", body)
	if err != nil {
		return err
	}
	if resp.StatusCode > 204 {
		return errors.New(resp.Status)
	}
	return nil
}

func (a *HTTPRepository) SendICECandidate(roomId string, id uint64, iceCandidate interface{}) error {
	if !a.isConfigured() {
		return errors.New("gg repository not initialized yet")
	}
	body, err := getReader(
		dto.SendIceCandidateReqModel{
			RoomPeerDTO: dto.RoomPeerDTO{
				RoomId: roomId,
				ID:     id,
			},
			ICECandidate: iceCandidate,
		})
	if err != nil {
		return err
	}
	resp, err := a.client.Post(a.svcAddr+"/room/ice", "application/json", body)
	if err != nil {
		return err
	}
	if resp.StatusCode > 204 {
		return errors.New(resp.Status)
	}
	return nil
}

func (a *HTTPRepository) SendAnswer(roomId string, peerId uint64, answer interface{}) error {
	if !a.isConfigured() {
		return errors.New("gg repository not initialized yet")
	}
	body, err := getReader(dto.SetSDPRPCModel{
		RoomPeerDTO: dto.RoomPeerDTO{
			RoomId: roomId,
			ID:     peerId,
		},
		SDP: answer,
	})
	if err != nil {
		return err
	}
	resp, err := a.client.Post(a.svcAddr+"/room/answer", "application/json", body)
	if err != nil {
		return err
	}
	if resp.StatusCode > 204 {
		return errors.New(resp.Status)
	}
	return nil
}

func (a *HTTPRepository) SendOffer(roomId string, peerId uint64, offer interface{}) error {
	if !a.isConfigured() {
		return errors.New("gg repository not initialized yet")
	}
	body, err := getReader(dto.SetSDPRPCModel{
		RoomPeerDTO: dto.RoomPeerDTO{
			RoomId: roomId,
			ID:     peerId,
		},
		SDP: offer,
	})
	if err != nil {
		return err
	}
	resp, err := a.client.Post(a.svcAddr+"/room/offer", "application/json", body)
	if err != nil {
		return err
	}
	if resp.StatusCode > 204 {
		return errors.New(resp.Status)
	}
	return nil
}

func (a *HTTPRepository) ClosePeer(roomId string, id uint64) error {
	if !a.isConfigured() {
		return errors.New("gg repository not initialized yet")
	}
	body, err := getReader(dto.RoomPeerDTO{
		RoomId: roomId,
		ID:     id,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodDelete, a.svcAddr+"/room/peer", body)
	if err != nil {
		return err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode > 204 {
		return errors.New(resp.Status)
	}
	return nil
}

func (a *HTTPRepository) ResetRoom(roomId string) error {
	if !a.isConfigured() {
		return errors.New("gg repository not initialized yet")
	}
	body, err := getReader(map[string]interface{}{"roomId": roomId})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodDelete, a.svcAddr+"/room/", body)
	if err != nil {
		return err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode > 204 {
		return errors.New(resp.Status)
	}
	return nil
}

func (a *HTTPRepository) Start() error {
	if a.svcAddr == "" {
		return errors.New("HTTPRepository instance is not initialized yet(waiting for goldgorilla hook)...")
	}
	resp, err := a.client.Post(a.svcAddr+"/room/", "application/json", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode > 204 {
		return errors.New(resp.Status)
	}
	return nil
}

func getReader(obj interface{}) (*bytes.Reader, error) {
	buffer, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(buffer)
	return body, nil
}
