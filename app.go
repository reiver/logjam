package main

import (
	"github.com/reiver/logjam/lib/logs"
	"github.com/reiver/logjam/routers"
	"github.com/reiver/logjam/srv/log"
)

type App struct {
	Logger logs.Logger
	Router *routers.Router
}

func (app *App) Init() {
	logger := logsrv.Logger

	app.Logger = logger.Tag("app")
	app.Logger.Info("app", "initializing logjam ..")
	app.Router = routers.NewRouter(logger)
	panicIfErr(app.Router.RegisterRoutes())
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
