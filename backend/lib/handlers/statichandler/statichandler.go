package statichandler

import "net/http"

// type httpHandler struct {
// }

// var Handler httpHandler

// func (receiver httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	w.Write([]byte("Hello!"))
// }
var Handler = http.FileServer(http.Dir("/public/statics/files"))
