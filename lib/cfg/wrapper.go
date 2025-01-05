package cfg

type wrapper struct {
	internal Model
}

var _ Configurer = wrapper{}

func Wrap(model Model) Configurer {
	return wrapper{model}
}


func (receiver wrapper) GoldGorillaBaseURL() string {
	return receiver.internal.GoldGorillaSVCAddr
}

func (receiver wrapper) WebServerTCPAddress() string {
	return receiver.internal.SrcListenAddr
}
