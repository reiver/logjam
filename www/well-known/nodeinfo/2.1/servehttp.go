package verboten

import (
	"net/http"

	"github.com/reiver/go-nodeinfo"
	nodeinfo2dot1 "github.com/reiver/go-nodeinfo/2.1"
	"github.com/reiver/go-opt"
	"github.com/reiver/go-path"

	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/srv/http"
)

func Path() string {
	return path.Join(nodeinfo.DefaultPath, "2.1")
}

func init() {
	var httphandler http.Handler = nodeinfo2dot1.NodeInfo{
		MetaData: map[string]any{
			"nodeName":        cfg.InstanceName(),
			"nodeDescription": cfg.InstanceDescription(),
		},
		Protocols: cfg.ProtocolsText(),
		Software: nodeinfo2dot1.Software{
			Name:     cfg.SoftwareNameText(),
			HomePage: opt.Something(cfg.SoftwareHomePageURL()),
		},
	}

	httpsrv.Router.Handle(Path(), httphandler)
}
