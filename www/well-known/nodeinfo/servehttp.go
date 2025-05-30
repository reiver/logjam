package verboten

import (
	"net/http"

	"github.com/reiver/go-nodeinfo"
	libpath "github.com/reiver/go-path"
	"github.com/reiver/go-opt"

	"github.com/reiver/logjam/srv/http"
	nodeinfo2     "github.com/reiver/logjam/www/well-known/nodeinfo/2.0"
	nodeinfo2dot1 "github.com/reiver/logjam/www/well-known/nodeinfo/2.1"
	nodeinfo2dot2 "github.com/reiver/logjam/www/well-known/nodeinfo/2.2"
)

const path string = nodeinfo.DefaultPath

func init() {
	var httphandler http.Handler = nodeinfo.WellKnown{
		NodeInfo1:     opt.Something(libpath.Join(nodeinfo.DefaultPath, "1.0")),
		NodeInfo1Dot1: opt.Something(libpath.Join(nodeinfo.DefaultPath, "1.1")),
		NodeInfo2:     opt.Something(nodeinfo2.Path()),
		NodeInfo2Dot1: opt.Something(nodeinfo2dot1.Path()),
		NodeInfo2Dot2: opt.Something(nodeinfo2dot2.Path()),
	}

	httpsrv.Router.Handle(path, httphandler)
}
