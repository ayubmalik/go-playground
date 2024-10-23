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

type ScheduleFinder interface {
	FindSchedules(ScheduleQuery) (ScheduleResult, error)
}

func trySchedule(sf ScheduleFinder, days int, origin, dest string) error {
	for i := range days {
		departure := time.Now().Add(time.Duration(i+1) * 24 * time.Hour).Format("2006-01-02")
		qry := newScheduleQuery(origin, dest, departure)

		result, err := sf.FindSchedules(qry)
		if err != nil {
			return err
		}

		if !result.IsEmpty() {
			log.Printf("Found a valid schedule for %s -> %s ", origin, dest)
			break
		}
	}
	return nil
}

func newScheduleQuery(origin, dest, departure string) ScheduleQuery {
	return ScheduleQuery{
		PurchaseType:    "SCHEDULE_BOOK",
		Origin:          Stop{origin},
		Destination:     Stop{dest},
		DepartDate:      departure,
		PassengerCounts: map[string]int{"Adult": 1},
	}
}
