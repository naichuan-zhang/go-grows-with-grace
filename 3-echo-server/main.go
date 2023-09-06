package main

import (
	"io"
	"log"
	"net"
)

const listenAddr = "localhost:4000"

func main() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Println("Server is running on: ", listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		written, err := io.Copy(conn, conn)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Written: ", written)
	}
}
