package routers

import (
	"github.com/gorilla/mux"
	"github.com/sparkscience/logjam/controllers"
	"net/http"
)

type auxiliaryNodeRouter struct {
	ctrl *controllers.AuxiliaryNodeController
}

func newAuxiliaryNodeRouter(ctrl *controllers.AuxiliaryNodeController) IRouteRegistrar {
	return &auxiliaryNodeRouter{
		ctrl: ctrl,
	}
}

func (r *auxiliaryNodeRouter) registerRoutes(router *mux.Router) {
	prefix := "/auxiliary-node"
	router.HandleFunc(prefix+"/ice", r.ctrl.SendICECandidate).
		Methods(http.MethodPost)
	router.HandleFunc(prefix+"/join", r.ctrl.Join).
		Methods(http.MethodGet)
}
