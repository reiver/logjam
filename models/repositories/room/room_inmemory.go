package roomRepository

import (
	"errors"
	"github.com/sparkscience/logjam/models"
	"github.com/sparkscience/logjam/models/contracts"
	"github.com/sparkscience/logjam/models/dto"
	"strconv"
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
		PeersTree: &models.PeerModel{},
		Members:   make(map[uint64]*models.MemberModel),
		MetaData:  make(map[string]interface{}),
	}
	return nil
}

func (r *roomRepository) HadGoldGorillaInTreeBefore(id string) bool {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(id) {
		return false
	}
	return r.rooms[id].HadGoldGorillaBefore
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

func (r *roomRepository) SetBroadcaster(roomId string, id uint64) error {
	r.Lock()
	defer r.Unlock()
	if r.doesRoomExists(roomId) {
		r.rooms[roomId].Lock()
		defer r.rooms[roomId].Unlock()
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
		MetaData:       map[string]interface{}{"streamId": streamId},
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
			return errors.New("member doesn't exists " + strconv.FormatUint(id, 10))
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

func (r *roomRepository) InsertMemberToTree(roomId string, memberId uint64, isGoldGorilla bool) (parentId *uint64, err error) {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return nil, errors.New("room doesn't exists")
	}
	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
	if r.rooms[roomId].GoldGorilla != nil && isGoldGorilla {
		gid := (*r.rooms[roomId].GoldGorilla).ID
		go r.RemoveMember(roomId, gid)
		r.rooms[roomId].GoldGorilla = nil
	}
	if r.rooms[roomId].GoldGorilla != nil {
		(*r.rooms[roomId].GoldGorilla).Children = append((*r.rooms[roomId].GoldGorilla).Children, &models.PeerModel{
			ID:            memberId,
			IsConnected:   true,
			Children:      []*models.PeerModel{},
			IsGoldGorilla: false,
		})
		parentId = &(*r.rooms[roomId].GoldGorilla).ID
		return parentId, nil
	}
	lastCheckedLevel := uint(0)
start:
	levelNodes, err := r.rooms[roomId].GetLevelMembers(lastCheckedLevel, false)
	if err != nil {
		return nil, err
	}
	lastCheckedLevel++
	if len(levelNodes) == 0 {
		return nil, errors.New("no node to connect to")
	}
	found := false
	for _, node := range levelNodes {
		if len((*node).Children) < 2 {
			newChild := &models.PeerModel{
				ID:            memberId,
				IsConnected:   true,
				Children:      []*models.PeerModel{},
				IsGoldGorilla: isGoldGorilla,
			}
			(*node).Children = append((*node).Children, newChild)
			parentId = &(*node).ID
			found = true
			if isGoldGorilla {
				r.rooms[roomId].GoldGorilla = &newChild
				r.rooms[roomId].HadGoldGorillaBefore = true
			}
		}
		if found {
			break
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

func (r *roomRepository) SetRoomMetaData(roomId string, metaData map[string]interface{}) error {
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
func (r *roomRepository) AddMessageToHistory(roomId string, senderId uint64, msg string) error {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return errors.New("room doesn't exists")
	}

	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()

	if _, exists := r.rooms[roomId].MetaData[models.RoomMessagesMetaDataKey]; !exists {
		r.rooms[roomId].MetaData[models.RoomMessagesMetaDataKey] = []dto.UserMessageModel{}
	}
	lastMessages := r.rooms[roomId].MetaData[models.RoomMessagesMetaDataKey].([]dto.UserMessageModel)
	r.rooms[roomId].MetaData[models.RoomMessagesMetaDataKey] = append(lastMessages, dto.UserMessageModel{
		Message:  msg,
		SenderId: senderId,
	})

	return nil
}

func (r *roomRepository) ClearMessageHistory(roomId string) error {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return errors.New("room doesn't exists")
	}

	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()

	r.rooms[roomId].MetaData[models.RoomMessagesMetaDataKey] = []dto.UserMessageModel{}
	return nil
}

func (r *roomRepository) GetRoomMetaData(roomId string) (map[string]interface{}, error) {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return nil, errors.New("room doesn't exists")
	}

	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
	copiedMap := make(map[string]interface{})
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

func (r *roomRepository) RemoveMember(roomId string, memberId uint64) (wasBroadcaster bool, nodeChildrenIdList []uint64, err error) {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return false, nil, errors.New("room doesn't exists")
	}
	defer func() {
		delete(r.rooms[roomId].Members, memberId)
	}()
	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
	if r.rooms[roomId].PeersTree.ID == memberId {
		if len(r.rooms[roomId].PeersTree.Children) > 0 {
			nodeChildrenIdList = append(nodeChildrenIdList, r.rooms[roomId].PeersTree.Children[0].ID)
		}
		if len(r.rooms[roomId].PeersTree.Children) > 1 {
			nodeChildrenIdList = append(nodeChildrenIdList, r.rooms[roomId].PeersTree.Children[1].ID)
		}
		r.rooms[roomId].PeersTree.IsConnected = false
		return true, nodeChildrenIdList, nil
	}
	var lastNodesList []**models.PeerModel
	//var targetNode ***models.PeerModel
	lastCheckedLevel := uint(0)

	found := false
start:
	lastNodesList, err = r.rooms[roomId].GetLevelMembers(lastCheckedLevel, true)
	if len(lastNodesList) == 0 {
		return false, nodeChildrenIdList, nil
	}
	lastCheckedLevel++
	for _, node := range lastNodesList {
		for i, nodeChild := range (*node).Children {
			if nodeChild.ID == memberId {
				for _, parentLostChild := range nodeChild.Children {
					nodeChildrenIdList = append(nodeChildrenIdList, parentLostChild.ID)
				}
				(*node).Children = append((*node).Children[:i], (*node).Children[i+1:]...)
				if memberId == models.GoldGorillaId {
					r.rooms[roomId].GoldGorilla = nil
				}
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		goto start
	}
	return false, nodeChildrenIdList, nil
}

func (r *roomRepository) GetChildrenIdList(roomId string, id uint64) ([]uint64, error) {
	r.Lock()
	defer r.Unlock()
	if !r.doesRoomExists(roomId) {
		return nil, errors.New("room doesn't exists")
	}

	r.rooms[roomId].Lock()
	defer r.rooms[roomId].Unlock()
	var list []uint64
	lastCheckedLevel := 0
start:
	levelNodes, err := r.rooms[roomId].GetLevelMembers(uint(lastCheckedLevel), true)
	if err != nil {
		return nil, err
	}
	if len(levelNodes) == 0 {
		return nil, nil
	}
	lastCheckedLevel++
	for _, node := range levelNodes {
		if (*node).ID == id {
			if (*node).Children[0] != nil {
				list = append(list, (*node).Children[0].ID)
			}
			if (*node).Children[1] != nil {
				list = append(list, (*node).Children[1].ID)
			}
			break
		}
	}
	if len(list) == 0 {
		goto start
	}

	return list, nil
}
