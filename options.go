package acrouter

type Option func(*RouterConfig)

func Address(ip string, port int) Option {
	return func(s *RouterConfig) {
		s.Addr = ip
		s.Port = port
	}
}

func DebugMode(debug bool) Option {
	return func(s *RouterConfig) {
		s.DebugMode = debug
	}
}

func Tls(key, cert string) Option {
	return func(s *RouterConfig) {
		s.Key = key
		s.Cert = cert
	}
}
