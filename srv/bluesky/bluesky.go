package blueskysrv

import (
	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/lib/bluesky"
)

var Repository = bluesky.NewHTTPRepository(cfg.Config.BlueSkyBaseURL())
