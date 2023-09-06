package main

import (
	"fmt"
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
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go match(conn)
	}
}

var partner = make(chan io.ReadWriteCloser)

// The match function simultaneously tries to send and receive a connection on a channel.
//   - If the send succeeds, the connection has been handed off to another goroutine,
//     so the function exits and the goroutine shuts down.
//   - If the receive succeeds, a connection has been received from another goroutine.
//     The current goroutine then has two connections, so it starts a chat session between them.
func match(c io.ReadWriteCloser) {
	fmt.Fprint(c, "Waiting for a partner...")
	select {
	case partner <- c:
		// now handled by the other goroutine
	case p := <-partner:
		chat(p, c)
	}
}

// The chat function sends a greeting to each connection
// and then copies data from one to the other, and vice versa.
// Notice that it launches another goroutine so that the copy operations may happen concurrently.
//func chat(a, b io.ReadWriteCloser) {
//	fmt.Fprintln(a, "Found one! Say hi.")
//	fmt.Fprintln(b, "Found one! Say hi.")
//	go io.Copy(a, b)
//	io.Copy(b, a)
//}

// It's important to clean up when the conversation is over.
// To do this we send the error value from each io.Copy call to the channel,
// log any non-nil errors, and close both connections.
func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "Found one! Say hi.")
	fmt.Fprintln(b, "Found one! Say hi.")
	errc := make(chan error, 1)
	go cp(a, b, errc)
	go cp(b, a, errc)
	if err := <-errc; err != nil {
		log.Println(err)
	}
	a.Close()
	b.Close()
}

func cp(w io.Writer, r io.Reader, errc chan<- error) {
	_, err := io.Copy(w, r)
	errc <- err
}
