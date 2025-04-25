package builder

import "crypto/tls"

type Options struct {
	Port      int
	Host      string
	TlsConfig *tls.Config
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

func WithTLSConfig(tlsOptions *tls.Config) Option {
	return func(options *Options) error {
		options.TlsConfig = tlsOptions
		return nil
	}
}
