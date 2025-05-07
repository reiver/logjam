package room

type IRoomRepository interface {
	CreateRoom(CreateRoomDTO) error
	GetRoom(UID string) (*RoomDTO, error)
	UpdateRoom(UpdateRoomDTO) error
	DeleteRoom(UID, ownerId string) error
	GetUserRooms(UserId string) ([]RoomDTO, error)
}
