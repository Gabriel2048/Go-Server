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

		if c, ok := err.(nonHttpsRequestError); ok {
			s.executeNonHttpsRequest(c.conn)
			return
		}

		return
	}

	s.executeRequest(conn, request)
}

func (s *Server) executeNonHttpsRequest(conn net.Conn) {
	if s.options.HttpsRedirect {
		const httpsRedirectUrlFormat = "https://%s/%s"
		redirectUrl := fmt.Sprintf(httpsRedirectUrlFormat, s.Addr(), "echo/cvcx")
		r := NewPermanentRedirect(redirectUrl)
		println(r.ToHttpString())
		conn.Write([]byte(r.ToHttpString()))
		conn.Close()
	} else {

	}
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
