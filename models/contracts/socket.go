package contracts

import "github.com/gorilla/websocket"

type ISocketService interface {
	Send(data interface{}, receiverIds ...uint64) error

	GetNewID() uint64
	OnConnect(conn *websocket.Conn) (uint64, error)
	OnDisconnect(conn *websocket.Conn, code int, error string) error
	GetSocketId(conn *websocket.Conn) (*uint64, error)
	Disconnect(socketId uint64) error
}
