package main

import (
	"flag"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

func main() {
	var creds string
	flag.StringVar(&creds, "creds", "", "creds file location")
	flag.Parse()

	if creds == "" {
		log.Fatalln("-creds flag needs to be set")
	}

	nc, err := nats.Connect("tls://connect.ngs.global", nats.UserCredentials(creds), nats.Name("pub.schedules.candidates"))
	if err != nil {
		log.Fatalf("could not connect to nats %s\n", err)
	}

	defer func() {
		_ = nc.Drain()
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	sub, err := nc.QueueSubscribe("tds.schedules.candidates", "candidates", func(msg *nats.Msg) {
		log.Printf("- %s - got msg: %s", msg.Header.Get(nats.MsgIdHdr), string(msg.Data))
	})
	if err != nil {
		wg.Done()
		log.Fatal(err)
	}
	wg.Wait()
	_ = sub.Drain()

}
