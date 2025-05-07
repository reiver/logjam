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

const path string = "/hapi/v1/users/profile"

func init() {
	httpsrv.RouterWithAuth.HandleFunc(path, serveHTTP).Methods(http.MethodPatch, http.MethodOptions)
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
	var req users.UpdateProfileDTO
	err = json.Unmarshal(reqBody, &req)
	if rest.HandleIfErr(responsewriter, err, 400) {
		return
	}
	user := rest.GetUser(request)
	if user == nil {
		http.Error(responsewriter, "not authenticated", http.StatusInternalServerError)
		return
	}
	req.UserID = user.ID
	err = userssrv.Repository.UpdateProfile(req)
	if rest.HandleIfErr(responsewriter, err, 500) {
		return
	}

	_ = rest.Write(responsewriter, nil, 204)
}
