package cfg

// wrapper exists to prevent a certain type of bug.
//
// Specially, a bug where some code changes the value of a config.
type wrapper struct {
	internal Model
}

var _ Configurer = wrapper{}

func Wrap(model Model) Configurer {
	return wrapper{model}
}

func (receiver wrapper) GoldGorillaBaseURL() string {
	return receiver.internal.GoldGorillaBaseURL
}

func (receiver wrapper) ProdMode() bool {
	return receiver.internal.ProdMode
}

func (receiver wrapper) WebServerTCPAddress() string {
	return receiver.internal.WebServerTCPAddress
}
