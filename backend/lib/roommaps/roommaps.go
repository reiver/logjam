package roommaps

import (
	"github.com/mmcomp/go-binarytree"
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

	for _, node := range room.Room.All() {
		nodeStreamId, ok := node.(*binarytreesrv.MySocket).MetaData["streamId"]
		if ok {
			if nodeStreamId == streamId {
				return node, nil
			}
		}
	}

	return nil, errNodeNotFound
}
