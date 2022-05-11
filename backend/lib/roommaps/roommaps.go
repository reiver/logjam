package roommaps

import (
	"github.com/mmcomp/go-binarytree"

	"sync"
)

type Type struct {
	mutex sync.Mutex
	roomMaps map[string]*binarytree.Tree
}

func (receiver *Type) Get(roomName string) (*binarytree.Tree, bool) {
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
		receiver.roomMaps = make(map[string]*binarytree.Tree)
	}

	if nil == mapptr {
		delete(receiver.roomMaps, roomName)
		return nil
	}

	receiver.roomMaps[roomName] = mapptr
	return nil
}
