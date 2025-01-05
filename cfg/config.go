package cfg

import (
	"github.com/reiver/logjam/flags"
	libcfg "github.com/reiver/logjam/lib/cfg"
)

var config libcfg.Model = libcfg.Model{
	GoldGorillaBaseURL:flags.GoldGorillaBaseURL,
	WebServerTCPAddress:flags.WebServerTCPAddress,
}

var Config libcfg.Configurer = libcfg.Wrap(config)
