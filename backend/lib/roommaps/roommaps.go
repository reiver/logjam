package roommaps

import (
	"encoding/json"

	"github.com/mmcomp/go-binarytree"
	"github.com/sparkscience/logjam/backend/lib/message"
	binarytreesrv "github.com/sparkscience/logjam/backend/srv/binarytree"

	"sync"
)

type RoomType struct {
	Room     *binarytree.Tree
	MetaData map[string]string
}

type Type struct {
	mutex    sync.Mutex
	roomMaps map[string]*RoomType
}

type User struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	StreamId string `json:"streamId"`
	Quality  string `json:"quality"`
}

func (receiver *Type) Get(roomName string) (*RoomType, bool) {
	if nil == receiver {
		return nil, false
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if nil == receiver.roomMaps {
		return nil, false
	}

	mapptr, found := receiver.roomMaps[roomName]
	if !found {
		return nil, false
	}
	return mapptr, true
}

func (receiver *Type) Set(roomName string, mapptr *binarytree.Tree) error {
	if nil == receiver {
		return errNilReceiver
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if nil == receiver.roomMaps && nil != mapptr {
		receiver.roomMaps = make(map[string]*RoomType)
	}

	var metaData = make(map[string]string)
	_, ok := receiver.roomMaps[roomName]
	if ok {
		metaData = receiver.roomMaps[roomName].MetaData
	}
	if nil == mapptr {
		delete(receiver.roomMaps, roomName)
		return nil
	}

	receiver.roomMaps[roomName] = &RoomType{
		Room:     mapptr,
		MetaData: metaData,
	}
	return nil
}

func (receiver *Type) GetFromMetaData(roomName string, key string) (*string, error) {
	if nil == receiver {
		return nil, errNilReceiver
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if receiver.roomMaps == nil {
		return nil, nil
	}

	room, exists := receiver.roomMaps[roomName]
	if !exists {
		return nil, errRoomNotFound
	}

	if metaValue, exists := room.MetaData[key]; exists {
		return &metaValue, nil
	}

	return nil, nil
}

func (receiver *Type) SetToMetaData(roomName string, key, value string) error {
	if nil == receiver {
		return errNilReceiver
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if nil == receiver.roomMaps {
		receiver.roomMaps = make(map[string]*RoomType)
	}

	room, ok := receiver.roomMaps[roomName]
	if !ok {
		return errRoomNotFound
	}
	room.MetaData[key] = value
	receiver.roomMaps[roomName] = room
	return nil
}

func (receiver *Type) DelFromMetaData(roomName string, key string) error {
	if nil == receiver {
		return errNilReceiver
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if nil == receiver.roomMaps {
		receiver.roomMaps = make(map[string]*RoomType)
	}

	room, ok := receiver.roomMaps[roomName]
	if !ok {
		return errRoomNotFound
	}
	delete(room.MetaData, key)
	receiver.roomMaps[roomName] = room
	return nil
}

func (receiver *Type) SetMetData(roomName string, metaData map[string]string) error {
	if nil == receiver {
		return errNilReceiver
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if nil == receiver.roomMaps {
		receiver.roomMaps = make(map[string]*RoomType)
	}

	room, ok := receiver.roomMaps[roomName]
	if !ok {
		return errRoomNotFound
	}
	room.MetaData = metaData
	receiver.roomMaps[roomName] = room
	return nil
}

func (receiver *Type) GetMetaDataJson(roomName string) (*string, error) {
	if nil == receiver {
		return nil, errNilReceiver
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if nil == receiver.roomMaps {
		receiver.roomMaps = make(map[string]*RoomType)
	}
	room, ok := receiver.roomMaps[roomName]
	if !ok {
		return nil, errRoomNotFound
	}
	bytes, err := json.Marshal(room.MetaData)
	if err != nil {
		return nil, err
	}
	jsonStr := string(bytes)
	return &jsonStr, nil
}

func (receiver *Type) GetSocketByStreamId(roomName, streamId string) (binarytree.SingleNode, error) {
	if nil == receiver {
		return nil, errNilReceiver
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if nil == receiver.roomMaps {
		receiver.roomMaps = make(map[string]*RoomType)
	}

	room, ok := receiver.roomMaps[roomName]
	if !ok {
		return nil, errRoomNotFound
	}

	for _, node := range room.Room.Nodes() {
		nodeStreamId, ok := node.(*binarytreesrv.MySocket).MetaData["streamId"]
		if ok {
			if nodeStreamId == streamId {
				return node, nil
			}
		}
	}

	return nil, errNodeNotFound
}

func (receiver *Type) GetUsers(roomName string) ([]User, error) {
	output := []User{}
	if nil == receiver {
		return nil, errRoomNotFound
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if nil == receiver.roomMaps {
		receiver.roomMaps = make(map[string]*RoomType)
	}

	room, ok := receiver.roomMaps[roomName]
	if !ok {
		return nil, errRoomNotFound
	}

	for _, node := range room.Room.Nodes() {
		nodeStreamId := node.(*binarytreesrv.MySocket).MetaData["streamId"]
		role := "audience"
		if node.(*binarytreesrv.MySocket).IsBroadcaster {
			role = "broadcaster"
		}

		if node.(*binarytreesrv.MySocket).Socket == nil {
			receiver.roomMaps[roomName].Room.Delete(node)
		} else {
			message := message.MessageContract{
				Type: "pong",
				Data: "pong",
			}
			messageTxt, _ := json.Marshal(message)
			err := node.(*binarytreesrv.MySocket).Writer.WriteMessage(1, messageTxt)
			if err != nil {
				// receiver.roomMaps[roomName].Room.Delete(node)
			} else {
				user := User{
					Id:       node.(*binarytreesrv.MySocket).ID,
					Name:     node.(*binarytreesrv.MySocket).Name,
					StreamId: nodeStreamId,
					Role:     role,
					Quality:  node.(*binarytreesrv.MySocket).MetaData["quality"],
				}
				output = append(output, user)
			}
		}
	}

	return output, nil
}
