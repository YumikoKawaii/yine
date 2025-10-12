package connection_registry

type Registry interface {
	Register(userIdentification string, serverIdentification string) error
	GetServers(userIdentifications []string) ([]string, error)
}

type redisImpl struct {
}
