package async

import "sync"

// CancelGroup cancels all operations in a group.
type CancelGroup struct {
	mu     sync.Mutex
	cc     []Canceller
	cancel bool
}

// NewCancelGroup creates a new cancel group.
func NewCancelGroup() *CancelGroup {
	return &CancelGroup{}
}

// Add adds a canceller to the group, immediately cancelling it if the group is cancelled.
func (g *CancelGroup) Add(c Canceller) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.cancel {
		c.Cancel()
		return
	}

	g.cc = append(g.cc, c)
}

// Cancel cancels all operations in the group.
func (g *CancelGroup) Cancel() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.cancel {
		return
	}

	g.cancel = true
	for _, c := range g.cc {
		c.Cancel()
	}
}

// CancelWait cancels all operations in the group and awaits them.
func (g *CancelGroup) CancelWait() {
	g.Cancel()
	g.Wait()
}

// Wait awaits all operations in the group.
func (g *CancelGroup) Wait() {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, c := range g.cc {
		w, ok := c.(Waiter)
		if !ok {
			continue
		}

		<-w.Wait()
	}
}
