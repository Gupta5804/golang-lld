package main

type Server struct {
	host    string
	port    int
	timeout int
	maxconn int
}

type Option func(*Server)

func WithPort(port int) Option{
	return func(s *Server){
		s.port = port
	}
}
func WithTimeout(timeout int) Option{
	return func(s *Server){
		s.timeout = timeout
	}
}

func NewServer(opts ...Option) *Server {
	svr := &Server{
		host:    "localhost",
		port:    8080,
		timeout: 30,
		maxconn: 100,
	}

	for _,opt := range opts{
		opt(svr)
	}
	return svr
}
