package statics

import (
	"github.com/sparkscience/logjam/backend/lib/handlers/statichandler"
	httproutersrv "github.com/sparkscience/logjam/backend/srv/http/router"
)

const Path = "/"

func init() {
	httproutersrv.Router.DelegatePath(statichandler.Handler, Path, "GET")
}
