package graphapi

import (
	"github.com/sparkscience/logjam/backend/lib/handlers/graphapihandler"
	httproutersrv "github.com/sparkscience/logjam/backend/srv/http/router"
	logsrv "github.com/sparkscience/logjam/backend/srv/log"
)

const Path = "/api/v2/graph"

func init() {
	log := logsrv.Begin()
	defer log.End()

	err := httproutersrv.Router.DelegatePath(graphapihandler.Handler(logsrv.Logger), Path, "GET")
	if nil != err {
		log.Errorf("Could not register http.Handler for %q %q: %v", "GET", Path, err)
		panic(err)
	}
	err = httproutersrv.Router.DelegatePath(graphapihandler.Handler(logsrv.Logger), Path, "POST")
	if nil != err {
		log.Errorf("Could not register http.Handler for %q %q: %v", "POST", Path, err)
		panic(err)
	}
	err = httproutersrv.Router.DelegatePath(graphapihandler.Handler(logsrv.Logger), Path, "DELETE")
	if nil != err {
		log.Errorf("Could not register http.Handler for %q %q: %v", "DELETE", Path, err)
		panic(err)
	}
}
