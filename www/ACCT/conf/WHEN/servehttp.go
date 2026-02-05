package verboten

import (
	"fmt"
	"io"
	"net/http"

	"github.com/reiver/go-actsock"
	"github.com/reiver/go-fediverseid"
	"github.com/reiver/go-http400"
	"github.com/reiver/go-http500"
	libpath "github.com/reiver/go-path"

	"codeberg.org/greatape/logjam/srv/http"
)

const acct string = "acct"
const when string = "when"

const path string = "/{"+acct+"}/conf/{"+when+"}"

func init() {
	httpsrv.Router.HandleFunc(path, ServeHTTP).Methods(http.MethodGet, http.MethodOptions)
}

func ServeHTTP(responsewriter http.ResponseWriter, request *http.Request) {
	if nil == responsewriter {
		log.Error("nil response-writer")
		return
	}
	if nil == request {
		http500.InternalServerError(responsewriter, request)
		log.Error("nil request")
		return
	}
	if nil == request.URL {
		http500.InternalServerError(responsewriter, request)
		log.Error("nil request-url")
		return
	}

	responsewriter.Header().Add("Access-Control-Allow-Origin", "*")

	var account string
	var unixtime string
	{
		vars := httpsrv.Vars(request)
		if 0 < len(vars) {
			account = vars[acct]
			unixtime = vars[when]
		}

		// For backwards compatibility reasons.
		// Can remove later.
		if "" == account {
			account = request.URL.Query().Get("room")
		}

		if "" == account {
			http400.BadRequest(responsewriter, request)
			log.Debugf("HTTP Bad Request — acct == %q", account)
			return
		}

		log.Debugf("account (fediverse-id): %q", account)
		log.Debugf("unix-time timestamp (seconds): %q", unixtime)
	}

	var acctURI string
	{
		fediverseID, err := fediverseid.ParseFediverseIDString(account)
		if nil != err {
			http400.BadRequest(responsewriter, request)
			log.Debugf("HTTP Bad Request — bad account — acct == %q: %s", account, err)
			return
		}

		acctURI = fediverseID.AcctURI()
	}

	var id string
	var inoutbox string
	{
		var uri = *request.URL
		uri.User = nil
		uri.Scheme = "https"
		uri.Host = request.Host

		id = uri.String()

		uri.Scheme  = "wss"
		uri.Path = libpath.Canonical(libpath.RemoveTrailingSeparators(libpath.Parent(uri.Path)))
		inoutbox = uri.String()
	}

	var object = actsock.Conference{
		Actor: acctURI,
		EndPoints: map[string]string{
			"inoutbox":inoutbox,
		},
		ID: id,
		Name: fmt.Sprintf("%s — GreatApe", account),
	}

	responsewriter.Header().Set("Content-Type", "application/activity+json")
	io.WriteString(responsewriter, object.String())
}
