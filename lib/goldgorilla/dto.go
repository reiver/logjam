package goldgorilla

type RoomPeerDTO struct {
	RoomId string `json:"roomId"`
	ID     uint64 `json:"id"`
}

type JoinReqModel struct {
	RoomId string `json:"roomId"`
}

// not the rtp direction like rtpsender or rtpreceiver; but the connection direction wether conn wants to send or recv
type ConnectionDirection string

const (
	CDSend = "sendOnly"
	CDRecv = "recvOnly"
)

type SendIceCandidateReqModel struct {
	RoomPeerDTO
	GGID          uint64              `json:"ggid"`
	ICECandidate  interface{}         `json:"iceCandidate"`
	ConnDirection ConnectionDirection `json:"connectionDirection"`
}

type CreatePeerRPCModel struct {
	RoomPeerDTO
	CanPublish    bool                `json:"canPublish"`
	IsCaller      bool                `json:"isCaller"`
	GGID          uint64              `json:"ggid"`
	ConnDirection ConnectionDirection `json:"connectionDirection"`
}

type SetSDPRPCModel struct {
	RoomPeerDTO
	GGID          uint64              `json:"ggid"`
	SDP           interface{}         `json:"sdp"`
}
