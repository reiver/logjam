package rooms_test

import (
	"fmt"

	"github.com/reiver/logjam/lib/rooms"
)

func ExampleRoomURL() {

	const roomID   string = "abc-123"
	const basePath string = "/api/v1/rooms"
	const host     string = "example.com"

	var roomURL string = rooms.RoomURL(roomID, basePath, host)

	fmt.Printf("room-url = %s\n", roomURL)

	// Output:
	// room-url = https://example.com/api/v1/rooms/abc-123
}
