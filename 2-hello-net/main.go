package main

import (
	"fmt"
	"log"
	"net"
)

const listenAddr = "localhost:4000"

func main() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Use Fprintln to write to a net connection.
		// It writes to an io.Writer, and net.Conn is an io.Writer.
		_, err = fmt.Fprintln(conn, "Hello!")
		if err != nil {
			log.Fatal(err)
		}
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}
