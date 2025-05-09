package server

import "net"

type nonHttpsRequestError struct {
	conn net.Conn
	path string
}

func (n nonHttpsRequestError) Error() string {
	return "request made using http instead of https"
}
