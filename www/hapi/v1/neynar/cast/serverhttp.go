package verboten

import (
	"encoding/json"
	"github.com/reiver/logjam/lib/neynar"
	"github.com/reiver/logjam/lib/rest"
	httpsrv "github.com/reiver/logjam/srv/http"
	neynarsrv "github.com/reiver/logjam/srv/neynar"
	"io"
	"net/http"
)

const path string = "/hapi/v1/neynar/cast"

func init() {
	httpsrv.Router.HandleFunc(path, serveHTTP).Methods(http.MethodPost, http.MethodOptions)
}

type createCastRequestModel struct {
	neynar.CastPayload
	FID        int64  `json:"fid"`
	SignerUUID string `json:"refreshToken"`
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
	var req createCastRequestModel
	err = json.Unmarshal(reqBody, &req)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}

	err = neynarsrv.Repository.CreateCast(req.FID, neynar.CastPayload{
		Text:      req.Text,
		ParentURL: req.ParentURL,
		Embeds:    req.Embeds,
	}, req.SignerUUID)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	_ = rest.Write(responsewriter, nil, 204)
}
