package contracts

import "github.com/gorilla/websocket"

type ISocketService interface {
	Send(data interface{}, receiverIds ...uint64) error

	OnConnect(conn *websocket.Conn) (uint64, error)
	OnDisconnect(conn *websocket.Conn) error
	GetSocketId(conn *websocket.Conn) (*uint64, error)
}
