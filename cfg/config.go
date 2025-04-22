package cfg

import (
	"github.com/reiver/logjam/flg"
	libcfg "github.com/reiver/logjam/lib/cfg"
)

var config libcfg.Model = libcfg.Model{
	GoldGorillaBaseURL:  flg.GoldGorillaBaseURL,
	ProdMode:            flg.ProdMode,
	WebServerTCPAddress: flg.WebServerTCPAddress,
	BlueSkyBaseURL:      flg.BlueSkyBaseURL,
	NeynarApiKey:        flg.NeynarApiKey,
}

var Config libcfg.Configurer = libcfg.Wrap(config)
