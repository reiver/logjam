package main

import (
	"github.com/reiver/logjam/controllers"
	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/lib/logs"
	roomRepository "github.com/reiver/logjam/models/repositories/room"
	"github.com/reiver/logjam/routers"
	"github.com/reiver/logjam/services"
	"github.com/reiver/logjam/srv/goldgorilla"
	"github.com/reiver/logjam/srv/log"
)

type App struct {
	Logger logs.Logger
	Router *routers.Router

	config cfg.Configurer
}

func (app *App) Init(config cfg.Configurer) {
	logger := logsrv.Logger

	app.Logger = logger.Tag("app")
	app.Logger.Info("app", "initializing logjam ..")
	app.config = config
	ggSVCRepo := goldgorillasrv.Repository
	roomRepo := roomRepository.NewRoomRepository()
	socketSVC := services.NewSocketService(logger)
	roomWSCtrl := controllers.NewRoomWSController(socketSVC, roomRepo, ggSVCRepo, logger)
	restHelper := &controllers.RestResponseHelper{}
	goldGorillaCtrl := controllers.NewGoldGorillaController(roomRepo, ggSVCRepo, socketSVC, app.config, restHelper, logger)
	app.Router = routers.NewRouter(roomWSCtrl, goldGorillaCtrl, roomRepo, socketSVC, logger)
	panicIfErr(app.Router.RegisterRoutes())
}

func (app *App) Run() {
	app.Logger.Info("app", "running ..")
	panicIfErr(app.Router.Serve(app.config.WebServerTCPAddress()))
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
