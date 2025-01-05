package flags

import (
	"flag"
	"os"
)

var (
	Src string
	GoldGorillaSVCAddr string
	ProdMode bool
	Help bool
)

func init() {
	flag.StringVar(&Src, "src", "0.0.0.0:8080", "source listen address")
	flag.StringVar(&GoldGorillaSVCAddr, "goldgorilla-svc-addr", "http://localhost:8080", "goldgorilla service address baseurl")
	flag.BoolVar(&ProdMode, "prod", false, "enable production mode ( its in dev mode by default )")
	flag.BoolVar(&Help, "h", false, "print help")

	flag.Parse()

	if Help {
		flag.PrintDefaults()
		os.Exit(0)
	}
}
