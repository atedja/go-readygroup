package readygroup

import (
	"sync/atomic"
)

// ReadyGroup is used in opposite situations to sync.WaitGroup.
// It blocks all goroutines until all are ready, then unblocks them all
// in no particular order.
type ReadyGroup struct {
	done        chan struct{}
	ready       chan struct{}
	readyTotal  uint64
	total       uint64
	waitRunning uint64
}

func New() *ReadyGroup {
	rg := &ReadyGroup{}
	rg.done = make(chan struct{})
	rg.ready = make(chan struct{})
	rg.readyTotal = 0
	rg.total = 0
	rg.waitRunning = 0
	return rg
}

// Add adds the total number of goroutines to wait, which cannot be negative, to the ReadyGroup counter.
func (self *ReadyGroup) Add(total uint64) {
	atomic.AddUint64(&self.total, total)
	go self.wait()
}

// Ready decrements the ReadyGroup counter, and blocks until counter reaches 0.
func (self *ReadyGroup) Ready() {
	self.ready <- struct{}{}
	<-self.done
}

func (self *ReadyGroup) wait() {
	// Check if there's already a wait() running
	if atomic.CompareAndSwapUint64(&self.waitRunning, 1, 1) {
		return
	}
	atomic.AddUint64(&self.waitRunning, 1)

	for range self.ready {
		self.readyTotal++
		atomic.AddUint64(&self.total, ^uint64(0))
		if atomic.CompareAndSwapUint64(&self.total, 0, 0) {
			for i := uint64(0); i < self.readyTotal; i++ {
				self.done <- struct{}{}
			}
			self.readyTotal = 0
			break
		}
	}

	atomic.AddUint64(&self.waitRunning, ^uint64(0))
}
