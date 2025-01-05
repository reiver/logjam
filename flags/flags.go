package flags

import (
	"flag"
	"fmt"
	"os"

	"github.com/reiver/logjam/env"
)

var (
	GoldGorillaBaseURL string
	help bool
	ProdMode bool
	WebServerTCPAddress string
)

func init() {
	var defaultSrc string = fmt.Sprintf(":%s", env.TcpPort)

	flag.StringVar(&GoldGorillaBaseURL, "goldgorilla-svc-addr", "http://localhost:8080", "goldgorilla service address baseurl")
	flag.BoolVar(&help, "h", false, "print help")
	flag.BoolVar(&ProdMode, "prod", false, "enable production mode ( its in dev mode by default )")
	flag.StringVar(&WebServerTCPAddress, "src", defaultSrc, "source listen address")

	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}
}
