package core

import (
	"errors"
	"strconv"
	"strings"
)

type HttpRequestHeaders map[string]string

func NewEmpty() HttpRequestHeaders {
	return make(map[string]string)
}

func (h HttpRequestHeaders) SetHeaderValue(key string, value string) {
	h[strings.ToLower(key)] = value
}

func (h HttpRequestHeaders) GetHeaderValue(headerKey string) (string, bool) {
	headerValue, hasHeader := h[strings.ToLower(headerKey)]

	return headerValue, hasHeader
}

func (h HttpRequestHeaders) GetContentLength() (int, error) {
	headerValue, hasHeader := h["content-length"]

	if !hasHeader {
		return -1, errors.New("requested header does not exists")
	}

	contentLength, err := strconv.Atoi(headerValue)

	if err != nil {
		return -1, errors.New("content-length is not a number")
	}

	return contentLength, nil
}
