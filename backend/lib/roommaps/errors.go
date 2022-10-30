package roommaps

import (
	"github.com/reiver/go-fck"
)

const (
	errNilReceiver  = fck.Error("roommaps: nil receiver")
	errRoomNotFound = fck.Error("roommaps: room not found")
)
