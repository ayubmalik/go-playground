package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	var creds string
	flag.StringVar(&creds, "creds", "", "creds file location")
	flag.Parse()

	if creds == "" {
		log.Fatalln("-creds flag needs to be set")
	}

	nc, err := nats.Connect("tls://connect.ngs.global", nats.UserCredentials(creds), nats.Name("hellosir.pub"))
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

	id, _ := nc.GetClientID()
	ip, _ := nc.GetClientIP()

	log.Printf("client id: %v, ip: %v\n", id, ip)
	log.Println("press ctrl-c to exit")

	err = Pub(ctx, nc, "tds.hello")
	log.Println(err)
	<-ctx.Done()
}

func Pub(ctx context.Context, nc *nats.Conn, subject string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			data := "Hello " + time.Now().String()
			err := nc.Publish(subject, []byte(data))
			if err != nil {
				return err
			}
		}
		time.Sleep(1 * time.Second)
	}

}
