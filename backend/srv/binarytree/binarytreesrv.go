package binarytreesrv

import (
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mmcomp/go-binarytree"
)

var (
	Map = binarytree.Default
)

type MySocket struct {
	mutex            sync.Mutex
	Socket           *websocket.Conn
	ID               uint64
	IsBroadcaster    bool
	Name             string
	HasStream        bool
	IsTURN           bool
	ConnectedSockets map[*websocket.Conn]MySocket
	MetaData         map[string]string
}

func (receiver *MySocket) Insert(node binarytree.SingleNode) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	socket := node.(*MySocket)
	receiver.ConnectedSockets[socket.Socket] = *socket
}

func (receiver *MySocket) Delete(node interface{}) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	delete(receiver.ConnectedSockets, node.(*websocket.Conn))
}

func (receiver *MySocket) Get(node interface{}) binarytree.SingleNode {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	result := receiver.ConnectedSockets[node.(*websocket.Conn)]
	return &result
}

func (receiver *MySocket) Length() int {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	return len(receiver.ConnectedSockets)
}

func (receiver *MySocket) IsHead() bool {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	return receiver.IsBroadcaster
}

func (receiver *MySocket) CanConnect() bool {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	return receiver.HasStream
}

func (receiver *MySocket) All() map[interface{}]binarytree.SingleNode {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	var output map[interface{}]binarytree.SingleNode = make(map[interface{}]binarytree.SingleNode)
	for indx := range receiver.ConnectedSockets {
		result := receiver.ConnectedSockets[indx]
		output[indx] = &result
	}
	return output
}

func (receiver *MySocket) ToggleHead() {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.IsBroadcaster = !receiver.IsBroadcaster
}

func (receiver *MySocket) ToggleCanConnect() {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.HasStream = !receiver.HasStream
}

func (receiver *MySocket) Index() interface{} {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	return receiver.Socket
}

func (receiver *MySocket) SetName(name string) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.Name = name
}

func (receiver *MySocket) SetIsTurn(isTurn bool) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.IsTURN = isTurn
}

func (receiver *MySocket) SetMetaData(metaData map[string]string) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.MetaData = metaData
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

func Tree(aMap binarytree.Tree) []TreeGraphElement {
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
