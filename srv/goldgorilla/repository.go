package goldgorillasrv

import (
	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/lib/goldgorilla"
	"github.com/reiver/logjam/models/repositories/goldgorilla"
)

var Repository goldgorilla.IGoldGorillaServiceRepository = GoldGorillaRepository.NewHTTPRepository(cfg.Config.GoldGorillaBaseURL())
