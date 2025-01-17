package rooms

type Repository interface {
	NumRooms() int
	RoomIDs() []string
	ForEachRoom(func(*RoomModel)error) error

	DoesRoomExists(id string) bool
	CreateRoom(id string) error
	GetRoom(id string) (*RoomModel, error)
	SetBroadcaster(roomId string, id uint64) error
	GetBroadcaster(roomId string) (*MemberModel, error)
	ClearBroadcasterSeat(roomId string) error
	AddMember(roomId string, id uint64, name, email, streamId string, isGoldGorilla bool) error
	GetMember(roomId string, id uint64) (*MemberModel, error)
	UpdateCanConnect(roomId string, id uint64, newState bool) error
	UpdateTurnStatus(roomId string, id uint64, newState bool) error
	UpdateMemberMeta(roomId string, id uint64, metaKey string, value string) error
	UpdateMemberName(roomId string, id uint64, name string) error
	GetAllMembersId(roomId string, excludeBroadcaster bool) ([]uint64, error)
	InsertMemberToTree(roomId string, memberId uint64, isGoldGorilla bool) (parentId *uint64, err error)
	RemoveMember(roomId string, memberId uint64) (wasBroadcaster bool, nodeChildrenIdList []uint64, err error)
	SetRoomMetaData(roomId string, metaData map[string]any) error
	GetRoomMetaData(roomId string) (map[string]any, error)
	AddMessageToHistory(roomId string, senderId uint64, msg string) error
	ClearMessageHistory(roomId string) error
	GetUserByStreamId(roomId string, streamId string) (*MemberModel, error)
	IsBroadcaster(roomId string, id uint64) (bool, error)
	GetMembersList(roomId string) ([]MemberDTO, error)
	GetChildrenIdList(roomId string, id uint64) ([]uint64, error)
	IsGGInstance(roomId string, id uint64) bool
	GetRoomGoldGorillaId(roomId string) (*uint64, error)
}
