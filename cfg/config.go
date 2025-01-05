package cfg

import (
	"github.com/reiver/logjam/flags"
	libcfg "github.com/reiver/logjam/lib/cfg"
)

var config libcfg.Model = libcfg.Model{
	SrcListenAddr:flags.Src,
	GoldGorillaSVCAddr:flags.GoldGorillaSVCAddr,
}

var Config libcfg.Configurer = libcfg.Wrap(config)
