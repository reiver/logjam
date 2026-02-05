package verboten

import (
	"net/http"

	"github.com/reiver/go-nodeinfo"
	nodeinfo2dot2 "github.com/reiver/go-nodeinfo/2.2"
	"github.com/reiver/go-opt"
	"github.com/reiver/go-path"

	"codeberg.org/greatape/logjam/cfg"
	"codeberg.org/greatape/logjam/srv/http"
)

func Path() string {
	return path.Join(nodeinfo.DefaultPath, "2.2")
}

func init() {
	var httphandler http.Handler = nodeinfo2dot2.NodeInfo{
		Instance: nodeinfo2dot2.Instance{
			Name:        opt.Something(cfg.InstanceName()),
			Description: opt.Something(cfg.InstanceDescription()),
		},
		Protocols: cfg.ProtocolsText(),
		Software: nodeinfo2dot2.Software{
			Name:     cfg.SoftwareNameText(),
			HomePage: opt.Something(cfg.SoftwareHomePageURL()),
		},
	}

	httpsrv.Router.Handle(Path(), httphandler)
}
