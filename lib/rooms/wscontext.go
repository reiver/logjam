package rooms

import (
	"codeberg.org/greatape/logjam/lib/msgs"
)

type WSContext struct {
	SocketID      uint64
	PureMessage   []byte
	ParsedMessage *msgs.Message
	RoomId        string
}
