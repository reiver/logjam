package GoldGorillaRepository

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"github.com/reiver/logjam/lib/goldgorilla"
	"github.com/reiver/logjam/models/dto"
	"time"
)

type HTTPRepository struct {
	client  *http.Client
	svcAddr string
}

func NewHTTPRepository(svcAddr string) goldgorilla.IGoldGorillaServiceRepository {
	return &HTTPRepository{
		client: &http.Client{
			Timeout: 8 * time.Second,
		},
		svcAddr: svcAddr,
	}
}

func (a *HTTPRepository) Init(svcAddr string) error {
	a.svcAddr = svcAddr
	return nil
}

func (a *HTTPRepository) isConfigured() bool {
	return len(a.svcAddr) > 0
}

func (a *HTTPRepository) CreatePeer(roomId string, id uint64, canPublish bool, isCaller bool, ggid uint64) error {
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
			GGID:       ggid,
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

func (a *HTTPRepository) ResetRoom(roomId string) (*uint64, error) {
	if !a.isConfigured() {
		return nil, errors.New("gg repository not initialized yet")
	}
	body, err := getReader(map[string]interface{}{"roomId": roomId})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodDelete, a.svcAddr+"/room/", body)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 204 {
		return nil, errors.New(resp.Status)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
	}
	defer resp.Body.Close()
	respModel := struct {
		GGID uint64 `json:"ggid"`
	}{}
	err = json.Unmarshal(respBody, &respModel)
	return &respModel.GGID, err
}

func (a *HTTPRepository) Start(roomId string) error {
	if a.svcAddr == "" {
		return errors.New("HTTPRepository instance is not initialized yet(waiting for goldgorilla hook)...")
	}

	body, _ := getReader(map[string]string{"roomId": roomId})
	resp, err := a.client.Post(a.svcAddr+"/room/", "application/json", body)
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
