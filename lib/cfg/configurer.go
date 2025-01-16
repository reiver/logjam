package cfg

type Configurer interface {
	GoldGorillaBaseURL() string
	PocketBaseURL() string
	ProdMode() bool
	WebServerTCPAddress() string
}
