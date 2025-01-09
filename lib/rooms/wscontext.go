package rooms

import (
	"github.com/reiver/logjam/models"
)

type WSContext struct {
	SocketID      uint64
	PureMessage   []byte
	ParsedMessage *models.MessageContract
	RoomId        string
}
