package models

type WSContext struct {
	SocketID      uint64
	PureMessage   []byte
	ParsedMessage *MessageContract
	RoomId        string
}

type MessageContract struct {
	Type   string
	Data   string
	Target string
	Name   string
}
