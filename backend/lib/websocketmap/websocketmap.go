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

// func CreateWebsocketMap() WebSocketMapType {
// 	websocketmap := WebSocketMapType{
// 		Connections: make(map[*websocket.Conn]MySocket),
// 	}

// 	return websocketmap
// }
