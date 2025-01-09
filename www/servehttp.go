package verboten

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"github.com/reiver/go-etag"

	"github.com/reiver/logjam/srv/http"
)

const path string = "/"

func init() {
	httpsrv.Router.HandleFunc(path, serveHTTP)
}

func serveHTTP(responsewriter http.ResponseWriter, request *http.Request) {
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

	var host string = request.Host
	if "" == host && nil != request.URL {
		host = request.URL.Host
	}
	if "" == host {
		const code int = http.StatusInternalServerError
		http.Error(responsewriter, http.StatusText(code), code)
		log.Error("empty host")
		return
	}

	var html string
	{
		const needle string = "{{HOST}}"
		html = strings.ReplaceAll(webpage, needle, host)
	}

	var digest string
	{
		digestBytes := sha256.Sum256([]byte(html))
		digest = hex.EncodeToString(digestBytes[:])
	}
	log.Debugf("digest: %s", digest)

	var eTag string = "sha256d-" + digest
	log.Debugf("eTag: %s", eTag)

	if etag.Handle(responsewriter, request, eTag) {
		log.Debug("etag caching HIT")
		return
	} else {
		log.Debug("etag caching MISS")
	}

	_, err := io.WriteString(responsewriter, html)
	if nil != err {
		log.Errorf("problem writing HTML content to client: %s", err)
	}
}
