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

func IsEmpty(t *testing.T, q queue.Queue, expected bool, context string) {
	obtained := q.IsEmpty()
	if obtained != expected {
		msg := fmt.Sprintf("wrong IsEmpty: expected %t, got %t",
			expected, obtained)
		error(t, context, msg)
	}
}

func IsFull(t *testing.T, q queue.Queue, expected bool, context string) {
	obtained := q.IsFull()
	if obtained != expected {
		msg := fmt.Sprintf("wrong IsFull: expected %t, got %t",
			expected, obtained)
		error(t, context, msg)
	}
}

func ErrorWhenCapIsReached(t *testing.T, q queue.Queue, context string) {
	Bounded(t, q, true, context)
	Len(t, q, 0, context)
	capacity, err := q.Cap()
	if err != nil {
		msg := fmt.Sprintf("unexpected error calling Cap: %q", err)
		error(t, context, msg)
	}
	// fill up the queue
	for i := range Seq(capacity) {
		if err := q.Enqueue(i); err != nil {
			msg := fmt.Sprintf(
				"unexpected error filling up queue:\n"+
					"on enqueue operation #%d: %s", i, err)
			error(t, context, msg)
		}
	}
	IsFull(t, q, true, context)
	// check that enqueueing once more gives ErrFull
	err = q.Enqueue(0)
	if err == nil {
		msg := fmt.Sprintf("enqueue on a full queue: return nil error")
		error(t, context, msg)
	}
	if err != queue.ErrFull {
		msg := fmt.Sprintf(
			"enqueue on a full queue: expected ErrFull, got %q", err)
		error(t, context, msg)
	}
}
