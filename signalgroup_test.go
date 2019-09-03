package signalgroup

import (
	"sync"
	"testing"
)

func TestGroup(t *testing.T) {
	var wg sync.WaitGroup
	g := New()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			last := 0
			c := g.Cursor()
			var s *Signal
			for {
				c, s = c.Wait()
				v, ok := s.Value.(int)
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
		g.Send(&Signal{Value: i})
	}
	wg.Wait()
}
