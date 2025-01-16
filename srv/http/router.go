package httpsrv

import (
	"net/http"

	"github.com/gorilla/mux"
)

var Router *mux.Router = mux.NewRouter()

func init() {
	Router = mux.NewRouter()
	if nil == Router {
		panic("nil HTTP gorilla mux router")
	}

	Router.Use(mux.CORSMethodMiddleware(Router))

	Router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			remoteAddr := request.RemoteAddr
			log.Debugf(`HTTP-request method=%q remote-addr=%q request-uri=%q`, request.Method, remoteAddr, request.URL.RequestURI())
			handler.ServeHTTP(writer, request)
		})
	})
}
