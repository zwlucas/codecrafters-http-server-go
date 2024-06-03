package main

import (
	"bufio"
	"io"
	"net"
	"regexp"
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

func (r *Router) write(writer io.Writer, response *Response) (err error) {
	_, err = writer.Write(response.StatusLine())
	if err != nil {
		return
	}

	_, err = writer.Write(response.Headers())
	if err != nil {
		return
	}

	_, err = writer.Write(response.Body())
	if err != nil {
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

	return r.write(conn, r.handler(request)(request))
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
