package verboten

import (
	"encoding/json"
	"github.com/reiver/logjam/lib/rest"
	blueskysrv "github.com/reiver/logjam/srv/bluesky"
	httpsrv "github.com/reiver/logjam/srv/http"
	"io"
	"net/http"
)

const path string = "/bluesky/post"

func init() {
	httpsrv.Router.HandleFunc(path, serveHTTP).Methods(http.MethodPost, http.MethodOptions)
}

type createPostRequestModel struct {
	DID  string `json:"did"`
	Text string `json:"text"`
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
	var reqModel createPostRequestModel
	err = json.Unmarshal(reqBody, &reqModel)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}

	err = blueskysrv.Repository.CreatePost(reqModel.DID, reqModel.Text)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	_ = rest.Write(responsewriter, nil, http.StatusNoContent)
}
