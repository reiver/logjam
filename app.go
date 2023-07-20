package main

import (
	"github.com/sparkscience/logjam/controllers"
	"github.com/sparkscience/logjam/models/contracts"
	"github.com/sparkscience/logjam/models/repositories"
	"github.com/sparkscience/logjam/routers"
	"github.com/sparkscience/logjam/services"
	"github.com/sparkscience/logjam/services/logger"
)

type App struct {
	Logger contracts.ILogger
	Router *routers.Router

	srcListenAddr string
}

func (app *App) Init(srcListenAddr string, prodMode bool) {
	app.Logger = logger.NewSTDOUTLogger(prodMode)
	_ = app.Logger.Log("app", contracts.LInfo, "initializing logjam ..")
	app.srcListenAddr = srcListenAddr

	roomRepo := repositories.NewRoomRepository()
	socketSVC := services.NewSocketService(app.Logger)
	roomWSCtrl := controllers.NewRoomWSController(socketSVC, roomRepo, app.Logger)
	auxiliaryNodeCtrl := controllers.NewAuxiliaryNodeController(roomRepo, socketSVC, app.Logger)
	app.Router = routers.NewRouter(roomWSCtrl, auxiliaryNodeCtrl, socketSVC, app.Logger)
	panicIfErr(app.Router.RegisterRoutes())
}

func (app *App) Run() {
	_ = app.Logger.Log("app", contracts.LInfo, "running ..")
	panicIfErr(app.Router.Serve(app.srcListenAddr))
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
