package verboten

import (
	"encoding/json"
	"github.com/reiver/logjam/lib/bluesky"
	"github.com/reiver/logjam/lib/rest"
	blueskysrv "github.com/reiver/logjam/srv/bluesky"
	httpsrv "github.com/reiver/logjam/srv/http"
	"io"
	"net/http"
)

const path string = "/bluesky/submit"

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
	var reqModel bluesky.AK
	err = json.Unmarshal(reqBody, &reqModel)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}

	err = blueskysrv.Repository.SaveLastTokens(reqModel)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	_ = rest.Write(responsewriter, nil, 204)
}
