package verboten

import (
	"github.com/reiver/logjam/lib/rest"
	httpsrv "github.com/reiver/logjam/srv/http"
	walletssrv "github.com/reiver/logjam/srv/wallets"
	"net/http"
)

const path string = "/hapi/v1/wallets"

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

	id := "" //read from token
	resp, err := walletssrv.Repository.GetUserWallet(id)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	_ = rest.Write(responsewriter, resp, http.StatusOK)
}
