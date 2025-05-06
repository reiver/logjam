package verboten

import (
	"encoding/json"
	"github.com/reiver/logjam/lib/rest"
	"github.com/reiver/logjam/lib/users"
	httpsrv "github.com/reiver/logjam/srv/http"
	userssrv "github.com/reiver/logjam/srv/users"
	"io"
	"net/http"
)

const path string = "/hapi/v1/users/otp"

func init() {
	httpsrv.Router.HandleFunc(path, serveHTTP).Methods(http.MethodPost, http.MethodOptions)
}

type createOtpReqModel struct {
	Email string `json:"email"`
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
	var req createOtpReqModel
	err = json.Unmarshal(reqBody, &req)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}

	err = userssrv.Repository.CreateOTP(users.CreateOTPDTO{
		Email: req.Email,
		Code:  "123456",
	})
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	_ = rest.Write(responsewriter, nil, 204)
}
