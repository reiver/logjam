package main

import (
	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/srv/log"
	_ "github.com/reiver/logjam/www"
)

func main() {
	log := logsrv.Tag("main")

	log.Info("LogJam ⚡")
	blur()

	app := App{}

	app.Init(cfg.Config)
	app.Run()
}
