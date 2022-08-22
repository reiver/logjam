package node

import (
	"fmt"
	"reflect"

	"github.com/gorilla/websocket"
)

type Node struct {
	id             uint64
	name           string
	isBroadcaster  bool
	hasStream      bool
	isTURN         bool
	socket         *websocket.Conn
	connectedNodes map[*websocket.Conn]Node
	r              *Room
}

func NewNode(id uint64, name string, isBroadcaster bool, hasStream bool, isTURN bool, socket *websocket.Conn) *Node {
	return &Node{
		id:             id,
		name:           name,
		isBroadcaster:  isBroadcaster,
		hasStream:      hasStream,
		isTURN:         isTURN,
		socket:         socket,
		connectedNodes: make(map[*websocket.Conn]Node),
	}
}

func (n *Node) Read() {
	defer n.socket.Close()
	for {
		_, msg, err := n.socket.ReadMessage()
		if err != nil {
			break
		}
		message, err := UnmarshalJSON(msg)
		if err != nil {
			fmt.Println("ERROR ", err)
			continue
		}

		fmt.Println("MESSAGE ", message, " ", reflect.TypeOf(message).Elem().Name())
		if reflect.TypeOf(message).Elem().Name() == "MessageContract[node.JoinMessageData]" {
			message.(*MessageContract[JoinMessageData]).handle()
		}

		fmt.Println("Node ", n.id, " read:", string(msg))
	}
}
