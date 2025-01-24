package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

func main() {
	after := 15
	if run := os.Getenv("RUN_FOR"); run != "" {
		i, err := strconv.Atoi(run)
		if err != nil {
			slog.Error("could not parse RUN_FOR", "error", err)
			os.Exit(1)
		}
		after = i
	}

	done := make(chan struct{})
	time.AfterFunc(time.Duration(after)*time.Second, func() {
		close(done)
	})

	msg := "world"
	if m := os.Getenv("MSG"); m != "" {
		msg = m
	}

	ticker := time.Tick(2 * time.Second)
	go func() {
		for next := range ticker {
			slog.Info(fmt.Sprintf("Hola %s", msg), "time", next.Format(time.RFC3339))
		}
	}()
	<-done
}
