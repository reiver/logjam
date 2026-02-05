package verboten

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/reiver/go-etag"

	"codeberg.org/greatape/logjam/srv/http"
)

const path string = "/ogimage.png"

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

	var digest string
	{
		digestBytes := sha256.Sum256(ogimage)
		digest = hex.EncodeToString(digestBytes[:])
	}
	log.Debugf("digest: %s", digest)

	var eTag string = "sha256-" + digest
	log.Debugf("eTag: %s", eTag)

	if etag.Handle(responsewriter, request, eTag) {
		log.Debug("etag caching HIT")
		return
	} else {
		log.Debug("etag caching MISS")
	}

	_, err := responsewriter.Write(ogimage)
	if nil != err {
		log.Errorf("problem writing ogimage.png content to client: %s", err)
	}
}
