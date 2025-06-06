package websock

import (
	"sync"

	"github.com/gorilla/websocket"
)

type SocketKeeper struct {
	mutex sync.Mutex
	wsConn *websocket.Conn
	ID     uint64
}

func (receiver *SocketKeeper) WriteTextMessage(data []byte) error {
	if nil == receiver {
		return errNilReceiver
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	err := receiver.wsConn.WriteMessage(websocket.TextMessage, data)
	return err
}

func (receiver *SocketKeeper) WriteMessage(messageType int, data []byte) error {
	if nil == receiver {
		return errNilReceiver
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	err := receiver.wsConn.WriteMessage(messageType, data)
	return err
}
