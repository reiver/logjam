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
	Connections map[*websocket.Conn]MySocket
}

func (receiver *Type) Delete(conn *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	delete(receiver.Connections, conn)
}

func (receiver *Type) Insert(conn *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.Connections[conn] = MySocket{
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

	mySocket := receiver.Connections[conn]
	connectedSockets := mySocket.ConnectedSockets

	_, ok := connectedSockets[connectedConn]
	if !ok {
		connectedSockets[connectedConn] = receiver.Connections[connectedConn]
		mySocket.ConnectedSockets = connectedSockets
		receiver.Connections[conn] = mySocket
	}
}

func (receiver *Type) DeleteConnected(conn *websocket.Conn, connectedConn *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	mySocket := receiver.Connections[conn]
	connectedSockets := mySocket.ConnectedSockets

	_, ok := connectedSockets[connectedConn]
	if ok {
		delete(connectedSockets, connectedConn)
		mySocket.ConnectedSockets = connectedSockets
		receiver.Connections[conn] = mySocket
	}
}

func (receiver *Type) RemoveBroadcasters() {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for conn := range receiver.Connections {
		if receiver.Connections[conn].IsBroadcaster {
			mySocket := receiver.Connections[conn]
			mySocket.IsBroadcaster = false
			receiver.Connections[conn] = mySocket
		}
	}
}

func (receiver *Type) SetBroadcaster(conn *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	mySocket := receiver.Connections[conn]
	mySocket.IsBroadcaster = true
	receiver.Connections[conn] = mySocket
}

func (receiver *Type) RemoveBroadcaster(conn *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	mySocket := receiver.Connections[conn]
	mySocket.IsBroadcaster = false
	receiver.Connections[conn] = mySocket
}

func (receiver *Type) GetBroadcaster() (MySocket, bool) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for conn := range receiver.Connections {
		if receiver.Connections[conn].IsBroadcaster {
			return receiver.Connections[conn], true
		}
	}

	return MySocket{}, false
}

func (receiver *Type) GetSocketByID(ID uint64) (MySocket, bool) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for conn := range receiver.Connections {
		if receiver.Connections[conn].ID == ID {
			return receiver.Connections[conn], true
		}
	}

	return MySocket{}, false
}
