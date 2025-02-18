package client

type InternalClient interface {
	BaseClient
	ChangeConn(targetAddress string)
}

func NewInternalClient(cfg Config) InternalClient {
	return &SimpleClient{Config: cfg}
}

func (c *SimpleClient) ChangeConn(targetAddress string) {
	c.Config.Address = targetAddress
}
