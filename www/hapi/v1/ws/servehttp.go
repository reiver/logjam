package verboten

import (
	"net/http"

	"codeberg.org/greatape/logjam/srv/http"
	tempTODO "codeberg.org/greatape/logjam/www/ACCT/conf"
)

const path string = "/hapi/v1/ws"

func init() {
	httpsrv.Router.HandleFunc(path, tempTODO.ServeHTTP).Methods(http.MethodGet, http.MethodOptions)
}
