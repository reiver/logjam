package websocketmap

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type MySocket struct {
	Socket           *websocket.Conn
	ID               uint64
	IsBroadcaster    bool
	ConnectedSockets map[*websocket.Conn]MySocket
	Name             string
}

var ConnectedSocketsIndex uint64 = 0

type Type struct {
	mutex       sync.Mutex
	connections map[*websocket.Conn]MySocket
}

var Map Type = Type{}

func (receiver *Type) Get(conn *websocket.Conn) MySocket {
	return receiver.connections[conn]
}

func (receiver *Type) Delete(conn *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	fmt.Println("Deleteing ", receiver.connections[conn].ID)
	for socket := range receiver.connections {
		delete(receiver.connections[socket].ConnectedSockets, conn)
	}
	for socket := range receiver.connections[conn].ConnectedSockets {
		socket.Close()
		fmt.Println("Deleting child ", receiver.connections[socket].ID)
		delete(receiver.connections[conn].ConnectedSockets, socket)
	}
	// conn.Close()
	delete(receiver.connections, conn)
	fmt.Println("Finish Deleting")
	check := receiver.connections[conn]
	fmt.Println("Check", check)
}

func (receiver *Type) Insert(conn *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if receiver.connections == nil {
		receiver.connections = map[*websocket.Conn]MySocket{}
	}

	receiver.connections[conn] = MySocket{
		Socket:           conn,
		ID:               ConnectedSocketsIndex,
		IsBroadcaster:    false,
		ConnectedSockets: make(map[*websocket.Conn]MySocket),
	}
	ConnectedSocketsIndex++
}

func (receiver *Type) InsertConnected(conn *websocket.Conn, connectedConn *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	mySocket := receiver.connections[conn]
	connectedSockets := mySocket.ConnectedSockets

	_, ok := connectedSockets[connectedConn]
	if !ok {
		if connectedSockets == nil {
			connectedSockets = make(map[*websocket.Conn]MySocket)
		}
		connectedSockets[connectedConn] = receiver.connections[connectedConn]
		mySocket.ConnectedSockets = connectedSockets
		receiver.connections[conn] = mySocket
	}
}

func (receiver *Type) DeleteConnected(conn *websocket.Conn, connectedConn *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	mySocket := receiver.connections[conn]
	connectedSockets := mySocket.ConnectedSockets

	_, ok := connectedSockets[connectedConn]
	if ok {
		delete(connectedSockets, connectedConn)
		mySocket.ConnectedSockets = connectedSockets
		receiver.connections[conn] = mySocket
	}
}

func (receiver *Type) RemoveBroadcasters() {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for conn := range receiver.connections {
		if receiver.connections[conn].IsBroadcaster {
			mySocket := receiver.connections[conn]
			mySocket.IsBroadcaster = false
			receiver.connections[conn] = mySocket
		}
	}
}

func (receiver *Type) SetBroadcaster(conn *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	mySocket := receiver.connections[conn]
	mySocket.IsBroadcaster = true
	receiver.connections[conn] = mySocket
}

func (receiver *Type) RemoveBroadcaster(conn *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	mySocket := receiver.connections[conn]
	mySocket.IsBroadcaster = false
	receiver.connections[conn] = mySocket
}

func (receiver *Type) GetBroadcaster() (MySocket, bool) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for conn := range receiver.connections {
		if receiver.connections[conn].IsBroadcaster {
			return receiver.connections[conn], true
		}
	}

	return MySocket{}, false
}

func (receiver *Type) GetSocketByID(ID uint64) (MySocket, bool) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for conn := range receiver.connections {
		if receiver.connections[conn].ID == ID {
			return receiver.connections[conn], true
		}
	}

	return MySocket{}, false
}

func (receiver *Type) Reset() {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for socket := range receiver.connections {
		for childSocket := range receiver.connections[socket].ConnectedSockets {
			childSocket.Close()
			delete(receiver.connections[socket].ConnectedSockets, childSocket)
		}
		socket.Close()
		delete(receiver.connections, socket)
	}
	ConnectedSocketsIndex = 0
}

func (receiver *Type) GetConnections() map[*websocket.Conn]MySocket {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	return receiver.connections
}

func (receiver *Type) SetName(conn *websocket.Conn, name string) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	socket := receiver.connections[conn]
	socket.Name = name
	receiver.connections[conn] = socket
}

func (receiver *Type) GetParent(conn *websocket.Conn) MySocket {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for socket := range receiver.connections {
		for child := range receiver.connections[socket].ConnectedSockets {
			if child == conn {
				return receiver.connections[socket]
			}
		}
	}

	return MySocket{}
}
