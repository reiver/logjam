package contracts

import (
	"sourcecode.social/greatape/logjam/models"
	"sourcecode.social/greatape/logjam/models/dto"
)

type IRoomRepository interface {
	DoesRoomExists(id string) bool
	CreateRoom(id string) error
	GetRoom(id string) (*models.RoomModel, error)
	SetBroadcaster(roomId string, id uint64) error
	GetBroadcaster(roomId string) (*models.MemberModel, error)
	ClearBroadcasterSeat(roomId string) error
	AddMember(roomId string, id uint64, name, email, streamId string) error
	GetMember(roomId string, id uint64) (*models.MemberModel, error)
	UpdateCanConnect(roomId string, id uint64, newState bool) error
	UpdateTurnStatus(roomId string, id uint64, newState bool) error
	UpdateMemberMeta(roomId string, id uint64, metaKey string, value string) error
	UpdateMemberName(roomId string, id uint64, name string) error
	GetAllMembersId(roomId string, excludeBroadcaster bool) ([]uint64, error)
	InsertMemberToTree(roomId string, memberId uint64, isGoldGorilla bool) (parentId *uint64, err error)
	RemoveMember(roomId string, memberId uint64) (wasBroadcaster bool, nodeChildrenIdList []uint64, err error)
	SetRoomMetaData(roomId string, metaData map[string]interface{}) error
	GetRoomMetaData(roomId string) (map[string]interface{}, error)
	AddMessageToHistory(roomId string, senderId uint64, msg string) error
	ClearMessageHistory(roomId string) error
	GetUserByStreamId(roomId string, streamId string) (*models.MemberModel, error)
	IsBroadcaster(roomId string, id uint64) (bool, error)
	GetMembersList(roomId string) ([]dto.MemberDTO, error)
	GetChildrenIdList(roomId string, id uint64) ([]uint64, error)
	HadGoldGorillaInTreeBefore(id string) bool
}

type IGoldGorillaServiceRepository interface {
	Init(svcAddr string) error
	Start() error
	ResetRoom(roomId string) error
	CreatePeer(roomId string, id uint64, canPublish bool, isCaller bool) error
	SendICECandidate(roomId string, id uint64, iceCandidate interface{}) error
	SendAnswer(roomId string, peerId uint64, answer interface{}) error
	SendOffer(roomId string, peerId uint64, offer interface{}) error
	ClosePeer(roomId string, id uint64) error
}
