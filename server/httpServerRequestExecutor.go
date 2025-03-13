package server

import (
	"fmt"
	"net"
	"strings"
)

func (s *Server) processTcpRequest(conn net.Conn) {
	defer writeInternalServerErrorOnPanic(conn)
	defer conn.Close()

	request, err := NewFromTCPConnection(conn, s.routes)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.executeRequest(conn, request)
}

func (s *Server) executeRequest(conn net.Conn, request *HttpRequest) {

	response := request.handler(*request)

	encoding, hasHeader := request.Headers.GetHeaderValue("Accept-Encoding")
	if hasHeader && strings.Contains(encoding, "gzip") && request.Url.Path != "/" {
		response.Headers().SetHeaderValue("Content-Encoding", "gzip")
	}
	conn.Write([]byte(response.ToHttpString()))
}

func writeInternalServerErrorOnPanic(c net.Conn) {
	if r := recover(); r != nil {
		c.Write([]byte(NewInternalServerError().ToHttpString()))
	}
}
