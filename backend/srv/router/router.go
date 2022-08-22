package router

import (
	"net/http"
	"strings"

	Node "github.com/sparkscience/logjam/backend/srv/node"
)

type Router struct {
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/rooms/") {
		roomName := strings.Split(strings.Replace(req.URL.Path, "/rooms/", "", 1), "/")[0]
		r := Node.NewRoom(roomName)
		r.ServeHTTP(w, req)
	} else {
		http.FileServer(http.Dir("./public")).ServeHTTP(w, req)
	}
}

func NewRouter() *Router {
	return &Router{}
}
