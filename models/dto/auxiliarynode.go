package dto

type RoomPeerDTO struct {
	RoomId string `json:"roomId"`
	ID     uint64 `json:"id"`
}

type JoinReqModel struct {
	RoomId      string `json:"roomId"`
	ServiceAddr string `json:"svcAddr"`
}

type SendIceCandidateReqModel struct {
	RoomPeerDTO
	ICECandidate any `json:"iceCandidate"`
}

type CreatePeerRPCModel struct {
	RoomPeerDTO
	CanPublish bool `json:"canPublish"`
	IsCaller   bool `json:"isCaller"`
}

type SetSDPRPCModel struct {
	RoomPeerDTO
	SDP any `json:"sdp"`
}
