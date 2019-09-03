package signalgroup

import "sync"

// Signal is a value broadcast to a group
type Signal struct {
	Type   string
	Value  interface{}
	Params map[string]string
}

// New creates a Group with the given bufSize
func New() *Group {
	g := &Group{
		cs: &Cursor{
			c: make(chan struct{}),
		},
	}
	return g
}

// Group orchestrates broadcasting signals to a group of consumers
type Group struct {
	in chan *Signal
	mu sync.Mutex
	cs *Cursor
}

// Send a signal to the group
func (g *Group) Send(s *Signal) {
	g.mu.Lock()
	if g.cs.next != nil {
		panic("signalgroup: trying to send to already populated Cursor")
	}

	newCursor := &Cursor{
		c: make(chan struct{}),
	}

	// give waiters holding the existing cursor a link to the new one
	g.cs.next = newCursor
	// take a reference to the old cursor
	oldCursor := g.cs
	// point to new cursor
	g.cs = newCursor
	// set the signal on that cursor
	oldCursor.s = s
	// unblock the waiters
	close(oldCursor.c)
	g.mu.Unlock()
}

// Cursor returns the current cursor for the signalgroup
func (g *Group) Cursor() *Cursor {
	g.mu.Lock()
	c := g.cs
	g.mu.Unlock()
	return c
}

// Cursor is used to recieve a signal
type Cursor struct {
	c    chan struct{}
	s    *Signal
	next *Cursor
}

// Wait blocks until a signal is recieved
func (c *Cursor) Wait() (*Cursor, *Signal) {
	<-c.c
	s := c.s
	return c.next, s
}