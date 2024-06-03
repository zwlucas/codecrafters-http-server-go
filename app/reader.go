package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type RequestReader struct {
	scanner *bufio.Scanner
	lineBuf []byte
}

func (r *RequestReader) read() ([]byte, error) {
	ok := r.scanner.Scan()
	if !ok {
		err := r.scanner.Err()
		if err == nil {
			err = io.EOF
		}

		return nil, err
	}

	return r.scanner.Bytes(), nil
}

func (r *RequestReader) line() ([]byte, error) {
	if r.lineBuf == nil {
		lineBuf, err := r.read()
		if err != nil {
			return nil, err
		}

		r.lineBuf = lineBuf
	}

	return r.lineBuf, nil
}

func (r *RequestReader) element(index int) ([]byte, error) {
	line, err := r.line()
	if err != nil {
		return nil, err
	}

	return bytes.Split(line, []byte(" "))[index], nil
}

func (r *RequestReader) Method() (string, error) {
	fmt.Println("Reading method...")
	element, err := r.element(0)
	if err != nil {
		return "", err
	}

	return string(element), nil
}

func (r *RequestReader) Target() (string, error) {
	fmt.Println("Reading target...")
	element, err := r.element(1)
	if err != nil {
		return "", err
	}

	return string(element), nil
}

func NewRequestReader(reader io.Reader) *RequestReader {
	return &RequestReader{
		scanner: bufio.NewScanner(reader),
	}
}
