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

type WebSocketMapType struct {
	mutex       sync.Mutex
	Connections map[*websocket.Conn]MySocket
}

func (receiver *WebSocketMapType) Delete(conn *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	delete(receiver.Connections, conn)
}

func (receiver *WebSocketMapType) Insert(conn *websocket.Conn) {
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

func (receiver *WebSocketMapType) InsertConnected(conn *websocket.Conn, connectedConn *websocket.Conn) {
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

func (receiver *WebSocketMapType) DeleteConnected(conn *websocket.Conn, connectedConn *websocket.Conn) {
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

func (receiver *WebSocketMapType) RemoveBroadcasters() {
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

func (receiver *WebSocketMapType) SetBroadcaster(conn *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	mySocket := receiver.Connections[conn]
	mySocket.IsBroadcaster = true
	receiver.Connections[conn] = mySocket
}

func (receiver *WebSocketMapType) RemoveBroadcaster(conn *websocket.Conn) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	mySocket := receiver.Connections[conn]
	mySocket.IsBroadcaster = false
	receiver.Connections[conn] = mySocket
}

func (receiver *WebSocketMapType) GetBroadcaster() (MySocket, bool) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for conn := range receiver.Connections {
		if receiver.Connections[conn].IsBroadcaster {
			return receiver.Connections[conn], true
		}
	}

	return MySocket{}, false
}

func (receiver *WebSocketMapType) GetSocketByID(ID uint64) (MySocket, bool) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for conn := range receiver.Connections {
		if receiver.Connections[conn].ID == ID {
			return receiver.Connections[conn], true
		}
	}

	return MySocket{}, false
}
