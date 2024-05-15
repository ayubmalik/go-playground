package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup
	c := sync.NewCond(&sync.Mutex{})

	subscribe := func(msg string) {
		go func() {
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fmt.Println("mama " + msg)
			wg.Done()
		}()
	}
	wg.Add(3)
	subscribe("a")
	subscribe("b")
	subscribe("c")
	time.Sleep(1 * time.Millisecond)
	c.Broadcast()
	wg.Wait()
}
