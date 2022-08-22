package node

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Room struct {
	name string
	// join is a channel for clients wishing to join the room.
	join chan *Node
	// leave is a channel for clients wishing to leave the room.
	leave chan *Node
	// clients holds all current clients in this room.
	clients map[*Node]bool
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.join:
			fmt.Println("run join ", r.name, client.id)
			// joining
			r.clients[client] = true
		case client := <-r.leave:
			fmt.Println("run leave ", r.name, client.id)
			// leaving
			delete(r.clients, client)
			// close(client.send)
			// case msg := <-r.forward:
			// 	// forward message to all clients
			// 	for client := range r.clients {
			// 		client.send <- msg
			// 	}
		}
	}
}

var rooms = make(map[string]*Room)

func NewRoom(name string) *Room {
	if _, ok := rooms[name]; ok {
		return rooms[name]
	}

	room := &Room{
		name:    name,
		join:    make(chan *Node),
		leave:   make(chan *Node),
		clients: make(map[*Node]bool),
	}
	rooms[name] = room
	go room.run()
	return room
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize}

var socketId = 0

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("room.ServeHTTP ", r.name, len(r.clients))
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		fmt.Println("ServeHTTP error ", err)
		return
	}
	node := NewNode(uint64(socketId), socket.RemoteAddr().String(), false, false, false, socket)
	go node.Read()
	fmt.Println("node ", node.id, " joined ", r.name)
	socketId++
	fmt.Println("OK1")
	r.join <- node
	fmt.Println("OK2")
	// defer func() { r.leave <- node }()
}
