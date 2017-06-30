package check

import (
	"testing"

	"github.com/alcortesm/queue"
)

func Seq(n int) []int {
	ret := make([]int, n)
	for i, _ := range ret {
		ret[i] = i
	}
	return ret
}

func Bounded(t *testing.T, q queue.Queue, expected bool, context string) {
	obtained := q.Bounded()
	if obtained != expected {
		t.Errorf("%swrong bounded info: expected %t, got %t",
			context, expected, obtained)
	}
}

func CapInfinite(t *testing.T, q queue.Queue, context string) {
	capacity, err := q.Cap()
	if err == nil {
		t.Errorf(
			"%snil error calling Cap, ErrInfinite was expected, capacity was %d",
			context, capacity)
	}
	if err != queue.ErrInfinite {
		t.Errorf("%swrong error calling Cap: %s", context, err)
	}
}

func CapBounded(t *testing.T, q queue.Queue, expected int, context string) {
	obtained, err := q.Cap()
	if err != nil {
		t.Errorf("%sunexpected error calling Cap: %s", context, err)
	}
	if obtained != expected {
		t.Errorf("%swrong Cap: expected %d, got %d",
			context, expected, obtained)
	}
}

func Len(t *testing.T, q queue.Queue, expected int, context string) {
	obtained := q.Len()
	if obtained != expected {
		t.Errorf("%swrong Len: expected %d, got %d",
			context, expected, obtained)
	}
}

func Empty(t *testing.T, q queue.Queue, expected bool, context string) {
	obtained := q.Empty()
	if obtained != expected {
		t.Errorf("%swrong Empty: expected %t, got %t",
			context, expected, obtained)
	}
}

func Full(t *testing.T, q queue.Queue, expected bool, context string) {
	obtained := q.Full()
	if obtained != expected {
		t.Errorf("%swrong Full: expected %t, got %t",
			context, expected, obtained)
	}
}
