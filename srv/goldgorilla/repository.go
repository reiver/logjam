package goldgorillasrv

import (
	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/lib/goldgorilla"
)

var Repository goldgorilla.IGoldGorillaServiceRepository = goldgorilla.NewHTTPRepository(cfg.Config.GoldGorillaBaseURL())
