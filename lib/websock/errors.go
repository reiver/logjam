package websock

import (
	"github.com/reiver/go-erorr"
)

const (
	errNilReceiver            = erorr.Error("websock: nil receiver")
	errNilWebSocketConnection = erorr.Error("websock: nil websocket connection")
)
