package main

import (
	"fmt"
	"net"
)

func main() {
	dialApple_Public()
	dialWiktor_Private()
}

// Wiktor runs
func dialApple_Public() {
	appleNodeAddr := "128.1.1.1:333"
	conn, err := net.Dial("tcp", appleNodeAddr)
	fmt.Println("err:", err)

	// go run .
	// err: <nil>
}

// Apple runs
func dialWiktor_Private() {
	wiktorNodeAddr := "224.1.1.1:222"
	conn, err := net.Dial("tcp", wiktorNodeAddr)
	fmt.Println("err:", err)

	// go run .
	// err: dial tcp 224.1.1.1:222: connect: operation timed out
}
