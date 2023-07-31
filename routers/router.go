package routers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sparkscience/logjam/controllers"
	"github.com/sparkscience/logjam/models/contracts"
	"net/http"
	"strings"
)

type IRouteRegistrar interface {
	registerRoutes(router *mux.Router)
}

type Router struct {
	router              *mux.Router
	roomWSRouter        IRouteRegistrar
	auxiliaryNodeRouter IRouteRegistrar
	logger              contracts.ILogger
}

func NewRouter(roomWSCtrl *controllers.RoomWSController, auxiliaryNodeCtrl *controllers.AuxiliaryNodeController, roomRepo contracts.IRoomRepository, socketSVC contracts.ISocketService, logger contracts.ILogger) *Router {
	return &Router{
		router:              mux.NewRouter(),
		roomWSRouter:        newRoomWSRouter(roomWSCtrl, roomRepo, socketSVC, logger),
		auxiliaryNodeRouter: newAuxiliaryNodeRouter(auxiliaryNodeCtrl),
		logger:              logger,
	}
}

func (r *Router) RegisterRoutes() error {
	r.roomWSRouter.registerRoutes(r.router)
	r.auxiliaryNodeRouter.registerRoutes(r.router)

	r.router.PathPrefix("/").Handler(http.FileServer(http.Dir("views/")))
	r.router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ip := request.RemoteAddr
			if strings.Index(ip, ":") > 0 {
				ip = ip[:strings.Index(ip, ":")]
			}
			r.logger.Log("HTTP", contracts.LDebug, fmt.Sprintf(`%s | %s "%s"`, ip, request.Method, request.URL.Path))
			handler.ServeHTTP(writer, request)
		})
	})
	return nil
}

func (r *Router) Serve(addr string) error {
	println("[HTTP] serving on", addr)
	return http.ListenAndServe(addr, r.router)
}
