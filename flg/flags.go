package flg

import (
	"flag"
	"fmt"
	"os"

	"github.com/reiver/go-erorr"

	"github.com/reiver/logjam/env"
)

var (
	GoldGorillaBaseURL string
	help bool
	PocketBaseURL string
	ProdMode bool
	WebServerTCPAddress string
)

const (
	pocketBaseURLFlag = "pburl"
)

func init() {
	var defaultSrc string = fmt.Sprintf(":%s", env.TcpPort)

	flag.StringVar(&GoldGorillaBaseURL, "goldgorilla-svc-addr", "http://localhost:8080", "goldgorilla service address baseurl")
	flag.BoolVar(&help, "h", false, "print help")
	flag.StringVar(&PocketBaseURL, pocketBaseURLFlag, env.PocketBaseURL, "pocketbase base API URL")
	flag.BoolVar(&ProdMode, "prod", false, "enable production mode ( its in dev mode by default )")
	flag.StringVar(&WebServerTCPAddress, "src", defaultSrc, "source listen address")

	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if "" == PocketBaseURL {
		var err error = erorr.Errorf("PocketBase URL cannot be empty. Must be set with %q environment-variable or --%s command-line flag / switch.", env.PocketBaseURLEnvVar, pocketBaseURLFlag)
		panic(err)
	}
}
