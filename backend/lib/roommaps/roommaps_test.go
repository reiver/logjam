package roommaps

import (
	"testing"

	"github.com/mmcomp/go-binarytree"
)

var theRoomMaps Type = Type{
	roomMaps: map[string]*RoomType{},
}
var emptyRoomMaps Type

// TestGetEmpty calls roommaps.Get from empty rooms
func TestGetEmpty(t *testing.T) {
	_, ok := emptyRoomMaps.Get("no-room")
	if ok {
		t.Error(`It should return FALSE`)
	}
}

// TestGetNotFound calls roommaps.Get from no rooms
func TestGetNotFound(t *testing.T) {
	_, ok := theRoomMaps.Get("no-room")
	if ok {
		t.Error(`It should return FALSE`)
	}
}

// TestSet calls roommaps.Set add a new room as room-one
func TestSet(t *testing.T) {
	beforeLength := len(theRoomMaps.roomMaps)
	theRoomMaps.Set(" room-one", &binarytree.Tree{})
	afterLength := len(theRoomMaps.roomMaps)
	if beforeLength != 0 || afterLength != 1 {
		t.Error(`'beforeLength' should be 0 but `, beforeLength)
		t.Error(`'afterLength' should be 1 but `, afterLength)
	}
}

// TestSetWithGet calls roommaps.Set add a new room as room-one
func TestSetWithGet(t *testing.T) {
	beforeLength := len(theRoomMaps.roomMaps)
	theRoomMaps.Set("room-two", &binarytree.Tree{})
	afterLength := len(theRoomMaps.roomMaps)
	if afterLength-beforeLength != 1 {
		t.Error(`'afterLength' should be 1 bigger than 'beforeLength' `, afterLength, beforeLength)
	}
	_, ok := theRoomMaps.Get("room-two")
	if !ok {
		t.Error(`'ok' should be TRUE' but is `, ok)
	}
}

// TestSetMetaData calls roommaps.Set add a new room meta data
func TestSetMetaData(t *testing.T) {
	// err := theRoomMaps.SetMetData("room-two", metadata.MetaData{
	// 	BackgroundURL: `https://google.com`,
	// })
	err := theRoomMaps.SetMetData("room-two", map[string]string{
		"BackgroundURL": `https://google.com`,
	})
	if err != nil {
		t.Error(`'err' should be nil but is `, err)
	}
	room, ok := theRoomMaps.Get("room-two")
	if !ok {
		t.Error(`'ok' should be TRUE' but is `, ok)
	}
	if room.MetaData["BackgroundURL"] != `https://google.com` {
		t.Error(`'room.MetaData.BackgroundURL' should be 'https://google.com' but is `, room.MetaData["BackgroundURL"])
	}
}
