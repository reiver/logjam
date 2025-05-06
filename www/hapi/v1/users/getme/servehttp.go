package verboten

import (
	"github.com/reiver/logjam/lib/rest"
	httpsrv "github.com/reiver/logjam/srv/http"
	userssrv "github.com/reiver/logjam/srv/users"
	"net/http"
)

const path string = "/hapi/v1/users/me"

func init() {
	httpsrv.Router.HandleFunc(path, serveHTTP).Methods(http.MethodGet, http.MethodOptions)
}

func serveHTTP(responsewriter http.ResponseWriter, request *http.Request) {
	if nil == responsewriter {
		return
	}
	if nil == request {
		const code int = http.StatusInternalServerError
		http.Error(responsewriter, http.StatusText(code), code)
		return
	}

	id := "" // get from token
	resp, err := userssrv.Repository.GetMe(id)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	_ = rest.Write(responsewriter, resp, http.StatusOK)
}
