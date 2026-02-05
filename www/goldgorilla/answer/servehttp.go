package verboten

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

const path string = "/goldgorilla/answer"

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
	var reqModel goldgorilla.SetSDPRPCModel
	err = json.Unmarshal(reqBody, &reqModel)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}
	websocksrv.WebSockSrv.Send(map[string]interface{}{
		"type":   msgs.TypeVideoAnswer,
		"target": strconv.FormatUint(reqModel.ID, 10),
		"name":   strconv.FormatUint(reqModel.GGID, 10),
		"sdp":    reqModel.SDP,
		"data":   strconv.FormatUint(reqModel.GGID, 10),
	}, reqModel.ID)
	_ = rest.Write(responsewriter, nil, 204)
}
