package goldgorillasrv

import (
	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/models/contracts"
	"github.com/reiver/logjam/models/repositories/goldgorilla"
)

var Repository contracts.IGoldGorillaServiceRepository = GoldGorillaRepository.NewHTTPRepository(cfg.Config.GoldGorillaBaseURL())
