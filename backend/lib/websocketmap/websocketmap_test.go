package websocketmap

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestWebSocketMap_Insert(t *testing.T) {
	var websocketmaps Type = Type{}

	tests := []struct {
		Input *websocket.Conn

		ExpectedLen int
	}{
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 1,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 2,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 3,
		},
	}

	for testNumber, test := range tests {
		websocketmaps.Insert(test.Input)
		if len(websocketmaps.connections) != test.ExpectedLen {
			t.Errorf("Test %d :  %d was expected but got %d", testNumber, test.ExpectedLen, len(websocketmaps.connections))
		}
	}
}

func TestWebSocketMap_Get(t *testing.T) {
	var websocketmaps Type = Type{}

	tests := []struct {
		Input *websocket.Conn
	}{
		{
			Input: &websocket.Conn{},
		},
		{
			Input: &websocket.Conn{},
		},
		{
			Input: &websocket.Conn{},
		},
	}
	for _, test := range tests {
		websocketmaps.Insert(test.Input)
	}

	for testNumber, test := range tests {
		mySocket := websocketmaps.Get(test.Input)
		if mySocket.Socket != test.Input {
			t.Errorf("Test %d : A socket was expected but got another one", testNumber)
		}
	}
}

func TestWebSocketMap_Delete(t *testing.T) {
	var websocketmaps Type = Type{}

	tests := []struct {
		Input *websocket.Conn
	}{
		{
			Input: &websocket.Conn{},
		},
		{
			Input: &websocket.Conn{},
		},
		{
			Input: &websocket.Conn{},
		},
	}
	for _, test := range tests {
		websocketmaps.Insert(test.Input)
	}

	for testNumber, test := range tests {
		websocketmaps.Delete(test.Input)
		deleteSocket := websocketmaps.Get(test.Input)
		if deleteSocket.Socket != nil {
			t.Errorf("Test %d : A nil was expected but got a socket", testNumber)
		}
	}
}

func TestWebSocketMap_InsertConnected(t *testing.T) {
	var websocketmaps Type = Type{}
	parent := &websocket.Conn{}
	websocketmaps.Insert(parent)
	tests := []struct {
		Input *websocket.Conn

		ExpectedLen int
	}{
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 1,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 2,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 3,
		},
	}

	for testNumber, test := range tests {
		websocketmaps.InsertConnected(parent, test.Input)
		parentSocket := websocketmaps.Get(parent)
		if len(parentSocket.ConnectedSockets) != test.ExpectedLen {
			t.Errorf("Test %d :  %d was expected but got %d", testNumber, test.ExpectedLen, len(parentSocket.ConnectedSockets))
		}
	}
}

func TestWebSocketMap_DeleteConnected(t *testing.T) {
	var websocketmaps Type = Type{}
	parent := &websocket.Conn{}
	websocketmaps.Insert(parent)
	tests := []struct {
		Input *websocket.Conn

		ExpectedLen int
	}{
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 0,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 0,
		},
		{
			Input:       &websocket.Conn{},
			ExpectedLen: 0,
		},
	}

	for testNumber, test := range tests {
		websocketmaps.InsertConnected(parent, test.Input)
		parentSocket := websocketmaps.Get(parent)
		if len(parentSocket.ConnectedSockets) != test.ExpectedLen+1 {
			t.Errorf("Test %d :  %d was expected but got %d", testNumber, test.ExpectedLen+1, len(parentSocket.ConnectedSockets))
		}
		websocketmaps.DeleteConnected(parent, test.Input)
		parentSocket = websocketmaps.Get(parent)
		if len(parentSocket.ConnectedSockets) != test.ExpectedLen {
			t.Errorf("Test %d :  %d was expected but got %d", testNumber, test.ExpectedLen, len(parentSocket.ConnectedSockets))
		}
	}
}

func TestWebSocketMap_DeleteWithConnectedSockets(t *testing.T) {
	var websocketmaps Type = Type{}

	sock1 := &websocket.Conn{}
	sock2 := &websocket.Conn{}
	sock3 := &websocket.Conn{}
	tests := []struct {
		Input            *websocket.Conn
		ConnectedSockets map[*websocket.Conn]MySocket
	}{
		{
			Input:            sock1,
			ConnectedSockets: make(map[*websocket.Conn]MySocket),
		},
		{
			Input:            sock2,
			ConnectedSockets: make(map[*websocket.Conn]MySocket),
		},
		{
			Input:            sock3,
			ConnectedSockets: make(map[*websocket.Conn]MySocket),
		},
	}
	for _, test := range tests {
		websocketmaps.Insert(test.Input)
	}
	websocketmaps.InsertConnected(sock1, sock2)
	websocketmaps.InsertConnected(sock2, sock3)
	websocketmaps.Delete(sock3)
	socket2 := websocketmaps.Get(sock2)
	if len(socket2.ConnectedSockets) != 0 {
		t.Errorf("Socket 2 must have no connectedSockets but have %d sockets", len(socket2.ConnectedSockets))
	}
	websocketmaps.Delete(sock2)
	socket1 := websocketmaps.Get(sock1)
	if len(socket1.ConnectedSockets) != 0 {
		t.Errorf("Socket 1 must have no connectedSockets but have %d sockets", len(socket1.ConnectedSockets))
	}
}
