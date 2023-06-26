package repositories

import (
	"errors"
	"github.com/sparkscience/logjam/models"
	"github.com/sparkscience/logjam/models/contracts"
	"github.com/sparkscience/logjam/models/dto"
	"sync"
	"time"
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
		PeersTree: &models.PeerModel{},
		Members:   make(map[uint64]*models.MemberModel),
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

func (r *roomRepository) SetBroadcaster(roomId string, id uint64, name string, streamId string) error {
	r.Lock()
	defer r.Unlock()
	email := ""
	if r.doesRoomExists(roomId) {
		r.rooms[roomId].Lock()
		defer r.rooms[roomId].Unlock()
		r.rooms[roomId].Members[id] = &models.MemberModel{
			ID:       id,
			Name:     name,
			Email:    email,
			MetaData: map[string]any{"streamId": streamId},
		}
		r.rooms[roomId].PeersTree.ID = id
		r.rooms[roomId].PeersTree.IsConnected = true
	} else {
		return errors.New("room does not exists")
	}
	return nil
}

func (r *roomRepository) GetBroadcaster(roomId string) (*models.MemberModel, error) {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return nil, errors.New("room doesnt exists")
	}
	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
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
		r.rooms[roomId].Lock()
		defer r.rooms[roomId].Unlock()
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
	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
	r.rooms[roomId].Members[id] = &models.MemberModel{
		ID:             id,
		Name:           name,
		Email:          email,
		MetaData:       make(map[string]any),
		CanAcceptChild: false,
	}
	return nil
}

func (r *roomRepository) GetMember(roomId string, id uint64) (*models.MemberModel, error) {
	r.Lock()
	defer r.Unlock()

	if r.doesRoomExists(roomId) {
		r.rooms[roomId].Lock()
		defer r.rooms[roomId].Unlock()
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
		r.rooms[roomId].Lock()
		defer r.rooms[roomId].Unlock()
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

	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
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
	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
	if _, exists := r.rooms[roomId].Members[id]; !exists {
		return errors.New("no such a member in this room")
	}
	r.rooms[roomId].Members[id].CanAcceptChild = newState
	return nil
}

func (r *roomRepository) InsertMemberToTree(roomId string, audienceId uint64) (parentId *uint64, err error) {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return nil, errors.New("room doesn't exists")
	}
	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
	tryCount := 0
	lastCheckedLevel := 0
start:

	levelNodes, err := r.rooms[roomId].GetLevelAudiences(uint(lastCheckedLevel))
	if err != nil {
		println(err.Error())
		return
	}
	lastCheckedLevel++
	if len(levelNodes) == 0 && tryCount < 10 {
		time.Sleep(1 * time.Second)
		tryCount++
		lastCheckedLevel = 0
		goto start
	} else if len(levelNodes) == 0 && tryCount > 10 {
		return nil, errors.New("no node to connect to")
	}
	found := false
	for _, node := range levelNodes {
		if node.Children[0] == nil {
			node.Children[0] = &models.PeerModel{
				ID:          audienceId,
				IsConnected: true,
				Children:    [2]*models.PeerModel{},
			}
			parentId = &node.ID
			found = true
		} else if node.Children[1] == nil {
			node.Children[1] = &models.PeerModel{
				ID:          audienceId,
				IsConnected: true,
				Children:    [2]*models.PeerModel{},
			}
			parentId = &node.ID
			found = true
		}
	}
	if !found {
		goto start
	}
	return parentId, nil
}

func (r *roomRepository) UpdateTurnStatus(roomId string, id uint64, newState bool) error {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return errors.New("room doesn't exists")
	}

	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
	if _, exists := r.rooms[roomId].Members[id]; !exists {
		return errors.New("member doesn't exists")
	}
	r.rooms[roomId].Members[id].IsUsingTurn = newState
	return nil
}

func (r *roomRepository) UpdateMemberName(roomId string, id uint64, name string) error {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return errors.New("room doesn't exists")
	}

	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
	if _, exists := r.rooms[roomId].Members[id]; !exists {
		return errors.New("member doesn't exists")
	}
	r.rooms[roomId].Members[id].Name = name
	return nil
}

func (r *roomRepository) SetRoomMetaData(roomId string, metaData map[string]any) error {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return errors.New("room doesn't exists")
	}

	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
	r.rooms[roomId].MetaData = metaData
	return nil
}

func (r *roomRepository) GetRoomMetaData(roomId string) (map[string]any, error) {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return nil, errors.New("room doesn't exists")
	}

	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
	copiedMap := make(map[string]any)
	for k, v := range r.rooms[roomId].MetaData {
		copiedMap[k] = v
	}
	return copiedMap, nil
}

func (r *roomRepository) GetUserByStreamId(roomId string, targetStreamId string) (*models.MemberModel, error) {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return nil, errors.New("room doesn't exists")
	}

	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
	var chosen *models.MemberModel
	for _, member := range r.rooms[roomId].Members {
		if streamId, exists := member.MetaData["streamId"]; exists && streamId != targetStreamId {
			continue
		}
		chosen = member
		break
	}
	if chosen == nil {
		return nil, nil
	}
	return chosen, nil
}

func (r *roomRepository) IsBroadcaster(roomId string, id uint64) (bool, error) {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return false, errors.New("room doesn't exists")
	}

	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
	return r.rooms[roomId].PeersTree.ID == id, nil
}

func (r *roomRepository) GetMembersList(roomId string) ([]dto.MemberDTO, error) {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return nil, errors.New("room doesn't exists")
	}

	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
	var list []dto.MemberDTO
	for _, member := range r.rooms[roomId].Members {
		role := "audience"
		if r.rooms[roomId].PeersTree.ID == member.ID {
			role = "broadcaster"
		}
		streamId := ""
		_streamId, exists := member.MetaData["streamId"]
		if exists {
			strStreamId, ok := _streamId.(string)
			if ok {
				streamId = strStreamId
			}
		}
		list = append(list, dto.MemberDTO{
			Id:       member.ID,
			Name:     member.Name,
			Role:     role,
			StreamId: streamId,
		})
	}
	return list, nil
}
