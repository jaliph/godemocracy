package democracy

type Option func(*NodeConfig)

func WithSource(source string) Option {
	return func(nc *NodeConfig) {
		nc.Source = source
	}
}

func WithTimeout(timeout int) Option {
	return func(nc *NodeConfig) {
		nc.Timeout = timeout
	}
}

func WithInterval(interval int) Option {
	return func(nc *NodeConfig) {
		nc.Interval = interval
	}
}

func WithPeers(peers []string) Option {
	return func(nc *NodeConfig) {
		nc.Peers = peers
	}
}

func WithMaxPacketSize(size int) Option {
	return func(nc *NodeConfig) {
		nc.MaxPacketSize = size
	}
}

func WithId(id string) Option {
	return func(nc *NodeConfig) {
		nc.Id = id
	}
}

func WithChannels(channels []string) Option {
	return func(nc *NodeConfig) {
		nc.Channels = channels
	}
}

func NewConfig(opts ...Option) *NodeConfig {
	nc := &NodeConfig{
		Interval:      1000,
		Timeout:       3000,
		Source:        "127.0.0.1:12345",
		MaxPacketSize: 2,
		Peers:         []string{},
		Channels:      []string{},
		Weight:        RandGenerator(),
	}
	nc.Id, _ = GenerateShortID()

	for _, opt := range opts {
		opt(nc)
	}

	return nc
}
