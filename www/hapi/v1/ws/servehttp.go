package verboten

import (
	"net/http"

	"github.com/reiver/logjam/srv/http"
	tempTODO "github.com/reiver/logjam/www/ACCT/conf"
)

const path string = "/hapi/v1/ws"

func init() {
	httpsrv.Router.HandleFunc(path, tempTODO.ServeHTTP).Methods(http.MethodGet, http.MethodOptions)
}
