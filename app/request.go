package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Request struct {
	Method  string
	Target  string
	Version string
}

func NewRequest(reader io.Reader) (*Request, error) {
	scanner := bufio.NewScanner(reader)

	ok := scanner.Scan()
	if !ok {
		err := scanner.Err()
		if err == nil {
			err = io.EOF
		}

		return nil, err
	}

	line := scanner.Bytes()
	elements := bytes.Split(line, []byte(" "))

	if len(elements) != 3 {
		return nil, fmt.Errorf("invalid request, excepted 3 elements, got %d", len(elements))
	}

	return &Request{Method: string(elements[0]), Target: string(elements[1]), Version: string(elements[2])}, nil
}
