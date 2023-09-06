package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const dialAddr = "localhost:4000"

func main() {
	conn, err := net.Dial("tcp", dialAddr)
	if err != nil {
		log.Fatal(err)
		return
	}
	resp, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Response from server: ", resp)
}
