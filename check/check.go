package check

import (
	"fmt"
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

func error(t *testing.T, ctx string, msg string) {
	t.Errorf("context: %q\n %s", ctx, msg)
}

func Bounded(t *testing.T, q queue.Queue, expected bool, context string) {
	obtained := q.Bounded()
	if obtained != expected {
		msg := fmt.Sprintf("wrong bounded info: expected %t, got %t",
			expected, obtained)
		error(t, context, msg)
	}
}

func CapInfinite(t *testing.T, q queue.Queue, context string) {
	capacity, err := q.Cap()
	if err == nil {
		msg := fmt.Sprintf("nil error calling Cap, "+
			"ErrInfinite was expected, capacity was %d",
			capacity)
		error(t, context, msg)
	}
	if err != queue.ErrInfinite {
		t.Errorf("%swrong error calling Cap: %s", context, err)
	}
}

func CapBounded(t *testing.T, q queue.Queue, expected int, context string) {
	obtained, err := q.Cap()
	if err != nil {
		msg := fmt.Sprintf("unexpected error calling Cap: %q", err)
		error(t, context, msg)
	}
	if obtained != expected {
		msg := fmt.Sprintf("wrong capacity: expected %d, got %d",
			expected, obtained)
		error(t, context, msg)
	}
}

func Len(t *testing.T, q queue.Queue, expected int, context string) {
	obtained := q.Len()
	if obtained != expected {
		msg := fmt.Sprintf("wrong Len: expected %d, got %d",
			expected, obtained)
		error(t, context, msg)
	}
}

func Empty(t *testing.T, q queue.Queue, expected bool, context string) {
	obtained := q.Empty()
	if obtained != expected {
		msg := fmt.Sprintf("wrong Empty: expected %t, got %t",
			expected, obtained)
		error(t, context, msg)
	}
}

func Full(t *testing.T, q queue.Queue, expected bool, context string) {
	obtained := q.Full()
	if obtained != expected {
		msg := fmt.Sprintf("wrong Full: expected %t, got %t",
			expected, obtained)
		error(t, context, msg)
	}
}
