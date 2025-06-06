package websock

import (
	"sync"

	"github.com/gorilla/websocket"
)

type SocketKeeper struct {
//@TODO: why is this a pointer to a mutex, rather than just a mutex.
	*sync.Mutex
	wsConn *websocket.Conn
	ID     uint64
}

func (receiver *SocketKeeper) WriteTextMessage(data []byte) error {
	if nil == receiver {
		return errNilReceiver
	}

	receiver.Lock()
	defer receiver.Unlock()

	err := receiver.wsConn.WriteMessage(websocket.TextMessage, data)
	return err
}

func (receiver *SocketKeeper) WriteMessage(messageType int, data []byte) error {
	if nil == receiver {
		return errNilReceiver
	}

	receiver.Lock()
	defer receiver.Unlock()

	err := receiver.wsConn.WriteMessage(messageType, data)
	return err
}
