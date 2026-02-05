package main

import (
	"codeberg.org/greatape/logjam/srv/log"

	// This import enables all the HTTP handlers.
	_ "codeberg.org/greatape/logjam/www"
)

func main() {
	log := logsrv.Tag("main")

	log.Info("LogJam ⚡")
	blur()

	log.Info("Here we go…")
	webserve()
}
