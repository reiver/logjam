package verboten

import (
	"net/http"

	"github.com/reiver/go-nodeinfo"
	"github.com/reiver/go-opt"

	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/srv/http"
	nodeinfo2     "github.com/reiver/logjam/www/well-known/nodeinfo/2.0"
	nodeinfo2dot1 "github.com/reiver/logjam/www/well-known/nodeinfo/2.1"
	nodeinfo2dot2 "github.com/reiver/logjam/www/well-known/nodeinfo/2.2"
)

const path string = nodeinfo.DefaultPath

const scheme string = "https"

func init() {
	nodeinfo.SetServerText(cfg.HTTPServerText())

	var httphandler http.Handler = nodeinfo.ResolvingWellKnown{
		WellKnown: nodeinfo.WellKnown{
			NodeInfo2:     opt.Something(nodeinfo2.Path()),
			NodeInfo2Dot1: opt.Something(nodeinfo2dot1.Path()),
			NodeInfo2Dot2: opt.Something(nodeinfo2dot2.Path()),
		},
		Scheme: scheme,
	}

	httpsrv.Router.Handle(path, httphandler)
}
