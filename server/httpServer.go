package server

import (
	"errors"
	"fmt"
	"http-server/server/routing"
	"net"
	"strings"
)

type routeWithHandler struct {
	verb    Verb
	route   routing.HttpRouteTemplate
	handler HttpHandler
}

type Server struct {
	routes []routeWithHandler
}

type HttpHandler func(request HttpRequest) HttpResponse

func NewServer() *Server {
	return &Server{
		routes: make([]routeWithHandler, 0),
	}
}

func (s *Server) MapGet(path string, handler HttpHandler) error {

	route, err := routing.NewRouteFromString(path)

	if err != nil {
		return errors.Join(errors.New("unable to register the route"), err)
	}

	s.routes = append(s.routes, routeWithHandler{
		verb:    GET,
		route:   route,
		handler: handler,
	})

	return nil
}

func (s *Server) MapPost(path string, handler HttpHandler) error {
	route, err := routing.NewRouteFromString(path)

	if err != nil {
		return errors.Join(errors.New("unable to register the route"), err)
	}

	s.routes = append(s.routes, routeWithHandler{
		verb:    POST,
		route:   route,
		handler: handler,
	})

	return nil
}

func (s *Server) RunOnPort(port string) error {
	const addressFormat string = "0.0.0.0:%s"
	listener, err := net.Listen("tcp", fmt.Sprintf(addressFormat, port))
	if err != nil {
		return fmt.Errorf("failed to bind to port %s", port)
	}
	fmt.Printf("Listening on port %s\n", port)

	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Unable to accept request")
		}

		go func() {

			request, err := NewFromTCPConnection(conn)
			if err != nil {
				fmt.Println(err)
				return
			}

			s.executeRequest(request.Verb, conn, request)
		}()
	}
}

func writeInternalServerErrorOnPanic(c net.Conn) {
	if r := recover(); r != nil {
		c.Write([]byte(NewInternalServerError().ToHttpString()))
	}
}

func (s *Server) executeRequest(verb Verb, conn net.Conn, request *HttpRequest) {
	defer writeInternalServerErrorOnPanic(conn)
	defer conn.Close()

	handler := s.findHandler(verb, request)

	if handler == nil {
		handler = notFound
	}

	response := handler(*request)

	encoding, hasHeader := request.Headers.GetHeaderValue("Accept-Encoding")
	if hasHeader && strings.Contains(encoding, "gzip") && request.Url.Path != "/" {
		response.Headers().SetHeaderValue("Content-Encoding", "gzip")
	}
	conn.Write([]byte(response.ToHttpString()))
}

func (s *Server) findHandler(verb Verb, request *HttpRequest) HttpHandler {
	for _, routeWithHandler := range s.routes {

		if routeWithHandler.verb != verb {
			continue
		}

		routeParameters, canHandle := routeWithHandler.route.CanHandlerPath(request.Url.Path)
		if canHandle {
			request.RouteParameters = routeParameters
			return routeWithHandler.handler
		}
	}

	return notFound
}

func notFound(request HttpRequest) HttpResponse {
	return *NewNotFound()
}
