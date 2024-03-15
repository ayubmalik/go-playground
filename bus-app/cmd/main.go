package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("Checking for a VPN")
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range ifaces {
		if i.Flags&net.FlagPointToPoint != 0 {
			if i.Flags&net.FlagRunning != 0 {
				log.Printf("looks like %s is a running VPN connection", i.Name)
			}
		}
	}
}
