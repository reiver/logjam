package auxiliarynode

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sparkscience/logjam/models/contracts"
	"github.com/sparkscience/logjam/models/dto"
	"net/http"
	"time"
)

type auxiliaryNodeRepository struct {
	client  *http.Client
	svcAddr string
}

func NewAuxiliaryNodeRepository() contracts.IAuxiliaryNodeServiceRepository {
	return &auxiliaryNodeRepository{
		client: &http.Client{
			Timeout: 8 * time.Second,
		},
		svcAddr: "",
	}
}

func (a *auxiliaryNodeRepository) Init(svcAddr string) error {
	a.svcAddr = svcAddr
	return nil
}

func (a *auxiliaryNodeRepository) CreatePeer(roomId string, id uint64, canPublish bool, isCaller bool) error {
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

func (a *auxiliaryNodeRepository) SendICECandidate(roomId string, id uint64, iceCandidate any) error {
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

func (a *auxiliaryNodeRepository) SendAnswer(roomId string, peerId uint64, answer any) error {
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

func (a *auxiliaryNodeRepository) SendOffer(roomId string, peerId uint64, offer any) error {
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

func (a *auxiliaryNodeRepository) ClosePeer(roomId string, id uint64) error {
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

func (a *auxiliaryNodeRepository) ResetRoom(roomId string) error {
	body, err := getReader(map[string]any{"roomId": roomId})
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

func (a *auxiliaryNodeRepository) Start() error {
	if a.svcAddr == "" {
		return errors.New("auxiliaryNodeRepository instance is not initialized yet(waiting for goldgorilla hook)...")
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

func getReader(obj any) (*bytes.Reader, error) {
	buffer, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(buffer)
	return body, nil
}
