package main

import (
	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/lib/logs"
	"github.com/reiver/logjam/routers"
	"github.com/reiver/logjam/srv/log"
	"github.com/reiver/logjam/srv/room"
	"github.com/reiver/logjam/srv/websock"
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
	roomRepo := roomsrv.Repository
	socketSVC := websocksrv.WebSockSrv
	app.Router = routers.NewRouter(roomRepo, socketSVC, logger)
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
