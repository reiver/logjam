package verboten

import (
	"net/http"

	"github.com/reiver/go-nodeinfo"
	nodeinfo2 "github.com/reiver/go-nodeinfo/2.0"
	"github.com/reiver/go-path"

	"codeberg.org/greatape/logjam/cfg"
	"codeberg.org/greatape/logjam/srv/http"
)

func Path() string {
	return path.Join(nodeinfo.DefaultPath, "2.0")
}

func init() {
	var httphandler http.Handler = nodeinfo2.NodeInfo{
		MetaData: map[string]any{
			"nodeName":        cfg.InstanceName(),
			"nodeDescription": cfg.InstanceDescription(),
		},
		Protocols: cfg.ProtocolsText(),
		Software: nodeinfo2.Software{
			Name: cfg.SoftwareNameText(),
		},
	}

	httpsrv.Router.Handle(Path(), httphandler)
}
