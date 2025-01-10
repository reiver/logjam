package verboten

import (
	"io"
	"io/fs"
	"net/http"

	"github.com/reiver/go-erorr"
	"github.com/reiver/go-path"

	"github.com/reiver/logjam/srv/http"
	"github.com/reiver/logjam/web-app"
)

const pathprefix string = "/src/"

func init() {
	httpsrv.Router.PathPrefix(pathprefix).HandlerFunc(serveHTTP)
}

func serveHTTP(responsewriter http.ResponseWriter, request *http.Request) {
	log.Debug("BEGIN")
	defer log.Debug("END")

	if nil == responsewriter {
		log.Error("nil response-writer")
		return
	}
	if nil == request {
		const code int = http.StatusInternalServerError
		http.Error(responsewriter, http.StatusText(code), code)
		log.Error("nil request")
		return
	}
	if nil == request.URL {
		const code int = http.StatusInternalServerError
		http.Error(responsewriter, http.StatusText(code), code)
		log.Error("nil request-url")
		return
	}

	var requestpath string = request.URL.Path

	log.Debugf("request-method = %q", request.Method)
	log.Debugf("request-path   = %q", requestpath)

	var fspath string = path.Canonical(requestpath)
	if 0 < len(fspath)&& '/' == fspath[0] {
		fspath = fspath[1:]
	}
	log.Debugf("fspath         = %q", requestpath)

	if http.MethodGet != request.Method {
		const code int = http.StatusMethodNotAllowed
		http.Error(responsewriter, http.StatusText(code), code)
		log.Debugf("method not allowed: %q", request.Method)
		return
	}

	file, err := webapp.SrcFS.Open(fspath)
	if nil != err {
		switch {
		case erorr.Is(err, fs.ErrNotExist):
			const code int = http.StatusNotFound
			http.Error(responsewriter, http.StatusText(code), code)
			log.Debugf("file %q not found", fspath)
			return
		default:
			const code int = http.StatusInternalServerError
			http.Error(responsewriter, http.StatusText(code), code)
			log.Errorf("problem opening %q: (%T) %s", fspath, err, err)
			return
		}
	}

	_, err = io.Copy(responsewriter, file)
	if nil != err {
		log.Errorf("problem sending content of file %q to client", fspath)
	}
}
