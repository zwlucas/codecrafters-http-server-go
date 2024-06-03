package main

import (
	"bufio"
	"bytes"
	"fmt"
)

type requestHeaders struct {
	headers map[string]string
}

type requestLine struct {
	method  string
	target  string
	version string
}

type Request struct {
	line    requestLine
	headers requestHeaders
	body    *bufio.Reader
}

func newRequestHeaders(reader *bufio.Reader) (requestHeaders, error) {
	request := requestHeaders{make(map[string]string)}

	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			return request, err
		}

		if bytes.Equal(buf, []byte{'\r', '\n'}) {
			return request, nil
		}

		buf = bytes.TrimSuffix(buf, []byte{'\r', '\n'})
		header := bytes.Split(buf, []byte{':', ' '})

		if len(header) != 2 {
			return request, fmt.Errorf("invalid header line")
		}

		request.headers[string(header[0])] = string(header[1])
	}
}

func newRequestLine(reader *bufio.Reader) (requestLine, error) {
	buf, err := reader.ReadBytes('\n')
	if err != nil {
		return requestLine{}, err
	}

	buf = bytes.TrimSuffix(buf, []byte{'\r', '\n'})
	elements := bytes.Split(buf, []byte{' '})

	if len(elements) != 3 {
		return requestLine{}, fmt.Errorf("invalid request, expected 3 elements, got %d", len(elements))
	}

	return requestLine{string(elements[0]), string(elements[1]), string(elements[2])}, nil
}

func NewRequest(reader *bufio.Reader) (*Request, error) {
	line, err := newRequestLine(reader)
	if err != nil {
		return nil, err
	}

	headers, err := newRequestHeaders(reader)
	if err != nil {
		return nil, err
	}

	return &Request{line: line, headers: headers, body: reader}, nil
}

func (r *Request) Body() *bufio.Reader {
	return r.body
}

func (r *Request) Method() string {
	return r.line.method
}

func (r *Request) Target() string {
	return r.line.target
}

func (r *Request) Header(name string) string {
	return r.headers.headers[name]
}
