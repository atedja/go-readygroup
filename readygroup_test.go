package readygroup

import (
	"sync"
	"testing"
)

func compare(a, b []int) bool {
	if &a == &b {
		return true
	}

	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if b[i] != a[i] {
			return false
		}
	}

	return true
}

func TestReadyGroup(t *testing.T) {
	var count uint64 = 5
	order := make(chan int, count*2)

	var wg sync.WaitGroup
	wg.Add(int(count))

	rg := New()
	rg.Add(count)
	for i := uint64(0); i < count; i++ {
		go func(n int) {
			order <- 1
			rg.Ready()
			order <- 0
			wg.Done()
		}(int(i))
	}

	wg.Wait()
	close(order)

	// check final sequence
	seq := make([]int, 0, count*2)
	for o := range order {
		seq = append(seq, o)
	}

	if !compare(seq, []int{1, 1, 1, 1, 1, 0, 0, 0, 0, 0}) {
		t.Fatal("Out of sync! Ready() does not block")
	}
}

func TestReadyGroupMultipleAdds(t *testing.T) {
	var count uint64 = 5
	order := make(chan int, count*2)

	var wg sync.WaitGroup
	wg.Add(int(count))

	rg := New()
	rg.Add(1)
	rg.Add(2)
	rg.Add(2)
	for i := uint64(0); i < count; i++ {
		go func(n int) {
			order <- 1
			rg.Ready()
			order <- 0
			wg.Done()
		}(int(i))
	}

	wg.Wait()
	close(order)

	// check final sequence
	seq := make([]int, 0, count*2)
	for o := range order {
		seq = append(seq, o)
	}

	if !compare(seq, []int{1, 1, 1, 1, 1, 0, 0, 0, 0, 0}) {
		t.Fatal("Out of sync! Ready() does not block")
	}
}

func ExampleReadyGroup() {
	var rg ReadyGroup
	rg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			rg.Ready() // decrements counter and blocks here until counter reaches 0
			// do something...
		}()
	}
}
