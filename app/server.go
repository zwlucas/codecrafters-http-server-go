package main

import (
	"fmt"
	"net"
	"os"
)

func createRouter() *Router {
	builder := NewRouterBuilder()
	Register(builder)
	return builder.Build()
}

func main() {
	router := createRouter()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	router.Handle(conn)
}
