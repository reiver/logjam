package main

import (
	"sourcecode.social/greatape/logjam/controllers"
	"sourcecode.social/greatape/logjam/models"
	"sourcecode.social/greatape/logjam/models/contracts"
	"sourcecode.social/greatape/logjam/models/repositories/goldgorilla"
	roomRepository "sourcecode.social/greatape/logjam/models/repositories/room"
	"sourcecode.social/greatape/logjam/routers"
	"sourcecode.social/greatape/logjam/services"
	"sourcecode.social/greatape/logjam/services/logger"
)

type App struct {
	Logger contracts.ILogger
	Router *routers.Router

	config *models.ConfigModel
}

func (app *App) Init(srcListenAddr string, prodMode bool, goldGorillaSVCAddr string) {
	app.Logger = logger.NewSTDOUTLogger(prodMode)
	_ = app.Logger.Log("app", contracts.LInfo, "initializing logjam ..")
	app.config = &models.ConfigModel{
		GoldGorillaSVCAddr: goldGorillaSVCAddr,
		SrcListenAddr:      srcListenAddr,
	}
	ggSVCRepo := GoldGorillaRepository.NewHTTPRepository(goldGorillaSVCAddr)
	roomRepo := roomRepository.NewRoomRepository()
	socketSVC := services.NewSocketService(app.Logger)
	roomWSCtrl := controllers.NewRoomWSController(socketSVC, roomRepo, ggSVCRepo, app.Logger)
	restHelper := &controllers.RestResponseHelper{}
	goldGorillaCtrl := controllers.NewGoldGorillaController(roomRepo, ggSVCRepo, socketSVC, app.config, restHelper, app.Logger)
	app.Router = routers.NewRouter(roomWSCtrl, goldGorillaCtrl, roomRepo, socketSVC, app.Logger)
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
