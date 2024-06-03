package main

import (
	"regexp"
	"strconv"
	"strings"
)

func basePathHandler(request *Request) *Response {
	return NewResponse(Ok, nil, nil)
}

func echoHandler(request *Request) *Response {
	message := strings.Split(request.Target(), "/")[2]
	return NewResponse(Ok, map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": strconv.Itoa(len(message)),
	}, []byte(message))
}

func userAgentHandler(request *Request) *Response {
	header := request.Header("User-Agent")
	return NewResponse(Ok, map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": strconv.Itoa(len(header)),
	}, []byte(header))
}

func Register(builder *RouterBuilder) {
	builder.Add("GET", regexp.MustCompile("^/$"), basePathHandler)
	builder.Add("GET", regexp.MustCompile("^/echo/[a-zA-Z]+$"), echoHandler)
	builder.Add("GET", regexp.MustCompile("^/user-agent$"), userAgentHandler)
}
