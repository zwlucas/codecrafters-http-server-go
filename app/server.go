package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
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

	request, err := NewRequest(conn)
	if err != nil {
		fmt.Println("Error getting target: ", err.Error())
		os.Exit(1)
	}

	if request.Target == "/" {
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		if err != nil {
			fmt.Println("Error connection write: ", err.Error())
		}
	} else if strings.HasPrefix(request.Target, "/echo/") {
		message := strings.SplitN(request.Target, "/", 3)[2]

		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			fmt.Println("Error connection write: ", err.Error())
		}

		_, err = conn.Write([]byte("Content-Type: text/plain\r\n"))
		if err != nil {
			fmt.Println("Error connection write: ", err.Error())
		}

		_, err = conn.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n", len(message))))
		if err != nil {
			fmt.Println("Error connection write: ", err.Error())
		}

		_, err = conn.Write([]byte("\r\n"))
		if err != nil {
			fmt.Println("Error connection write: ", err.Error())
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error connection write: ", err.Error())
		}
	} else {
		_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		if err != nil {
			fmt.Println("Error connection write: ", err.Error())
		}
	}
}
