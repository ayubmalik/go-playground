package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/nats-io/nats.go"
)

func main() {
	var creds string
	flag.StringVar(&creds, "creds", "", "creds file location")
	flag.Parse()

	if creds == "" {
		log.Fatalln("-creds flag needs to be set")
	}

	nc, err := nats.Connect("tls://connect.ngs.global", nats.UserCredentials(creds), nats.Name("hellosir.sub"))
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer func() {
		nc.Close()
		cancel()
		stop()
	}()

	log.Println("press ctrl-c to exit")
	err = Sub(ctx, nc, "tds.hello")
	log.Println(err)
	<-ctx.Done()
}

func Sub(ctx context.Context, nc *nats.Conn, subject string) error {
	sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
		_ = msg.Ack()
		log.Printf("msg received: %v", string(msg.Data))
	})
	if err != nil {
		return err
	}

	<-ctx.Done()
	return sub.Drain()
}
