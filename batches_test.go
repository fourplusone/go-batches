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
	recv := (<-out).In

	select {
	case ch <- i3:
		t.Error("Expected channel to be full")
	default:

	}

	if recv != 3 {
		t.Error("Expected to receive 3, got", recv)
	}

	i2.In <- 5
	recv = (<-out).In
	if recv != 5 {
		t.Error("Expected to receive 5, got", recv)
	}
}

func TestCombiner(t *testing.T) {
	combiner := Combiner{
		CombineFunc: func(i []In) Out {
			return 43
		},
		Input: make(chan Item),
	}

	go func() { combiner.Process() }()

	a, x := combiner.Announce()
	b, y := combiner.Announce()
	c, z := combiner.Announce()

	a <- 1
	b <- 2
	c <- 3
	if <-x != 43 || <-y != 43 || <-z != 43 {
		t.Error("Expected to receive 43")
	}
}

func TestCancel(t *testing.T) {
	combiner := Combiner{
		CombineFunc: func(i []In) Out {
			return 43
		},
		Input: make(chan Item),
	}

	go func() { combiner.Process() }()

	a, x := combiner.Announce()
	b, _ := combiner.Announce()
	c, z := combiner.Announce()

	a <- 1
	close(b)
	c <- 3
	if <-x != 43 || <-z != 43 {
		t.Error("Expected to receive 43")
	}
}
