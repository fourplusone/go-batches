package batches

import "testing"

func TestCollectAndResolveEmptyInput(t *testing.T) {
	ch := make(chan Item)
	out := collectAndResolve(ch)

	close(ch)
	_, more := <-out

	if more {
		t.Error("Channel was not closed")
	}
}

func TestCollectAndResolveAll(t *testing.T) {
	ch := make(chan Item)
	out := collectAndResolve(ch)

	i1 := Item{
		In:  make(chan In),
		Out: make(chan Out),
	}

	i2 := Item{
		In:  make(chan In),
		Out: make(chan Out),
	}

	i3 := Item{
		In:  make(chan In),
		Out: make(chan Out),
	}

	ch <- i1
	ch <- i2
	ch <- i3
	//close(ch)

	i1.In <- 3
	recv := (<-out).In

	if recv != 3 {
		t.Error("Expected to receive 3, got", recv)
	}

	i2.In <- 5
	recv = (<-out).In
	if recv != 5 {
		t.Error("Expected to receive 5, got", recv)
	}

	i3.In <- 10
	recv = (<-out).In
	if recv != 10 {
		t.Error("Expected to receive 10, got", recv)
	}

}

func TestCollectAndResolvePartial(t *testing.T) {
	ch := make(chan Item)
	out := collectAndResolve(ch)

	i1 := Item{
		In:  make(chan In),
		Out: make(chan Out),
	}

	i2 := Item{
		In:  make(chan In),
		Out: make(chan Out),
	}

	i3 := Item{
		In:  make(chan In),
		Out: make(chan Out),
	}

	ch <- i1
	ch <- i2
	i1.In <- 3
	select {
	case ch <- i3:
		t.Error("Expected channel to be full")
	default:

	}

	recv := (<-out).In
	if recv != 3 {
		t.Error("Expected to receive 3, got", recv)
	}

	i2.In <- 5
	recv = (<-out).In
	if recv != 5 {
		t.Error("Expected to receive 5, got", recv)
	}

}
