package websocketmap

import (
	"sync"

	"github.com/gorilla/websocket"
)

type MySocket struct {
	Socket           *websocket.Conn
	ID               uint64
	IsBroadcaster    bool
	ConnectedSockets map[*websocket.Conn]MySocket
}

var ConnectedSocketsIndex uint64 = 0

type Type struct {
	mutex       sync.Mutex
	connections map[*websocket.Conn]MySocket
}

func (receiver *Type) Get(conn *websocket.Conn) MySocket {
	return receiver.connections[conn]
}

func (receiver *Type) Delete(conn *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for socket := range receiver.connections {
		delete(receiver.connections[socket].ConnectedSockets, conn)
	}
	delete(receiver.connections, conn)
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
