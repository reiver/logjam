package room

import (
	cErrors "github.com/reiver/logjam/lib/errors"
	"github.com/reiver/logjam/lib/marshal"
	dbsrv "github.com/reiver/logjam/srv/db"
	"net/http"
	"strings"
)

type roomRepo struct {
}

const (
	roomTbl = "rooms"

	UIDKey     = "UID"
	ownerIdKey = "ownerId"
)

func NewRoomRepository() IRoomRepository {
	return &roomRepo{}
}

func (r *roomRepo) CreateRoom(dto CreateRoomDTO) error {
	if len(dto.UID) < 6 || len(dto.UID) > 32 {
		return cErrors.NewErrorWithMsg(http.StatusBadRequest, "invalid UID len. min: 6, max: 32")
	}
	room, _ := r.GetRoom(strings.ToLower(dto.UID))
	if room != nil && len(room.UID) > 0 {
		return cErrors.NewErrorWithMsg(http.StatusConflict, "this room UID is taken")
	}
	dto.UID = strings.ToLower(dto.UID)
	data, err := marshal.ObjToMap(dto)
	if err != nil {
		return err
	}
	_, err = dbsrv.Repository.Insert(roomTbl, data)
	return err
}

func (r *roomRepo) GetRoom(UID string) (result *RoomDTO, err error) {
	rows, err := dbsrv.Repository.GetByFilter(roomTbl, map[string]any{
		UIDKey: UID,
	})
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
	room, err := r.GetRoom(dto.UID)
	if err != nil {
		return err
	}
	if room.OwnerID != dto.OwnerID {
		return cErrors.NewError(http.StatusForbidden)
	}
	data, err := marshal.ObjToMap(dto)
	if err != nil {
		return err
	}
	return dbsrv.Repository.UpdateByFilter(roomTbl, map[string]any{
		UIDKey:     dto.UID,
		ownerIdKey: dto.OwnerID,
	}, data)
}

func (r *roomRepo) DeleteRoom(UID, ownerId string) error {
	room, err := r.GetRoom(UID)
	if err != nil {
		return err
	}
	if room.OwnerID != ownerId {
		return cErrors.NewError(http.StatusForbidden)
	}
	return dbsrv.Repository.DeleteByFilter(roomTbl, map[string]any{
		UIDKey:     UID,
		ownerIdKey: ownerId,
	})
}

func (r *roomRepo) GetUserRooms(UserId string) (result []RoomDTO, err error) {
	rows, err := dbsrv.Repository.GetByFilter(roomTbl, map[string]any{
		ownerIdKey: UserId,
	})
	if err != nil {
		return nil, err
	}
	err = marshal.MapArrayToObjArray(rows, &result)
	return
}
