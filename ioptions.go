package rpc

type Options interface {
	URIConnection() string
}

type RabbitURI string

func (this RabbitURI) URIConnection() string {
	return string(this)
}
