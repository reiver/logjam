package main

import (
	"github.com/reiver/logjam/srv/log"
	_ "github.com/reiver/logjam/www"
)

func main() {
	log := logsrv.Tag("main")

	log.Info("LogJam ⚡")
	blur()

	log.Info("Here we go…")
	webserve()
}
