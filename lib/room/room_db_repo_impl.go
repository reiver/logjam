package room

import (
	"github.com/reiver/logjam/lib/marshal"
	dbsrv "github.com/reiver/logjam/srv/db"
)

type roomRepo struct {
}

const (
	roomTbl = "rooms"

	UIDKey = "UID"
)

func NewRoomRepository() IRoomRepository {
	return &roomRepo{}
}

func (r *roomRepo) CreateRoom(dto CreateRoomDTO) error {
	data, err := marshal.ObjToMap(dto)
	if err != nil {
		return err
	}
	_, err = dbsrv.Repository.Insert(roomTbl, data)
	return err
}

func (r *roomRepo) GetRoom(UID, ownerId string) (result *RoomDTO, err error) {
	rows, err := dbsrv.Repository.GetByFilter(roomTbl, map[string]any{UIDKey: UID})
	if err != nil {
		return nil, err
	}
	if rows == nil || len(rows) == 0 {
		return nil, nil
	}

	err = marshal.MapToObj(rows[0], &result)
	return
}

func (r *roomRepo) UpdateRoom(dto UpdateRoomDTO) error {
	//TODO implement me
	panic("implement me")
}

func (r *roomRepo) DeleteRoom(UID, ownerId string) error {
	//TODO implement me
	panic("implement me")
}

func (r *roomRepo) GetUserRooms(UserId string) ([]RoomDTO, error) {
	return nil, nil
}
