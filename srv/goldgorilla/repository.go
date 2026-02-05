package goldgorillasrv

import (
	"codeberg.org/greatape/logjam/cfg"
	"codeberg.org/greatape/logjam/lib/goldgorilla"
)

var Repository goldgorilla.IGoldGorillaServiceRepository = goldgorilla.NewHTTPRepository(cfg.Config.GoldGorillaBaseURL())
