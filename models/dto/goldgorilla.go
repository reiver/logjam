package dto

type RoomPeerDTO struct {
	RoomId string `json:"roomId"`
	ID     uint64 `json:"id"`
}

type JoinReqModel struct {
	RoomId string `json:"roomId"`
}

type SendIceCandidateReqModel struct {
	RoomPeerDTO
	GGID         uint64      `json:"ggid"`
	ICECandidate interface{} `json:"iceCandidate"`
}

type CreatePeerRPCModel struct {
	RoomPeerDTO
	CanPublish bool   `json:"canPublish"`
	IsCaller   bool   `json:"isCaller"`
	GGID       uint64 `json:"ggid"`
}

type SetSDPRPCModel struct {
	RoomPeerDTO
	GGID uint64      `json:"ggid"`
	SDP  interface{} `json:"sdp"`
}
