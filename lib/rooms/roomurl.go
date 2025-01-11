package rooms

import (
	"net/url"
	"strings"

	libpath "github.com/reiver/go-path"
)

// RoomURL is used to construct the URL that represents the canonical URL for a room.
func RoomURL(roomID string, basePath string, host string) string {

	if "" == roomID {
		return ""
	}
	if "" == basePath {
		return ""
	}
	if "" == host {
		return ""
	}

	var path string = libpath.Join(basePath, roomID)
	if "." == path {
		path = "/"
	}


	var urloc url.URL

	urloc.Scheme = "https"
	urloc.Host   = strings.ToLower(host)
	urloc.Path   = path

	return urloc.String()
}
