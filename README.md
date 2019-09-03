# signalgroup

This library provides a basic tool for broadcasting data to an arbitrary number of goroutines, with minimal blocking on the send side, and deterministic ordering for recievers. Under the hood the library uses a linked list of channels that are closed to unblock when data is ready to be recieved. this allows us to repeatedly add signals to be recieved, even if consumers are slow to unblock and consume.

As a tradeoff there's a slightly inconvenient API, where you are given a new cursor handle after every wait call, here is an example of the group in action:

```go
package main

import (
	"fmt"
	"sync"

	"github.com/arussellsaw/signalgroup"
)

func main() {
	g := signalgroup.New()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(g.Cursor(), &wg)
	}
	for i := 0; i < 10; i++ {
		fmt.Println("sending", i)
		g.Send(i)
	}
	wg.Wait()
}

func worker(c *signalgroup.Cursor, wg *sync.WaitGroup) {
	defer wg.Done()
	var s interface{}
	for {
		c, s = c.Wait()
		fmt.Println("got", s)
		v := s.(int)
		if v == 9 {
			return
		}
	}
}
```

## codegen
this lib uses `interface{}` for ease of use, but would be a great candidate for code generation to have typed signalgroups, i might work on this if i get bored but i welcome contributors.
