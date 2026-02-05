package cfg

import (
	"codeberg.org/greatape/logjam/flg"
	libcfg "codeberg.org/greatape/logjam/lib/cfg"
)

var config libcfg.Model = libcfg.Model{
	GoldGorillaBaseURL:flg.GoldGorillaBaseURL,
	ProdMode:flg.ProdMode,
	WebServerTCPAddress:flg.WebServerTCPAddress,
}

var Config libcfg.Configurer = libcfg.Wrap(config)
