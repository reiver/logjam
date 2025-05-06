package room

type IRoomRepository interface {
	CreateRoom(CreateRoomDTO) error
	GetRoom(UID, ownerId string) (*RoomDTO, error)
	UpdateRoom(UpdateRoomDTO) error
	DeleteRoom(UID, ownerId string) error
	GetUserRooms(UserId string) ([]RoomDTO, error)
}
