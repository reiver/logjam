package main

import (
	"net/http"

	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/srv/http"
	"github.com/reiver/logjam/srv/log"
	_ "github.com/reiver/logjam/www"
)

func webserve() {
	log := logsrv.Tag("webserve")

	var tcpaddr string = cfg.Config.WebServerTCPAddress()
	log.Infof("serving HTTP on TCP address: %q", tcpaddr)

	err := http.ListenAndServe(tcpaddr, httpsrv.Router)
	if nil != err {
		log.Errorf("ERROR: problem with serving HTTP on TCP address %q: %s", tcpaddr, err)
		panic(err)
	}
}
