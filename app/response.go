package main

import "fmt"

type Response struct {
	status  Status
	headers map[string]string
	body    []byte
}

func (r *Response) StatusLine() []byte {
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.status.Code(), r.status.Text()))
}

func (r *Response) Headers() []byte {
	buf := make([]byte, 0)

	for key, value := range r.headers {
		buf = append(buf, fmt.Sprintf("%s: %s\r\n", key, value)...)
	}

	return append(buf, "\r\n"...)
}

func (r *Response) Body() []byte {
	return r.body
}

func NewResponse(code Status, headers map[string]string, body []byte) *Response {
	return &Response{code, headers, body}
}
