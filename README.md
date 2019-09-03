# signalgroup

This library provides a basic tool for broadcasting data to an arbitrary number of goroutines, with minimal blocking on the send side, and deterministic ordering for recievers. Under the hood the library uses a linked list of channels that are closed to unblock when data is ready to be recieved. this allows us to repeatedly add signals to be recieved, even if consumers are slow to unblock and consume.

As a tradeoff there's a slightly inconvenient API, where you are given a new cursor handle after every wait call, here is an example of the group in action:

```go
func TestGroup(t *testing.T) {
	var wg sync.WaitGroup
	g := New()

	// start 10 workers consuming signals
	for i := 0; i < 10; i++ {
		wg.Add(1)
		// get the current cursor for the group
		c := g.Cursor()
		go func() {
			defer wg.Done()
			last := 0
			for {
				var s interface{}

				// wait for a signal, updating the cursor when we recieve it
				c, s = c.Wait()
				v, ok := s.(int)
				if !ok {
					t.Errorf("expected ok, got !ok")
				}

				t.Logf("got: %v", v)

				// this assertion validates ordering
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

	// send 10 integer signals to the group
	for i := 1; i <= 10; i++ {
		t.Logf("sending: %v", i)
		g.Send(&Signal{Value: i})
	}
	wg.Wait()
}
```

## codegen
this lib uses `interface{}` for ease of use, but would be a great candidate for code generation to have typed signalgroups, i might work on this if i get bored but i welcome contributors.
