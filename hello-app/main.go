package main

import (
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

	ticker := time.Tick(3 * time.Second)
	go func() {
		for next := range ticker {
			slog.Info("Hello ", "time", next.Format(time.RFC3339))
		}
	}()
	<-done
}
