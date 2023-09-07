package main

import "fmt"

type A struct {
}

func (A) Hello() {
	fmt.Println("Hello!")
}

type B struct {
	A
}

// Implicitly
//func (b B) Hello() {
//	b.A.Hello()
//}

// Go supports a kind of "mix-in" functionality with a feature known as "struct embedding".
// The embedding struct delegates calls to the embedded type's methods.
func main() {
	var b B
	b.Hello()
}
