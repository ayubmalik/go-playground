package main

import (
	"fmt"
	"net"
	"os"
)

const (
	Black   = "\033[0;30m"
	Red     = "\033[0;31m"
	Green   = "\033[0;32m"
	Yellow  = "\033[0;33m"
	Blue    = "\033[0;34m"
	Magenta = "\033[0;35m"
	Cyan    = "\033[0;36m"
	White   = "\033[0;37m"
	Default = "\033[0;39m"
	Reset   = "\033[0m"
)

func main() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("could not list interfaces, %v", err)
		os.Exit(1)
	}

	for _, i := range ifaces {
		if i.Flags&net.FlagPointToPoint == net.FlagPointToPoint {
			fmt.Println(Green + "pointtopoint connection detected: " + i.Name + Reset)
			return
		}
	}

	fmt.Println(Red + "pointtopoint connection not detected" + Reset)
	os.Exit(1)

}
