package services

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"sourcecode.social/greatape/logjam/models/contracts"
	"sync"
)

type SocketKeeper struct {
	*sync.Mutex
	wsConn *websocket.Conn
	ID     uint64
}

func (s *SocketKeeper) WriteMessage(data []byte) error {
	s.Lock()
	defer s.Unlock()
	err := s.wsConn.WriteMessage(websocket.TextMessage, data)
	return err
}

type socketService struct {
	*sync.Mutex
	logger      contracts.ILogger
	lastId      uint64
	sockets     map[*websocket.Conn]*SocketKeeper
	socketsById map[uint64]*websocket.Conn
}

func NewSocketService(logger contracts.ILogger) contracts.ISocketService {
	return &socketService{
		Mutex:       &sync.Mutex{},
		logger:      logger,
		lastId:      0,
		sockets:     make(map[*websocket.Conn]*SocketKeeper),
		socketsById: make(map[uint64]*websocket.Conn),
	}
}

func (s *socketService) GetSocketId(conn *websocket.Conn) (*uint64, error) {
	if socket, exists := s.sockets[conn]; exists {
		return &socket.ID, nil
	}
	return nil, nil
}

func (s *socketService) Send(data interface{}, receiverIds ...uint64) error {
	var jsonData []byte
	var err error
	if b, isByte := data.([]byte); isByte {
		jsonData = b
	} else {
		jsonData, err = json.Marshal(data)
		if err != nil {
			println(err.Error())
			return err
		}
	}
	s.Lock()
	defer s.Unlock()
	for _, id := range receiverIds {
		if socket, exists := s.socketsById[id]; exists {
			keeper := s.sockets[socket]
			_ = keeper.WriteMessage(jsonData)
		}
	}
	return nil
}

func (s *socketService) OnConnect(conn *websocket.Conn) (uint64, error) {
	s.logger.Log("socket_svc", contracts.LDebug, "new socket connected")
	s.Lock()
	defer s.Unlock()

	id := s.lastId
	s.lastId++
	s.sockets[conn] = &SocketKeeper{
		Mutex:  &sync.Mutex{},
		wsConn: conn,
		ID:     id,
	}
	s.socketsById[id] = conn

	return id, nil
}

func (s *socketService) OnDisconnect(conn *websocket.Conn) error {
	s.logger.Log("socket_svc", contracts.LDebug, "a socket got disconnected")
	s.Lock()
	defer s.Unlock()
	if keeper, exists := s.sockets[conn]; exists {
		delete(s.socketsById, keeper.ID)
		delete(s.sockets, conn)
	}

	return nil
}
