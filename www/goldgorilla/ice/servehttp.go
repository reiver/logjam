package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"codeberg.org/greatape/logjam/lib/goldgorilla"
	"codeberg.org/greatape/logjam/lib/msgs"
	"codeberg.org/greatape/logjam/lib/rest"
	"codeberg.org/greatape/logjam/srv/http"
	"codeberg.org/greatape/logjam/srv/websock"
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
	var reqModel goldgorilla.SendIceCandidateReqModel
	err = json.Unmarshal(reqBody, &reqModel)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}
	_ = websocksrv.WebSockSrv.Send(map[string]interface{}{
		"type":      msgs.TypeNewIceCandidate,
		"Target":    strconv.FormatUint(reqModel.ID, 10),
		"candidate": reqModel.ICECandidate,
		"data":      strconv.FormatUint(reqModel.GGID, 10),
	}, reqModel.ID)
	_ = rest.Write(responsewriter, nil, 204)
}
