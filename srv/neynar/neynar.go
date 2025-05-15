package neynarsrv

import (
	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/lib/neynar"
)

var Repository = neynar.NewHTTPRepository("https://api.neynar.com/", cfg.Config.NeynarApiKey())
