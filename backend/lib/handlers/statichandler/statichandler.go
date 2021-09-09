package statichandler

import (
	"net/http"
)

func Handler(path string) http.Handler {
	return http.FileServer(http.Dir(path))
}
