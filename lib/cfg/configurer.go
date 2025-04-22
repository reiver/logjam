package cfg

type Configurer interface {
	NeynarApiKey() string
	BlueSkyBaseURL() string
	GoldGorillaBaseURL() string
	PocketBaseURL() string
	ProdMode() bool
	WebServerTCPAddress() string
}
