package binarytreesrv

import (
	"strconv"
	"time"

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
	IsTURN           bool
	ConnectedSockets map[*websocket.Conn]MySocket
}

func (receiver *MySocket) Insert(node binarytree.SingleNode) {
	socket := node.(*MySocket)
	receiver.ConnectedSockets[socket.Socket] = *socket
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

func (receiver *MySocket) SetName(name string) {
	receiver.Name = name
}

func (receiver *MySocket) SetIsTurn(isTurn bool) {
	receiver.IsTURN = isTurn
}

func fillFunction(node interface{}, socketIndex uint64) binarytree.SingleNode {
	conn := node.(*websocket.Conn)
	result := MySocket{
		Socket:           conn,
		ID:               socketIndex,
		Name:             "Socket " + strconv.FormatUint(socketIndex, 10),
		IsTURN:           true,
		ConnectedSockets: make(map[*websocket.Conn]MySocket),
	}
	return &result
}

func GetMap() binarytree.Tree {
	Map.SetFillNode(fillFunction)
	return Map
}

func InsertChild(socket *websocket.Conn, aMap binarytree.Tree) (binarytree.SingleNode, error) {
	tryCount := 0
	var result binarytree.SingleNode
	var err error
	result, err = aMap.InsertChild(socket, false)
	for err != nil && err.Error() == "no nodes to connect" && tryCount < 10 {
		result, err = aMap.InsertChild(socket, false)
		time.Sleep(1000 * time.Millisecond)
		tryCount++
	}
	return result, err
}

type TreeGraphElement struct {
	Name     string             `json:"name"`
	Parent   string             `json:"parent"`
	Children []TreeGraphElement `json:"children"`
}

func addSubSockets(socket MySocket, children *[]TreeGraphElement, aMap binarytree.Tree) {
	for child := range socket.ConnectedSockets {
		childSocket := aMap.Get(child).(*MySocket)
		var turnState string = "no-TURN"
		if childSocket.IsTURN {
			turnState = "TURN"
		}
		*children = append(*children, TreeGraphElement{
			Name:     childSocket.Name + "[" + turnState + "]",
			Parent:   "null",
			Children: []TreeGraphElement{},
		})
		addSubSockets(*childSocket, &(*children)[len(*children)-1].Children, aMap)
	}
}

func GetTree(aMap binarytree.Tree) []TreeGraphElement {
	treeData := []TreeGraphElement{}
	broadcasterLevel := aMap.LevelNodes(1)
	if len(broadcasterLevel) == 0 {
		return treeData
	}
	broadcaster := broadcasterLevel[0].(*MySocket)
	var turnState string = "no-TURN"
	if broadcaster.IsTURN {
		turnState = "TURN"
	}
	treeData = append(treeData, TreeGraphElement{

		Name:     broadcaster.Name + "[" + turnState + "]",
		Parent:   "null",
		Children: []TreeGraphElement{},
	})
	addSubSockets(*broadcaster, &treeData[0].Children, aMap)
	return treeData
}
