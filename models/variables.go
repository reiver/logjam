package models

import "sync"

var (
	goldGorillaId uint64 = 18446744073709551615
	ggidLock      *sync.Mutex
)

func init() {
	ggidLock = &sync.Mutex{}
}

func GetGoldGorillaId() uint64 {
	ggidLock.Lock()
	defer ggidLock.Unlock()
	return goldGorillaId
}
func IncreaseGoldGorillaId() {
	ggidLock.Lock()
	defer ggidLock.Unlock()
	goldGorillaId++
}
func DecreaseGoldGorillaId() {
	ggidLock.Lock()
	defer ggidLock.Unlock()
	goldGorillaId--
}

const (
	RoomMessagesMetaDataKey string = "messages"
)
