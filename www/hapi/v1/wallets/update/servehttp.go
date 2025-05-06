package verboten

import (
	"encoding/json"
	"github.com/reiver/logjam/lib/rest"
	"github.com/reiver/logjam/lib/wallets"
	httpsrv "github.com/reiver/logjam/srv/http"
	walletssrv "github.com/reiver/logjam/srv/wallets"
	"io"
	"net/http"
)

const path string = "/hapi/v1/wallets"

func init() {
	httpsrv.Router.HandleFunc(path, serveHTTP).Methods(http.MethodPatch, http.MethodOptions)
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
	reqBody, err := io.ReadAll(request.Body)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}
	var req wallets.UpdateWalletDTO
	err = json.Unmarshal(reqBody, &req)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}

	req.OwnerId = "" //read from token
	err = walletssrv.Repository.UpdateWallet(req)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	_ = rest.Write(responsewriter, nil, http.StatusNoContent)
}
