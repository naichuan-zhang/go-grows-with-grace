package main

import (
	"fmt"
	"time"
)

// Goroutines are lightweight threads that are managed by the Go runtime.
// To run a function in a new goroutine, just put "go" before the function call.
func main() {
	go say("let's go!", 3)
	go say("ho!", 2)
	go say("hey!", 1)
	time.Sleep(4 * time.Second)
}

func say(text string, secs int) {
	time.Sleep(time.Duration(secs) * time.Second)
	fmt.Println(text)
}
