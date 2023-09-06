package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

const listenAddr = "localhost:4000"

func main() {
	http.Handle("/", websocket.Handler(handler))
	http.ListenAndServe(listenAddr, nil)
}

func handler(conn *websocket.Conn) {
	var s string
	fmt.Fscan(conn, &s)
	fmt.Println("Received: ", s)
	fmt.Fprint(conn, "How do you do?")
}
