package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	opts := []nats.Option{
		nats.UserCredentials(creds),
		nats.Name("sub.schedules.candidates"),
		nats.ConnectHandler(func(c *nats.Conn) {
			log.Println("Connected to", c.ConnectedUrl())
		}),
		nats.ReconnectWait(3 * time.Second),
		nats.ReconnectHandler(func(c *nats.Conn) {
			log.Println("Reconnected to", c.ConnectedUrl())
		}),
		nats.DisconnectHandler(func(c *nats.Conn) {
			log.Println("Disconnected from NATS")
		}),
		nats.ClosedHandler(func(c *nats.Conn) {
			log.Println("NATS connection is closed.")
		}),
		nats.NoReconnect(),
	}

	nc, err := nats.Connect("tls://connect.ngs.global", opts...)
	if err != nil {
		log.Fatalf("could not connect to nats %s\n", err)
	}

	defer func() {
		_ = nc.Drain()
	}()

	sub, err := nc.QueueSubscribe("tds.schedules.candidates", "candidates", func(msg *nats.Msg) {
		log.Printf("- %s - got msg: %s", msg.Header.Get(nats.MsgIdHdr), string(msg.Data))
	})
	if err != nil {
		stop()
		log.Fatal(err)
	}

	<-ctx.Done()
	log.Println("draining sub and conn")
	_ = sub.Drain()

}

func trySchedule(schedules ScheduleAPI, days int, origin, dest string) error {
	var dates []string
	for i := range days {
		dt := time.Now().Add(time.Duration(i+1) * 24 * time.Hour)
		dates = append(dates, dt.Format("2006-01-02"))
	}
	log.Printf("dates: %v", dates)
	return nil
}

func newScheduleRequest(origin, dest, departure string) ScheduleRequest {
	return ScheduleRequest{
		PurchaseType:    "SCHEDULE_BOOK",
		Origin:          Stop{origin},
		Destination:     Stop{dest},
		DepartDate:      departure,
		PassengerCounts: map[string]int{"Adult": 1},
	}
}
