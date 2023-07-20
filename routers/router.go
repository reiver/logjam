package routers

import (
	"github.com/gorilla/mux"
	"github.com/sparkscience/logjam/controllers"
	"github.com/sparkscience/logjam/models/contracts"
	"net/http"
)

type IRouteRegistrar interface {
	registerRoutes(router *mux.Router)
}

type Router struct {
	router              *mux.Router
	roomWSRouter        IRouteRegistrar
	auxiliaryNodeRouter IRouteRegistrar
}

func NewRouter(roomWSCtrl *controllers.RoomWSController, auxiliaryNodeCtrl *controllers.AuxiliaryNodeController, socketSVC contracts.ISocketService, logger contracts.ILogger) *Router {
	return &Router{
		router:              mux.NewRouter(),
		roomWSRouter:        newRoomWSRouter(roomWSCtrl, socketSVC, logger),
		auxiliaryNodeRouter: newAuxiliaryNodeRouter(auxiliaryNodeCtrl),
	}
}

func (r *Router) RegisterRoutes() error {
	/*r.router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			println("start")
			request.Header.Add("startTime", strconv.FormatInt(time.Now().UnixNano(), 10))
			handler.ServeHTTP(writer, request)
		})
	})*/
	r.roomWSRouter.registerRoutes(r.router)
	r.auxiliaryNodeRouter.registerRoutes(r.router)

	r.router.PathPrefix("/").Handler(http.FileServer(http.Dir("views/")))
	println("routes registered")
	return nil
}

func (r *Router) Serve(addr string) error {
	println("[HTTP] serving on", addr)
	return http.ListenAndServe(addr, r.router)
}
