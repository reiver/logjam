package goldgorilla

type IGoldGorillaServiceRepository interface {
	Start(roomId string) error
	ResetRoom(roomId string) (*uint64, error)
	CreatePeer(roomId string, id uint64, canPublish bool, isCaller bool, ggId uint64, connDirection ConnectionDirection) error
	SendICECandidate(roomId string, id uint64, iceCandidate interface{}, connDirection ConnectionDirection) error
	SendAnswer(roomId string, peerId uint64, answer interface{}) error
	SendOffer(roomId string, peerId uint64, offer interface{}) error
	ClosePeer(roomId string, id uint64) error
}
