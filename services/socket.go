package services

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/reiver/logjam/lib/logs"
	"github.com/reiver/logjam/models/contracts"
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

type socketService struct {
	*sync.Mutex
	logger      logs.Logger
	lastId      uint64
	sockets     map[*websocket.Conn]*SocketKeeper
	socketsById map[uint64]*websocket.Conn
	pingTimeout time.Duration
}

func NewSocketService(logger logs.Logger) contracts.ISocketService {
	return &socketService{
		Mutex:       &sync.Mutex{},
		logger:      logger,
		lastId:      0,
		sockets:     make(map[*websocket.Conn]*SocketKeeper),
		socketsById: make(map[uint64]*websocket.Conn),
		pingTimeout: 5 * time.Second,
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
			_ = keeper.WriteTextMessage(jsonData)
		}
	}
	return nil
}

func (s *socketService) OnConnect(conn *websocket.Conn) (uint64, error) {
	s.logger.Debug("socket_svc", "new socket connected")
	s.Lock()
	defer s.Unlock()

	id := s.getNewId()
	s.sockets[conn] = &SocketKeeper{
		Mutex:  &sync.Mutex{},
		wsConn: conn,
		ID:     id,
	}
	s.socketsById[id] = conn
	conn.SetCloseHandler(func(code int, text string) error {
		_ = s.OnDisconnect(conn, code, text)
		return nil
	})
	lastPong := time.Now()
	conn.SetPongHandler(func(appData string) error {
		lastPong = time.Now()
		return nil
	})
	go func() {
		for {
			s.Lock()
			if keeper, exists := s.sockets[conn]; exists {
				s.Unlock()
				err := keeper.WriteMessage(websocket.PingMessage, []byte("keepalive"))
				if err != nil {
					return
				}
				time.Sleep(s.pingTimeout / 2)
				if time.Since(lastPong) > s.pingTimeout {
					_ = conn.Close()
					break
				}
			} else {
				s.Unlock()
				break
			}
		}
	}()
	return id, nil
}

func (s *socketService) getNewId() uint64 {
	id := s.lastId
	s.lastId++
	return id
}
func (s *socketService) GetNewID() uint64 {
	s.Lock()
	defer s.Unlock()
	return s.getNewId()
}

func (s *socketService) OnDisconnect(conn *websocket.Conn, code int, error string) error {
	s.Lock()
	defer s.Unlock()
	if keeper, exists := s.sockets[conn]; exists {
		s.logger.Debug("socket_svc", "a socket got disconnected ["+strconv.FormatUint(keeper.ID, 10)+"]", strconv.Itoa(code), ":", error)
		delete(s.socketsById, keeper.ID)
		delete(s.sockets, conn)
	}

	return nil
}

func (s *socketService) Disconnect(socketId uint64) error {
	if conn, exists := s.socketsById[socketId]; exists {
		return conn.Close()
	}
	return nil
}
