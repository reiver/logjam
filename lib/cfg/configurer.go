package cfg

type Configurer interface {
	GoldGorillaBaseURL() string
	ProdMode() bool
	WebServerTCPAddress() string
}
