package flags

import (
	"flag"
	"fmt"
	"os"

	"github.com/reiver/logjam/env"
)

var (
	GoldGorillaBaseURL string
	WebServerTCPAddress string
	ProdMode bool
	Help bool
)

func init() {
	var defaultSrc string = fmt.Sprintf(":%s", env.TcpPort)

	flag.StringVar(&GoldGorillaBaseURL, "goldgorilla-svc-addr", "http://localhost:8080", "goldgorilla service address baseurl")
	flag.StringVar(&WebServerTCPAddress, "src", defaultSrc, "source listen address")
	flag.BoolVar(&ProdMode, "prod", false, "enable production mode ( its in dev mode by default )")
	flag.BoolVar(&Help, "h", false, "print help")

	flag.Parse()

	if Help {
		flag.PrintDefaults()
		os.Exit(0)
	}
}
