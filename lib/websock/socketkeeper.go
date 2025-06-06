package websock

import (
	"sync"

	"github.com/gorilla/websocket"
)

type SocketKeeper struct {
	*sync.Mutex
	wsConn *websocket.Conn
	ID     uint64
}

func (s *SocketKeeper) WriteTextMessage(data []byte) error {
	s.Lock()
	defer s.Unlock()
	err := s.wsConn.WriteMessage(websocket.TextMessage, data)
	return err
}

func (s *SocketKeeper) WriteMessage(messageType int, data []byte) error {
	s.Lock()
	defer s.Unlock()
	err := s.wsConn.WriteMessage(messageType, data)
	return err
}
