package binarytreesrv

import (
	"github.com/gorilla/websocket"
	"github.com/mmcomp/go-binarytree"
)

var (
	Map = binarytree.Default
)

type MySocket struct {
	Socket           *websocket.Conn
	ID               uint64
	IsBroadcaster    bool
	Name             string
	HasStream        bool
	ConnectedSockets map[*websocket.Conn]MySocket
}

func (receiver *MySocket) Insert(node interface{}, socketIndex uint64) {
	receiver.ConnectedSockets[node.(*websocket.Conn)] = MySocket{
		Socket:           node.(*websocket.Conn),
		ID:               socketIndex,
		IsBroadcaster:    false,
		HasStream:        false,
		ConnectedSockets: make(map[*websocket.Conn]MySocket),
	}
}

func (receiver *MySocket) Delete(node interface{}) {
	delete(receiver.ConnectedSockets, node.(*websocket.Conn))
}

func (receiver *MySocket) Get(node interface{}) binarytree.SingleNode {
	result := receiver.ConnectedSockets[node.(*websocket.Conn)]
	return &result
}

func (receiver *MySocket) GetLength() int {
	return len(receiver.ConnectedSockets)
}

func (receiver *MySocket) IsHead() bool {
	return receiver.IsBroadcaster
}

func (receiver *MySocket) CanConnect() bool {
	return receiver.HasStream
}

func (receiver *MySocket) GetAll() map[interface{}]binarytree.SingleNode {
	var output map[interface{}]binarytree.SingleNode = make(map[interface{}]binarytree.SingleNode)
	for indx := range receiver.ConnectedSockets {
		result := receiver.ConnectedSockets[indx]
		output[indx] = &result
	}
	return output
}

func (receiver *MySocket) ToggleHead() {
	receiver.IsBroadcaster = !receiver.IsBroadcaster
}

func (receiver *MySocket) ToggleCanConnect() {
	receiver.HasStream = !receiver.HasStream
}

func (receiver *MySocket) GetIndex() interface{} {
	return receiver.Socket
}
