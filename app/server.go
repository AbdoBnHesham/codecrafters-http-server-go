package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		panic("Failed to bind to port 4221")
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			message := fmt.Errorf("error accepting connection: %v", err.Error())
			panic(message)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) error {
	defer conn.Close()

	_, err := conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	if err != nil {
		return fmt.Errorf("error writing response: %v", err.Error())
	}

	return nil
}
