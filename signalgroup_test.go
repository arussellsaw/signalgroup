package signalgroup

import (
	"sync"
	"testing"
)

func TestSend(t *testing.T) {
	var wg sync.WaitGroup
	g := New()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		c := g.Cursor()
		go func() {
			defer wg.Done()
			last := 0
			var s interface{}
			for {
				c, s = c.Wait()
				v, ok := s.(int)
				if !ok {
					t.Errorf("expected ok, got !ok")
				}
				t.Logf("got: %v", v)
				if v != last+1 {
					t.Errorf("expected %v, got %v", last+1, v)
				}
				last = v
				if v == 10 {
					return
				}
			}
		}()
	}
	for i := 1; i <= 10; i++ {
		t.Log(i)
		g.Send(i)
	}
	wg.Wait()
}

func TestBlockingSend(t *testing.T) {
	var wg sync.WaitGroup
	g := New()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		c := g.Cursor()
		go func() {
			defer wg.Done()
			last := 0
			for {
				newC, s := c.Wait()
				v, ok := s.(int)
				if !ok {
					t.Errorf("expected ok, got !ok")
				}
				t.Logf("got: %v", v)
				if v != last+1 {
					t.Errorf("expected %v, got %v", last+1, v)
				}
				last = v
				c.Done()
				c = newC
				if v == 10 {
					return
				}
			}
		}()
	}
	for i := 1; i <= 10; i++ {
		t.Log(i)
		g.BlockingSend(i, 10)
	}
	wg.Wait()
}
