package server

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"http-server/server/core"
	c "http-server/server/core"
	"strconv"
)

const httpResponseStringFromat string = "HTTP/1.1 %d %s\r\n%s\r\n%s"

type HttpResponse struct {
	statusCode c.StatusCode
	headers    c.HttpRequestHeaders
	body       string
}

func (h *HttpResponse) Headers() c.HttpRequestHeaders {
	return h.headers
}

func NewInternalServerError() *HttpResponse {
	return &HttpResponse{
		statusCode: 500,
		headers:    c.NewEmpty(),
	}
}

func NewHttpResponse(statusCode c.StatusCode) *HttpResponse {
	return &HttpResponse{
		statusCode: statusCode,
		headers:    c.NewEmpty(),
	}
}

func NewNotFound() *HttpResponse {
	return &HttpResponse{
		statusCode: c.NotFound,
		headers:    c.NewEmpty(),
	}
}

func NewPermanentRedirect(location string) *HttpResponse {
	return &HttpResponse{
		statusCode: 308,
		headers: core.HttpRequestHeaders{
			"Content-Length": "0",
			"Location":       location,
		},
	}
}

func (h *HttpResponse) SetJsonBody(body any) error {
	h.headers.SetHeaderValue("Content-Type", "application/json")
	json, err := json.Marshal(body)
	if err != nil {
		return err
	}
	jsonString := string(json)
	h.headers.SetHeaderValue("Content-Length", strconv.Itoa(len(jsonString)))
	h.body = jsonString

	return nil
}

func (h *HttpResponse) SetTextBody(body string) {
	h.headers.SetHeaderValue("Content-Type", "text/plain")
	h.headers.SetHeaderValue("Content-Length", strconv.Itoa(len(body)))
	h.body = body
}

func (h *HttpResponse) SetOctetStreamBody(body []byte) {
	value := string(body)
	h.headers.SetHeaderValue("Content-Type", "application/octet-stream")
	h.headers.SetHeaderValue("Content-Length", strconv.Itoa(len(value)))
	h.body = value
}

func (h HttpResponse) ToHttpString() string {

	body := h.body

	encoding, hasEncoding := h.headers.GetHeaderValue("Content-Encoding")

	if hasEncoding && encoding == "gzip" {
		body = encodeToGzip(h.body)
		h.headers.SetHeaderValue("Content-Length", strconv.Itoa(len(body)))
	}

	return fmt.Sprintf(httpResponseStringFromat, h.statusCode, h.statusCode.GetReason(), formatResponseHeaders(h.headers), body)
}

func formatResponseHeaders(headers c.HttpRequestHeaders) string {

	if len(headers) == 0 {
		return ""
	}

	result := ""
	for k, v := range headers {
		result += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	return result
}

func encodeToGzip(data string) string {
	var bytesBuffer bytes.Buffer
	gzipper := gzip.NewWriter(&bytesBuffer)
	gzipper.Write([]byte(data))
	gzipper.Close()

	return bytesBuffer.String()
}
