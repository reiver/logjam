package main

import (
	"github.com/reiver/logjam/controllers"
	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/lib/logs"
	"github.com/reiver/logjam/models/repositories/goldgorilla"
	roomRepository "github.com/reiver/logjam/models/repositories/room"
	"github.com/reiver/logjam/routers"
	"github.com/reiver/logjam/services"
	"github.com/reiver/logjam/srv/log"
)

type App struct {
	Logger logs.Logger
	Router *routers.Router

	config cfg.Configurer
}

func (app *App) Init(config cfg.Configurer) {
	app.Logger = logsrv.Logger
	app.Logger.Info("app", "initializing logjam ..")
	app.config = config
	ggSVCRepo := GoldGorillaRepository.NewHTTPRepository(config.GoldGorillaBaseURL())
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
	panicIfErr(app.Router.Serve(app.config.WebServerTCPAddress()))
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
