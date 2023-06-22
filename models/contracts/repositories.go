package contracts

import "github.com/sparkscience/logjam/models"

type IRoomRepository interface {
	DoesRoomExists(id string) bool
	CreateRoom(id string) error
	GetRoom(id string) (*models.RoomModel, error)
	SetBroadcaster(roomId string, id uint64, name string, stream string) error
	GetBroadcaster(roomId string) (*models.UserModel, error)
	ClearBroadcasterSeat(roomId string) error
	AddMember(roomId string, id uint64, name, email, streamId string) error
	GetMember(roomId string, id uint64) (*models.UserModel, error)
	UpdateCanConnect(roomId string, id uint64, newState bool) error
	UpdateMemberMeta(roomId string, id uint64, metaKey string, value string) error
	GetAllMembersId(roomId string, excludeBroadcaster bool) ([]uint64, error)
	InsertMemberToTree(roomId string, audienceId uint64) (parentId *uint64, err error)
}
