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
