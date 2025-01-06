package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/reiver/logjam/lib/rest"
	"github.com/reiver/logjam/models/dto"
	"github.com/reiver/logjam/srv/http"
	"github.com/reiver/logjam/srv/websock"
)

const path string = "/goldgorilla/ice"

func init() {
        httpsrv.Router.HandleFunc(path, serveHTTP)
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
	var reqModel dto.SendIceCandidateReqModel
	err = json.Unmarshal(reqBody, &reqModel)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}
	_ = websocksrv.WebSockSrv.Send(map[string]interface{}{
		"Type":      "new-ice-candidate",
		"Target":    strconv.FormatUint(reqModel.ID, 10),
		"candidate": reqModel.ICECandidate,
		"data":      strconv.FormatUint(reqModel.GGID, 10),
	}, reqModel.ID)
	_ = rest.Write(responsewriter, nil, 204)
}
