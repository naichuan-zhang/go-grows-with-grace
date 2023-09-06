package main

import "fmt"

// Goroutines communicate via channels.
// A channel is a typed conduit that may be synchronous (unbuffered) or asynchronous (buffered).
func main() {
	ch := make(chan int)
	go fibs(ch)
	for i := 0; i < 20; i++ {
		fmt.Println(<-ch)
	}
}

func fibs(ch chan int) {
	i, j := 0, 1
	for {
		ch <- j
		i, j = j, i+j
		//println(i, j)
	}
}
