package server

import (
	"fmt"
	"net"
	"strings"
)

func (s *Server) processTcpRequest(conn net.Conn) {
	defer writeInternalServerErrorOnPanic(conn)
	defer conn.Close()

	request, err := NewFromTCPConnection(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.executeRequest(request.Verb, conn, request)
}

func (s *Server) executeRequest(verb Verb, conn net.Conn, request *HttpRequest) {
	handler := s.findHandler(verb, request)

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

func writeInternalServerErrorOnPanic(c net.Conn) {
	if r := recover(); r != nil {
		c.Write([]byte(NewInternalServerError().ToHttpString()))
	}
}

func notFound(request HttpRequest) HttpResponse {
	return *NewNotFound()
}
