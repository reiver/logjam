package websocketmap

import (
	"sync"

	"github.com/gorilla/websocket"
)

type mySocket struct {
    Socket        *websocket.Conn
    ID            uint64
    IsBroadcaster bool
}

var ConnectedSocketsIndex uint64 = 0

type WebSocketMapType struct {
    mutex sync.Mutex
    connections map[*websocket.Conn]mySocket
}

func (receiver *WebSocketMapType) Delete(conn *websocket.Conn) {
    receiver.mutex.Lock()
    defer receiver.mutex.Unlock()

    delete(receiver.connections, conn)
}

func (receiver *WebSocketMapType) Insert(conn *websocket.Conn) {
    receiver.mutex.Lock()
    defer receiver.mutex.Unlock()

    receiver.connections[conn]  = mySocket{
		Socket:        conn,
		ID:            ConnectedSocketsIndex,
		IsBroadcaster: false,
	}
	ConnectedSocketsIndex++
}

func CreateWebsocketMap() WebSocketMapType {
	websocketmap := WebSocketMapType{
		connections: make(map[*websocket.Conn]mySocket),
	};

	return websocketmap
}