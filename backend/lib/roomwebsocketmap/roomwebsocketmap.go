package roomwebsocketmap

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type MySocket struct {
	Socket    *websocket.Conn
	ID        uint64
	Name      string
	HasStream bool
}

type Room struct {
	Description string
	Sockets     map[*websocket.Conn]MySocket
	Name        string
}

var ConnectedSocketsIndex uint64 = 0

type Type struct {
	mutex sync.Mutex
	rooms map[string]Room
}

var Map Type = Type{}

func (receiver *Type) GetRoom(roomName string) Room {
	return receiver.rooms[roomName]
}

func (receiver *Type) GetSocketRoom(socket *websocket.Conn) (Room, error) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for room := range receiver.rooms {
		for asocket := range receiver.rooms[room].Sockets {
			if socket == asocket {
				return receiver.rooms[room], nil
			}
		}
	}
	return Room{}, errors.New("socket not found")
}

func (receiver *Type) GetSocket(socket *websocket.Conn) (Room, MySocket, error) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for room := range receiver.rooms {
		for asocket := range receiver.rooms[room].Sockets {
			if socket == asocket {
				return receiver.rooms[room], receiver.rooms[room].Sockets[asocket], nil
			}
		}
	}
	return Room{}, MySocket{}, errors.New("socket not found")
}

func (receiver *Type) AddToRoom(roomName string, socket MySocket) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if receiver.rooms == nil {
		receiver.rooms = map[string]Room{}
	}
	room := receiver.rooms[roomName]
	if room.Sockets == nil {
		room.Sockets = map[*websocket.Conn]MySocket{}
	}
	socket.ID = ConnectedSocketsIndex
	ConnectedSocketsIndex++
	room.Sockets[socket.Socket] = socket
	receiver.rooms[roomName] = room
}

func (receiver *Type) RemoveFromRoom(rooName string, socket MySocket) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	room := receiver.rooms[rooName]
	delete(room.Sockets, socket.Socket)
	receiver.rooms[rooName] = room
}

func (receiver *Type) GetSocketByID(id uint64) (Room, MySocket, error) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	for room := range receiver.rooms {
		for asocket := range receiver.rooms[room].Sockets {
			if id == receiver.rooms[room].Sockets[asocket].ID {
				return receiver.rooms[room], receiver.rooms[room].Sockets[asocket], nil
			}
		}
	}
	return Room{}, MySocket{}, errors.New("socket not found")
}
