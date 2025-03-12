package server

import (
	"fmt"
	"http-server/server/builder"
	"http-server/server/routing"
	"net"
)

type routeWithHandler struct {
	verb    Verb
	route   routing.HttpRouteTemplate
	handler HttpHandler
}

type Server struct {
	routes  []routeWithHandler
	options builder.Options
}

type HttpHandler func(request HttpRequest) HttpResponse

func NewServer(options ...builder.Option) (*Server, error) {
	var serverOptions builder.Options

	for _, option := range options {
		err := option(&serverOptions)
		if err != nil {
			return nil, err
		}
	}

	return &Server{
		routes:  make([]routeWithHandler, 0),
		options: serverOptions,
	}, nil
}

func (s *Server) Run() error {
	const addressFormat string = "%s:%d"
	serverAddress := fmt.Sprintf(addressFormat, s.options.Host, s.options.Port)
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		return fmt.Errorf("failed to bind to port %s", serverAddress)
	}
	fmt.Printf("Listening on port %s\n", serverAddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Unable to accept request")
		}

		go s.processTcpRequest(conn)
	}
}
