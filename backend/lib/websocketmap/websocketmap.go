package websocketmap

import (
	"sync"

	"github.com/gorilla/websocket"
)

type MySocket struct {
	Socket        *websocket.Conn
	ID            uint64
	IsBroadcaster bool
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
		Socket:        conn,
		ID:            ConnectedSocketsIndex,
		IsBroadcaster: false,
	}
	ConnectedSocketsIndex++
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

func (receiver *WebSocketMapType) GetBroadcaster() (*websocket.Conn, bool) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for conn := range receiver.Connections {
		if receiver.Connections[conn].IsBroadcaster {
			return conn, true
		}
	}

	return nil, false
}

func (receiver *WebSocketMapType) GetSocketByID(ID uint64) (*websocket.Conn, bool) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for conn := range receiver.Connections {
		if receiver.Connections[conn].ID == ID {
			return conn, true
		}
	}

	return nil, false
}
