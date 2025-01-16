package rooms_test

import (
	"testing"

	"github.com/reiver/logjam/lib/rooms"
)

func TestRoomURL(t *testing.T) {

	tests := []struct{
		RoomID   string
		BasePath string
		Host     string
		Expected string
	}{
		{

		},



		{
			RoomID: "apple",
		},
		{
			RoomID: "Banana",
		},
		{
			RoomID: "CHERRY",
		},



		{
			BasePath: "/",
		},
		{
			BasePath: "/apple",
		},
		{
			BasePath: "/apple/",
		},
		{
			BasePath: "/apple/banana",
		},
		{
			BasePath: "/apple/banana/",
		},
		{
			BasePath: "/apple/banana/cherry",
		},
		{
			BasePath: "/apple/banana/cherry/",
		},



		{
			Host: "example.com",
		},



		{
			RoomID: "ONE",
			BasePath: "/",
			Host: "example.com",
			Expected: "https://example.com/ONE",
		},
		{
			RoomID: "ONE",
			BasePath: "/apple",
			Host: "example.com",
			Expected: "https://example.com/apple/ONE",
		},
		{
			RoomID: "ONE",
			BasePath: "/apple/",
			Host: "example.com",
			Expected: "https://example.com/apple/ONE",
		},
		{
			RoomID: "ONE",
			BasePath: "/apple/Banana",
			Host: "example.com",
			Expected: "https://example.com/apple/Banana/ONE",
		},
		{
			RoomID: "ONE",
			BasePath: "/apple/Banana/",
			Host: "example.com",
			Expected: "https://example.com/apple/Banana/ONE",
		},
		{
			RoomID: "ONE",
			BasePath: "/apple/Banana/CHERRY",
			Host: "example.com",
			Expected: "https://example.com/apple/Banana/CHERRY/ONE",
		},
		{
			RoomID: "ONE",
			BasePath: "/apple/Banana/CHERRY/",
			Host: "example.com",
			Expected: "https://example.com/apple/Banana/CHERRY/ONE",
		},



		{
			RoomID: "ONE",
			BasePath: "/",
			Host: "EXAMPLE.COM",
			Expected: "https://example.com/ONE",
		},
		{
			RoomID: "ONE",
			BasePath: "/apple",
			Host: "EXAMPLE.COM",
			Expected: "https://example.com/apple/ONE",
		},
		{
			RoomID: "ONE",
			BasePath: "/apple/",
			Host: "EXAMPLE.COM",
			Expected: "https://example.com/apple/ONE",
		},
		{
			RoomID: "ONE",
			BasePath: "/apple/Banana",
			Host: "EXAMPLE.COM",
			Expected: "https://example.com/apple/Banana/ONE",
		},
		{
			RoomID: "ONE",
			BasePath: "/apple/Banana/",
			Host: "EXAMPLE.COM",
			Expected: "https://example.com/apple/Banana/ONE",
		},
		{
			RoomID: "ONE",
			BasePath: "/apple/Banana/CHERRY",
			Host: "EXAMPLE.COM",
			Expected: "https://example.com/apple/Banana/CHERRY/ONE",
		},
		{
			RoomID: "ONE",
			BasePath: "/apple/Banana/CHERRY/",
			Host: "EXAMPLE.COM",
			Expected: "https://example.com/apple/Banana/CHERRY/ONE",
		},
	}

	for testNumber, test := range tests {

		actual := rooms.RoomURL(test.RoomID, test.BasePath, test.Host)

		expected := test.Expected

		if expected != actual {
			t.Errorf("For test #%d, the actual room-url is not what was expected." , testNumber)
			t.Logf("EXPECTED: %q", expected)
			t.Logf("ACTUAL:   %q", actual)
			t.Logf("HOST:         %q", test.Host)
			t.Logf("BASE-PATH:    %q", test.BasePath)
			t.Logf("ROOM-ID:      %q", test.RoomID)
			continue
		}
	}
}
