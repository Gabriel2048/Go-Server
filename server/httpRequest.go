package server

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	c "http-server/server/core"
	r "http-server/server/routing"
	"io"
	"net"
	"net/url"
	"strings"
)

type Verb string

const (
	GET    Verb = "GET"
	POST   Verb = "POST"
	PATCH  Verb = "PATCH"
	DELETE Verb = "DELETE"
)

type HttpRequest struct {
	Verb            Verb
	Url             *url.URL
	RouteParameters r.RouteParameters
	Headers         c.HttpRequestHeaders
	Body            io.Reader
	handler         HttpHandler
}

type httpStatusLine struct {
	verb        Verb
	target      *url.URL
	httpVersion string
}

func NewFromTCPConnection(conn net.Conn, routes []routeWithHandler) (*HttpRequest, error) {
	reader := bufio.NewReader(conn)

	statusLine, err := readStatusLine(reader)

	if err != nil {
		return nil, err
	}

	headers, err := readHeaders(reader)
	if err != nil {
		return nil, err
	}

	body := bufio.NewReader(reader)

	handler, routeParameters := matchHandlerAndRouteParams(*statusLine.target, statusLine.verb, routes)

	return &HttpRequest{
		Verb:            statusLine.verb,
		Url:             statusLine.target,
		Headers:         *headers,
		Body:            body,
		RouteParameters: routeParameters,
		handler:         handler,
	}, nil
}

// Reads the entire request body according to the Content-Length header
func (h *HttpRequest) ReadBody() ([]byte, error) {
	cotentLength, err := h.Headers.GetContentLength()

	if err != nil {
		return nil, errors.Join(errors.New("unable to read request body"), err)
	}

	body := make([]byte, cotentLength)

	readAmount, err := h.Body.Read(body)

	if err != nil {
		return nil, errors.Join(errors.New("unable to read request body"), err)
	}

	if readAmount != cotentLength {
		return nil, errors.New("unable to read request body due to body size to content-length missmatch")
	}

	return body, nil
}

var ErrMalformedStatusLine = errors.New("malformed status line")

func readStatusLine(reader *bufio.Reader) (*httpStatusLine, error) {
	rawRequestLine, err := reader.ReadString('\n')

	if err != nil {
		if r, isRecordHeaderError := err.(tls.RecordHeaderError); isRecordHeaderError {
			// t := r.Conn.(*net.TCPConn)

			reader := bufio.NewReader(r.Conn)
			println(reader.Size())
			// qq := make([]byte, reader.Size())
			// a, err := reader.Read(qq)

			return nil, nonHttpsRequestError{
				conn: r.Conn,
				path: r.Conn.RemoteAddr().String(),
			}
		} else {
			return nil, ErrMalformedStatusLine
		}
	}

	requestLine := rawRequestLine[0 : len(rawRequestLine)-2]
	parts := strings.Split(requestLine, " ")

	if len(parts) != 3 {
		return nil, ErrMalformedStatusLine
	}

	verb, err := castToVerb(parts[0])

	if err != nil {
		return nil, ErrMalformedStatusLine
	}

	url, err := url.Parse(parts[1])

	if err != nil {
		return nil, ErrMalformedStatusLine
	}

	return &httpStatusLine{
		verb:        verb,
		target:      url,
		httpVersion: parts[2],
	}, nil
}

func readHeaders(reader *bufio.Reader) (*c.HttpRequestHeaders, error) {
	const headerLineEnding string = "\r\n"
	const headerSeparator string = ": "

	result := c.NewEmpty()
	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if headerLine == headerLineEnding {
			break
		}
		if len(headerLine) > 1 {
			rawHeaderKeyValue := string(headerLine)[0 : len(headerLine)-2]
			headerParts := strings.Split(rawHeaderKeyValue, headerSeparator)

			headerKey := headerParts[0]
			headerValue := headerParts[1]
			result.SetHeaderValue(headerKey, headerValue)
		}
	}

	return &result, nil
}

func matchHandlerAndRouteParams(url url.URL, verb Verb, routes []routeWithHandler) (HttpHandler, r.RouteParameters) {
	for _, routeWithHandler := range routes {

		if routeWithHandler.verb != verb {
			continue
		}

		routeParameters, canHandle := routeWithHandler.route.CanHandlerPath(url.Path)
		if canHandle {
			return routeWithHandler.handler, routeParameters
		}
	}

	return notFound, r.RouteParameters{}
}

func castToVerb(s string) (Verb, error) {
	switch Verb(s) {
	case GET, POST, PATCH, DELETE:
		return Verb(s), nil
	default:
		return "", fmt.Errorf("invalid Verb: %s", s)
	}
}

func notFound(request HttpRequest) HttpResponse {
	return *NewNotFound()
}
