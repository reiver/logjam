package verboten

import (
	"encoding/json"
	"errors"
	"github.com/reiver/logjam/lib/neynar"
	"github.com/reiver/logjam/lib/rest"
	httpsrv "github.com/reiver/logjam/srv/http"
	neynarsrv "github.com/reiver/logjam/srv/neynar"
	"io"
	"net/http"
)

const path string = "/hapi/v1/neynar/submit"

func init() {
	httpsrv.Router.HandleFunc(path, serveHTTP).Methods(http.MethodPost, http.MethodOptions)
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
	var reqModel neynar.AK
	err = json.Unmarshal(reqBody, &reqModel)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}

	ok, err := neynarsrv.Repository.VerifySigner(reqModel.SignerUUID)
	if rest.HandleIfErr(responsewriter, err, http.StatusInternalServerError) {
		return
	}
	if !ok {
		rest.Error(responsewriter, errors.New("invalid signerUUID"), http.StatusUnauthorized)
		return
	}

	err = neynarsrv.Repository.SaveAccountKeys(reqModel)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	_ = rest.Write(responsewriter, nil, http.StatusCreated)
}
