package readygroup

import (
	"testing"
)

func TestReadyGroup(t *testing.T) {
	var count uint = 100
	res := make(chan int, count+1)
	var rg ReadyGroup
	rg.Add(count)
	for i := uint(0); i < count; i++ {
		go func(n int) {
			rg.Ready()
			res <- n
		}(int(i))
	}
	res <- -1
	rg.Go()

	first := <-res
	if first != -1 {
		t.Fatal("Ready does not block!")
	}
}

func ExampleReadyGroup() {
	var rg ReadyGroup
	rg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			rg.Ready() // blocks here until Go() gets unblocked.
			// do something...
		}()
	}
	rg.Go() // blocks and only unblocks when all goroutines are ready.
}
