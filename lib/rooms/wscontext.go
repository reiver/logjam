package rooms

import (
	"github.com/reiver/logjam/lib/msgs"
)

type WSContext struct {
	SocketID      uint64
	PureMessage   []byte
	ParsedMessage *msgs.MessageContract
	RoomId        string
}
