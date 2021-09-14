package statichandler

import (
	"net/http"
)

// var Handler http.Handler = http.FileServer(http.Dir("./public/statics/files"))
func Handler(path string) http.Handler {
	return http.FileServer(http.Dir(path))
}
