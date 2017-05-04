package readygroup

// ReadyGroup is used in opposite situation to sync.WaitGroup.
// sync.WaitGroup is typically used to wait for goroutines to finish their executions,
// while ReadyGroup is used to make all goroutines wait before they can begin execution.
// It guarantees to block all goroutines until they are ready.
type ReadyGroup struct {
	groups chan struct{}
	done   chan struct{}
	total  uint
}

// Adds delta, which cannot be negative, to the ReadyGroup counter.
func (self *ReadyGroup) Add(delta uint) {
	self.total = delta
	self.groups = make(chan struct{}, delta)
	self.done = make(chan struct{}, delta)
}

// Ready decrements the ReadyGroup counter, and blocks until Go gets unblocked.
func (self *ReadyGroup) Ready() {
	self.groups <- struct{}{}
	<-self.done
}

// Go blocks until the counter reaches 0.
func (self *ReadyGroup) Go() {
	var counter uint = 0
	for counter < self.total {
		<-self.groups
		counter++
	}

	for counter > 0 {
		self.done <- struct{}{}
		counter--
	}
}
