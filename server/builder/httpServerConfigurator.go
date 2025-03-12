package builder

type Options struct {
	Port int
	Host string
}

type Option func(options *Options) error

func WithPort(port int) Option {
	return func(options *Options) error {
		options.Port = port
		return nil
	}
}

func WithHost(host string) Option {
	return func(options *Options) error {
		options.Host = host
		return nil
	}
}
