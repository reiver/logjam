package binarytreesrv

import (
	"fmt"
	"strconv"

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
	fmt.Println("ToggleCanConnect ", receiver.ID)
	receiver.HasStream = !receiver.HasStream
	fmt.Println("HasStream ", receiver.HasStream)
}

func (receiver *MySocket) GetIndex() interface{} {
	return receiver.Socket
}

func fillFunction(node interface{}, socketIndex uint64) binarytree.SingleNode {
	conn := node.(*websocket.Conn)
	result := MySocket{
		Socket:           conn,
		ID:               socketIndex,
		Name:             "Socket " + strconv.FormatUint(socketIndex, 10),
		ConnectedSockets: make(map[*websocket.Conn]MySocket),
	}
	return &result
}

func GetMap() binarytree.Tree {
	Map.SetFillNode(fillFunction)
	return Map
}

type TreeGraphElement struct {
	Name     string             `json:"name"`
	Parent   string             `json:"parent"`
	Children []TreeGraphElement `json:"children"`
}

func addSubSockets(socket MySocket, children *[]TreeGraphElement, aMap binarytree.Tree) {
	fmt.Println("Adding CHilds of ", socket.ID, socket.Name)
	for child := range socket.ConnectedSockets {
		childSocket := aMap.Get(child).(*MySocket)
		*children = append(*children, TreeGraphElement{
			Name:     childSocket.Name,
			Parent:   "null",
			Children: []TreeGraphElement{},
		})
		fmt.Println("Child ", childSocket.ID, childSocket.Name)
		addSubSockets(*childSocket, &(*children)[len(*children)-1].Children, aMap)
	}
}
func GetTree(aMap binarytree.Tree) []TreeGraphElement {
	fmt.Println("GetTree()")
	treeData := []TreeGraphElement{}
	broadcasterLevel := aMap.LevelNodes(1)
	fmt.Println("broadcasterLevel ", len(broadcasterLevel))
	broadcaster := broadcasterLevel[0].(*MySocket)
	fmt.Println("Broadcaster ", broadcaster.ID, broadcaster.Name)
	treeData = append(treeData, TreeGraphElement{
		Name:     broadcaster.Name,
		Parent:   "null",
		Children: []TreeGraphElement{},
	})
	addSubSockets(*broadcaster, &treeData[0].Children, aMap)
	return treeData
}
