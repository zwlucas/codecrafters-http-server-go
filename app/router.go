package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Router struct {
	routes []Route
}

func (r *Router) handler(request *Request) func(*Request) *Response {
	for _, route := range r.routes {
		if route.method == request.Method() && route.path.MatchString(request.Target()) {
			return route.handler
		}
	}

	return func(request *Request) *Response {
		return NewResponse(NotFound, nil, nil)
	}
}

func (r *Router) write(writer io.Writer, response *Response, request *Request) (err error) {
	if slices.Contains(strings.Split(request.Header("Accept-Encoding"), ", "), "gzip") {
		var b bytes.Buffer

		gz := gzip.NewWriter(&b)
		if _, err := gz.Write(response.body); err != nil {
			panic(err)
		}

		if err := gz.Close(); err != nil {
			panic(err)
		}

		response.headers["Content-Length"] = strconv.Itoa(b.Len())
		response.headers["Content-Encoding"] = "gzip"
		response.body = b.Bytes()
	}

	_, err = writer.Write(response.StatusLine())
	if err != nil {
		fmt.Println("Error writing status:", err)
		return
	}

	_, err = writer.Write(response.Headers())
	if err != nil {
		fmt.Println("Error writing headers:", err)
		return
	}

	_, err = writer.Write(response.Body())
	if err != nil {
		fmt.Println("Error writing body:", err)
		return
	}

	return
}

func (r *Router) Handle(conn net.Conn) error {
	defer conn.Close()

	request, err := NewRequest(bufio.NewReader(conn))
	if err != nil {
		return err
	}

	return r.write(conn, r.handler(request)(request), request)
}

type Route struct {
	method  string
	path    *regexp.Regexp
	handler func(*Request) *Response
}

type RouterBuilder struct {
	routes []Route
}

func (r *RouterBuilder) Add(method string, path *regexp.Regexp, handler func(*Request) *Response) {
	r.routes = append(r.routes, Route{
		method:  method,
		path:    path,
		handler: handler,
	})
}

func (r *RouterBuilder) Build() *Router {
	return &Router{
		routes: r.routes,
	}
}

func NewRouterBuilder() *RouterBuilder {
	return &RouterBuilder{
		routes: make([]Route, 0),
	}
}
