package main

import (
	"github.com/sparkscience/logjam/controllers"
	"github.com/sparkscience/logjam/models"
	"github.com/sparkscience/logjam/models/contracts"
	"github.com/sparkscience/logjam/models/repositories/auxiliarynode"
	roomRepository "github.com/sparkscience/logjam/models/repositories/room"
	"github.com/sparkscience/logjam/routers"
	"github.com/sparkscience/logjam/services"
	"github.com/sparkscience/logjam/services/logger"
)

type App struct {
	Logger contracts.ILogger
	Router *routers.Router

	config *models.ConfigModel
}

func (app *App) Init(srcListenAddr string, prodMode bool) {
	app.Logger = logger.NewSTDOUTLogger(prodMode)
	_ = app.Logger.Log("app", contracts.LInfo, "initializing logjam ..")
	app.config = &models.ConfigModel{
		AuxiliaryNodeSVCAddr: "",
		SrcListenAddr:        srcListenAddr,
	}
	anSVCRepo := auxiliarynode.NewAuxiliaryNodeRepository()
	roomRepo := roomRepository.NewRoomRepository()
	socketSVC := services.NewSocketService(app.Logger)
	roomWSCtrl := controllers.NewRoomWSController(socketSVC, roomRepo, anSVCRepo, app.Logger)
	restHelper := &controllers.RestResponseHelper{}
	auxiliaryNodeCtrl := controllers.NewAuxiliaryNodeController(roomRepo, anSVCRepo, socketSVC, app.config, restHelper, app.Logger)
	app.Router = routers.NewRouter(roomWSCtrl, auxiliaryNodeCtrl, roomRepo, socketSVC, app.Logger)
	panicIfErr(app.Router.RegisterRoutes())
}

func (app *App) Run() {
	_ = app.Logger.Log("app", contracts.LInfo, "running ..")
	panicIfErr(app.Router.Serve(app.config.SrcListenAddr))
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}