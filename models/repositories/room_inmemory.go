package repositories

import (
	"errors"
	"github.com/sparkscience/logjam/models"
	"github.com/sparkscience/logjam/models/contracts"
	"sync"
)

type roomRepository struct {
	*sync.Mutex
	rooms map[string]*models.RoomModel
}

func NewRoomRepository() contracts.IRoomRepository {
	return &roomRepository{
		Mutex: &sync.Mutex{},
		rooms: make(map[string]*models.RoomModel),
	}
}

func (r *roomRepository) doesRoomExists(id string) bool {
	if _, exists := r.rooms[id]; exists {
		return true
	}
	return false
}

func (r *roomRepository) DoesRoomExists(id string) bool {
	r.Lock()
	defer r.Unlock()
	return r.doesRoomExists(id)
}

func (r *roomRepository) CreateRoom(id string) error {
	r.Lock()
	defer r.Unlock()
	if r.doesRoomExists(id) {
		return nil
	}

	r.rooms[id] = &models.RoomModel{
		Mutex:     &sync.Mutex{},
		Title:     "",
		PeersTree: &models.PeerTreeModel{},
		Members:   make(map[uint64]*models.UserModel),
		MetaData:  make(map[string]any),
	}
	return nil
}

func (r *roomRepository) GetRoom(id string) (*models.RoomModel, error) {
	r.Lock()
	defer r.Unlock()
	if room, exists := r.rooms[id]; exists {
		return room, nil
	} else {
		return nil, nil
	}
}

func (r *roomRepository) SetBroadcaster(roomId string, id uint64, name string, stream string) error {
	r.Lock()
	defer r.Unlock()
	email := ""
	if r.doesRoomExists(roomId) {
		r.rooms[roomId].Members[id] = &models.UserModel{
			ID:       id,
			Name:     name,
			Email:    email,
			StreamId: stream,
			MetaData: make(map[string]any),
		}
		r.rooms[roomId].PeersTree.ID = id
		r.rooms[roomId].PeersTree.IsConnected = true
	} else {
		return errors.New("room does not exists")
	}
	return nil
}

func (r *roomRepository) GetBroadcaster(roomId string) (*models.UserModel, error) {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return nil, errors.New("room doesnt exists")
	}
	if r.rooms[roomId].PeersTree.IsConnected {
		return r.rooms[roomId].Members[r.rooms[roomId].PeersTree.ID], nil
	} else {
		return nil, nil
	}
}

func (r *roomRepository) ClearBroadcasterSeat(roomId string) error {
	r.Lock()
	defer r.Unlock()
	if r.doesRoomExists(roomId) {
		r.rooms[roomId].PeersTree.ID = 0
		r.rooms[roomId].PeersTree.IsConnected = false
	}
	return nil
}

func (r *roomRepository) AddMember(roomId string, id uint64, name, email, streamId string) error {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return errors.New("room doesnt' exists")
	}
	r.rooms[roomId].Members[id] = &models.UserModel{
		ID:             id,
		Name:           name,
		Email:          email,
		StreamId:       streamId,
		MetaData:       make(map[string]any),
		CanAcceptChild: false,
	}
	return nil
}

func (r *roomRepository) GetMember(roomId string, id uint64) (*models.UserModel, error) {
	r.Lock()
	defer r.Unlock()

	if r.doesRoomExists(roomId) {
		if user, exists := r.rooms[roomId].Members[id]; exists {
			return user, nil
		} else {
			return nil, nil
		}
	} else {
		return nil, errors.New("room doesnt exists")
	}
}

func (r *roomRepository) UpdateMemberMeta(roomId string, id uint64, metaKey string, value string) error {
	r.Lock()
	defer r.Unlock()

	if r.doesRoomExists(roomId) {
		if user, exists := r.rooms[roomId].Members[id]; exists {
			user.MetaData[metaKey] = value
		} else {
			return errors.New("audience doesn't exists")
		}
	} else {
		return errors.New("room doesnt exists")
	}
	return nil
}

func (r *roomRepository) GetAllMembersId(roomId string, excludeBroadcaster bool) ([]uint64, error) {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return nil, errors.New("room doesn't exists")
	}

	memberIds := make([]uint64, 0, len(r.rooms[roomId].Members))
	for id := range r.rooms[roomId].Members {
		if excludeBroadcaster && id == r.rooms[roomId].PeersTree.ID {
			continue
		}
		memberIds = append(memberIds, id)
	}
	return memberIds, nil
}

func (r *roomRepository) UpdateCanConnect(roomId string, id uint64, newState bool) error {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return errors.New("room doesn't exists")
	}
	if _, exists := r.rooms[roomId].Members[id]; !exists {
		return errors.New("no such a member in this room")
	}
	r.rooms[roomId].Members[id].CanAcceptChild = newState
	return nil
}

func (r *roomRepository) InsertMemberToTree(roomId string, audienceId uint64) (parentId *uint64, err error) {
	r.rooms[roomId].PeersTree.Children[0] = &models.PeerTreeModel{
		ID:          audienceId,
		IsConnected: true,
		Children:    [2]*models.PeerTreeModel{},
	}
	return &r.rooms[roomId].PeersTree.ID, nil
}
