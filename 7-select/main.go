package main

import (
	"fmt"
	"time"
)

// A select statement is like a switch,
// but it selects over channel operations (and chooses exactly one of them).
func main() {
	ticker := time.NewTicker(250 * time.Millisecond)
	boom := time.After(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			fmt.Println("tick")
		case <-boom:
			fmt.Println("boom!")
			return
		}
	}
}
