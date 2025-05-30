package verboten

import (
	"net/http"

	"github.com/reiver/go-nodeinfo"
	libpath "github.com/reiver/go-path"
	"github.com/reiver/go-opt"

	"github.com/reiver/logjam/srv/http"
)

const path string = nodeinfo.DefaultPath

func init() {
	var httphandler http.Handler = nodeinfo.WellKnown{
		NodeInfo1:     opt.Something(libpath.Join(nodeinfo.DefaultPath, "1.0")),
		NodeInfo1Dot1: opt.Something(libpath.Join(nodeinfo.DefaultPath, "1.1")),
		NodeInfo2:     opt.Something(libpath.Join(nodeinfo.DefaultPath, "2.0")),
		NodeInfo2Dot1: opt.Something(libpath.Join(nodeinfo.DefaultPath, "2.1")),
		NodeInfo2Dot2: opt.Something(libpath.Join(nodeinfo.DefaultPath, "2.2")),
	}

	httpsrv.Router.Handle(path, httphandler)
}
