package main

import (
	"flag"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	var creds string
	flag.StringVar(&creds, "creds", "", "creds file location")
	flag.Parse()

	if creds == "" {
		log.Fatalln("-creds flag needs to be set")
	}

	nc, err := nats.Connect("tls://connect.ngs.global", nats.UserCredentials(creds))
	if err != nil {
		log.Fatalln(err)
	}

	defer nc.Close()

	id, _ := nc.GetClientID()
	ip, _ := nc.GetClientIP()

	log.Printf("Client ID: %v, IP: %v", id, ip)
}
