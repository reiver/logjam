package cfg

type Configurer interface {
	GoldGorillaBaseURL() string
	WebServerTCPAddress() string
}
