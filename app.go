package main

import (
	"github.com/sparkscience/logjam/controllers"
	"github.com/sparkscience/logjam/models"
	"github.com/sparkscience/logjam/models/contracts"
	"github.com/sparkscience/logjam/models/repositories/goldgorilla"
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
		GoldGorillaSVCAddr: "",
		SrcListenAddr:      srcListenAddr,
	}
	ggSVCRepo := GoldGorilla.NewGoldGorillaRepository()
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
