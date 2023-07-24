package routers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sparkscience/logjam/controllers"
	"github.com/sparkscience/logjam/models/contracts"
	"net/http"
	"strconv"
	"time"
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
	r.router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			request.Header.Add("startTime", strconv.FormatInt(time.Now().UnixMicro(), 10))
			handler.ServeHTTP(writer, request)
		})
	})
	r.roomWSRouter.registerRoutes(r.router)
	r.auxiliaryNodeRouter.registerRoutes(r.router)

	r.router.PathPrefix("/").Handler(http.FileServer(http.Dir("views/")))
	r.router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			//startTime, err := strconv.ParseUint(request.Header.Get("startTime"), 10, 64)
			//if err != nil {
			//println(err.Error())
			//return
			//}
			/*a := time.UnixMicro(int64(startTime))
			b := time.Now()
			diff := b.Sub(a)
			println(diff.String())
			r.logger.Log("HTTP", contracts.LDebug, fmt.Sprintf(`%ds %s`, diff.Milliseconds(), request.URL.Path))*/
			r.logger.Log("HTTP", contracts.LDebug, fmt.Sprintf(`%s`, request.URL.Path))
			handler.ServeHTTP(writer, request)
		})
	})
	println("routes registered")
	return nil
}

func (r *Router) Serve(addr string) error {
	println("[HTTP] serving on", addr)
	return http.ListenAndServe(addr, r.router)
}
