package websock

import (
	"sync"

	"github.com/gorilla/websocket"
)

type SocketKeeper struct {
	mutex sync.Mutex
	wsConn *websocket.Conn
	id     uint64
}

func newSocketKeeper(wsConn *websocket.Conn, id uint64) *SocketKeeper {
	if nil == wsConn {
		return nil
	}

	return &SocketKeeper{
		wsConn:wsConn,
		id:id,
	}
}

func (receiver *SocketKeeper) ID() uint64 {
	if nil == receiver {
		var nada uint64
		return nada
	}

	return receiver.id
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
