package main

import (
	"github.com/reiver/logjam/controllers"
	"github.com/reiver/logjam/lib/logs"
	"github.com/reiver/logjam/models"
	"github.com/reiver/logjam/models/repositories/goldgorilla"
	roomRepository "github.com/reiver/logjam/models/repositories/room"
	"github.com/reiver/logjam/routers"
	"github.com/reiver/logjam/services"
)

type App struct {
	Logger logs.Logger
	Router *routers.Router

	config *models.ConfigModel
}

func (app *App) Init(srcListenAddr string, prodMode bool, goldGorillaSVCAddr string) {
	app.Logger = logs.NewStdOutLogger(prodMode)
	app.Logger.Info("app", "initializing logjam ..")
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
	app.Logger.Info("app", "running ..")
	panicIfErr(app.Router.Serve(app.config.SrcListenAddr))
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
