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

const path string = "/hapi/v1/users/signup"

func init() {
	httpsrv.Router.HandleFunc(path, serveHTTP).Methods(http.MethodPost, http.MethodOptions)
}

type completeSignUpReqModel struct {
	Email string `json:"email"`
	Code  string `json:"code"`
	Name  string `json:"name"`
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
	var req completeSignUpReqModel
	err = json.Unmarshal(reqBody, &req)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}

	resp, err := userssrv.Repository.CompleteSignUp(users.CompleteSignUpDTO{
		Email: req.Email,
		Code:  req.Code,
		Name:  req.Name,
	})
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	_ = rest.Write(responsewriter, resp, http.StatusOK)
}
