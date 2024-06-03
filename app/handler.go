package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var directory = ""

func init() {
	args := os.Args
	index := slices.Index(args, "--directory")
	if index > 0 {
		directory = args[index+1]
	}

	fmt.Println("directory: ", directory)
}

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

func getFile(request *Request) *Response {
	file, err := os.Open(filepath.Join(directory, strings.Split(request.Target(), "/")[2]))
	if err != nil {
		return NewResponse(NotFound, nil, nil)
	}

	buf, _ := io.ReadAll(file)
	return NewResponse(Ok, map[string]string{
		"Content-Type":   "application/octet-stream",
		"Content-Length": strconv.Itoa(len(buf)),
	}, buf)
}

func createFile(request *Request) *Response {
	fPath := filepath.Join(directory, strings.Split(request.Target(), "/")[2])
	fmt.Println("create file: ", fPath)

	file, err := os.Create(fPath)
	if err != nil {
		return NewResponse(InternalServerError, nil, nil)
	}

	defer file.Close()

	length, err := strconv.Atoi(request.Header("Content-Length"))
	if err != nil {
		return NewResponse(InternalServerError, nil, nil)
	}

	_, err = io.CopyN(file, request.Body(), int64(length))
	if err != nil {
		return NewResponse(InternalServerError, nil, nil)
	}

	return NewResponse(Created, nil, nil)
}

func Register(builder *RouterBuilder) {
	builder.Add("GET", regexp.MustCompile("^/$"), basePathHandler)
	builder.Add("GET", regexp.MustCompile("^/echo/[a-zA-Z]+$"), echoHandler)
	builder.Add("GET", regexp.MustCompile("^/user-agent$"), userAgentHandler)

	if directory != "" {
		builder.Add("GET", regexp.MustCompile("^/files/.+$"), getFile)
		builder.Add("POST", regexp.MustCompile("^/files/.+$"), createFile)
	}
}
