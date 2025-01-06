package main

import (
	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/srv/log"
)

func main() {
	logsrv.Info("main", "GreatApe LogJam ⚡")

	app := App{}

	app.Init(cfg.Config)
	app.Run()
}
