package routers

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/reiver/logjam/controllers"
)

type GoldGorillaRouter struct {
	ctrl *controllers.GoldGorillaController
}

func newGoldGorillaRouter(ctrl *controllers.GoldGorillaController) IRouteRegistrar {
	return &GoldGorillaRouter{
		ctrl: ctrl,
	}
}

func (r *GoldGorillaRouter) registerRoutes(router *mux.Router) {
	prefix := "/goldgorilla"
	router.HandleFunc(prefix+"/ice", r.ctrl.SendICECandidate).
		Methods(http.MethodPost)
	router.HandleFunc(prefix+"/join", r.ctrl.Join).
		Methods(http.MethodPost)
	router.HandleFunc(prefix+"/rejoin", r.ctrl.RejoinGoldGorilla).
		Methods(http.MethodPost)
	router.HandleFunc(prefix+"/offer", r.ctrl.SendOffer).
		Methods(http.MethodPost)
	router.HandleFunc(prefix+"/answer", r.ctrl.SendAnswer).
		Methods(http.MethodPost)

}
