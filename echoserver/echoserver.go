package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// Example echo server to use as a test service
func main() {
	var port string
	flag.StringVar(&port, "port", "", "Port to listen on")
	flag.Parse()

	if port == "" {
		fmt.Printf("Port required\n")
		flag.Usage()
		os.Exit(1)
	}

	// create a tcp listener on the given port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("failed to create listener, err:", err)
		os.Exit(1)
	}

	// listen for new connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("failed to accept connection, err:", err)
			continue
		}

		// pass an accepted connection to a handler goroutine
		go handleConnection(conn)
	}
}

// handleConnection handles the lifetime of a connection
func handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	io.Copy(conn, conn)
}
