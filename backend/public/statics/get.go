package statics

import (
	"github.com/sparkscience/logjam/backend/cfg"
	"github.com/sparkscience/logjam/backend/lib/handlers/statichandler"
	httproutersrv "github.com/sparkscience/logjam/backend/srv/http/router"
	logsrv "github.com/sparkscience/logjam/backend/srv/log"
)

const Path = "/"
const Method = "GET"

func init() {
	log := logsrv.Begin()
	defer log.End()

	// err := httproutersrv.Router.DelegatePath(statichandler.Handler(cfg.WebStaticsPath), Path, Method)
	// if nil != err {
	// 	log.Errorf("Could not register http.Handler for %q %q: %v", Method, Path, err)
	// 	panic(err)
	// }

	err := httproutersrv.Router.DelegatePath(statichandler.Handler(cfg.LogsPath), "/logs", Method)
	if nil != err {
		log.Errorf("Could not register http.Handler for %q %q: %v", Method, Path, err)
		panic(err)
	}
}
