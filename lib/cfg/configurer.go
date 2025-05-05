package cfg

type Configurer interface {
	NeynarApiKey() string
	BlueSkyBaseURL() string
	GoldGorillaBaseURL() string
	PocketBaseURL() string
	PocketBaseAuthToken() string
	ProdMode() bool
	WebServerTCPAddress() string
}
