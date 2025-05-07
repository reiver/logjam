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

const path string = "/hapi/v1/bluesky/submit"

func init() {
	httpsrv.RouterWithAuth.HandleFunc(path, serveHTTP).Methods(http.MethodPost, http.MethodOptions)
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
	user := rest.GetUser(request)
	if user == nil {
		http.Error(responsewriter, "not authenticated", http.StatusInternalServerError)
		return
	}
	err = blueskysrv.Repository.SaveLastTokens(reqModel, user.ID)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	_ = rest.Write(responsewriter, nil, 204)
}
