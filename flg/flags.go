package flg

import (
	"flag"
	"fmt"
	"os"

	"github.com/reiver/logjam/env"
)

var (
	GoldGorillaBaseURL string
	help bool
	InstanceName string
	ProdMode bool
	WebServerTCPAddress string
)

func init() {
	var defaultInstanceName string = env.InstanceName
	var defaultSrc string = fmt.Sprintf(":%s", env.TcpPort)

	flag.StringVar(&GoldGorillaBaseURL, "goldgorilla-svc-addr", "http://localhost:8080", "goldgorilla service address baseurl")
	flag.BoolVar(&help, "h", false, "print help")
	flag.StringVar(&InstanceName, "instance-name", defaultInstanceName, "the name of the instance server")
	flag.BoolVar(&ProdMode, "prod", false, "enable production mode ( its in dev mode by default )")
	flag.StringVar(&WebServerTCPAddress, "src", defaultSrc, "source listen address")

	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}
}
