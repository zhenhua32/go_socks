package main

import (
	"fmt"
	"net"
)

func main() {
	ip := net.ParseIP("127.0.0.1")
	ip = ip.To4()
	for i, v := range ip {
		fmt.Println(i, v)
	}
	fmt.Println("------")
	ip = ip.To16()
	for i, v := range ip {
		fmt.Println(i, v)
	}
}
