package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	var credentials string
	flag.StringVar(&credentials, "credentials", "", "credentials file location")
	flag.Parse()

	if credentials == "" {
		log.Fatalln("-credentials flag needs to be set")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	opts := []nats.Option{
		nats.UserCredentials(credentials),
		nats.Name("sub.schedules.candidates"),
		nats.ConnectHandler(func(c *nats.Conn) {
			log.Println("Connected to", c.ConnectedUrl())
		}),
		nats.ReconnectWait(3 * time.Second),
		nats.ReconnectHandler(func(c *nats.Conn) {
			log.Println("Reconnected to", c.ConnectedUrl())
		}),
		nats.DisconnectErrHandler(func(conn *nats.Conn, err error) {
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

	finder := NewClient()
	sub, err := nc.QueueSubscribe("tds.schedules.candidates", "candidates", func(msg *nats.Msg) {
		log.Printf("- %s - got msg: %s", msg.Header.Get(nats.MsgIdHdr), string(msg.Data))
		stops := strings.Split(string(msg.Data), " ")
		scheduleErr := trySchedule(finder, 7, stops[0], stops[1])
		if scheduleErr != nil {
			log.Println(scheduleErr)
		}
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
	FindSchedules(ctx context.Context, qry ScheduleQuery) (ScheduleResult, error)
}

type OriginDestinationInserter interface {
}

func trySchedule(sf ScheduleFinder, days int, origin, dest string) error {
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := range days {
		wg.Add(1)
		departure := time.Now().Add(time.Duration(i+1) * 24 * time.Hour).Format("2006-01-02")
		qry := newScheduleQuery(origin, dest, departure)

		go func() {
			defer wg.Done()
			result, err := sf.FindSchedules(ctx, qry)
			if err != nil && !errors.Is(err, context.Canceled) {
				log.Printf("error finding schedules: %s", err)
				return
			}

			if !result.IsEmpty() {
				log.Printf("Found a valid schedule for %s %s %s", departure, origin, dest)
				cancel()
			}
		}()
		time.Sleep(50 * time.Millisecond)
	}
	wg.Wait()
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
