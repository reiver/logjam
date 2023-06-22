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
	Logger     contracts.ILogger
	roomRouter *routers.RoomRouter

	listenHost string
	listenPort int
}

func (app *App) Init(listenHost string, listenPort int, mode string) {
	app.Logger = logger.NewSTDOUTLogger()
	_ = app.Logger.Log("app", contracts.LInfo, "initializing logjam ..")
	app.listenHost = listenHost
	app.listenPort = listenPort

	roomRepo := repositories.NewRoomRepository()
	socketSVC := services.NewSocketService(app.Logger)
	roomCtrl := controllers.NewRoomController(socketSVC, roomRepo, app.Logger)
	app.roomRouter = routers.NewRoomRouter(roomCtrl, socketSVC, app.Logger)
}

func (app *App) Run() {
	_ = app.Logger.Log("app", contracts.LInfo, "running ..")
	panicIfErr(app.roomRouter.Serve(app.listenHost, app.listenPort))
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
