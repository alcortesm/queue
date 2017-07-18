package check

import (
	"fmt"
	"testing"

	"github.com/alcortesm/queue"
)

func error(t *testing.T, ctx string, msg string) {
	t.Errorf("context: %q\n %s", ctx, msg)
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

// HeadErrEmpty checks that calling the Head method on the given queue
// returns an ErrEmpty error.  If not, the test is failed and an error
// message, based on the context string, is reported to the testing
// library.
func HeadErrEmpty(t *testing.T, q queue.Queue, context string) {
	_, err := q.Head()
	if err == nil {
		error(t, context, "Head return a nil error")
	}
	if err != queue.ErrEmpty {
		msg := fmt.Sprintf(
			"Head return an error different than ErrEmpty: %q", err)
		error(t, context, msg)
	}
}

// Head checks that calling Head ont the given queue is successful and
// returns the given integer.  Otherwise, an error is notified to the
// test library using an error message prefixed with the context string.
func Head(t *testing.T, q queue.Queue, e int, context string) {
	o, err := q.Head()
	if err != nil {
		msg := fmt.Sprintf("unexpected error: %q", err)
		error(t, context, msg)
	}
	oint, ok := o.(int)
	if !ok {
		msg := fmt.Sprintf("head: returned value cannot be cast to int")
		error(t, context, msg)
	}
	if oint != e {
		msg := fmt.Sprintf("head: expected %d, got %d", e, oint)
		error(t, context, msg)
	}
}

// Enqueue checks that enqueing the given integer into the given queue returns
// no error.  Otherwise, an error is notified to the test library using an
// error message prefixed with the context string.
func Enqueue(t *testing.T, q queue.Queue, e int, context string) {
	if err := q.Enqueue(e); err != nil {
		msg := fmt.Sprintf("enqueue: unexpected error: %q", err)
		error(t, context, msg)
	}
}

// EnqueueErrFull checks that trying to enqueue a number in the queue
// returns the error queue.ErrFull.  Otherwise, an error is notified to
// the test library using an error message prefixed with the context
// string.
func EnqueueErrFull(t *testing.T, q queue.Queue, context string) {
	if err := q.Enqueue(0); err != queue.ErrFull {
		msg := fmt.Sprintf("enqueue: expected ErrFull, got %q", err)
		error(t, context, msg)
	}
}

// Dequeue checks that dequeing from the given queue is successful and
// returns the given integer.  Otherwise, an error is notified to the
// test library using an error message prefixed with the context string.
func Dequeue(t *testing.T, q queue.Queue, e int, context string) {
	o, err := q.Dequeue()
	if err != nil {
		msg := fmt.Sprintf("dequeue: unexpected error: %q", err)
		error(t, context, msg)
	}
	oint, ok := o.(int)
	if !ok {
		msg := fmt.Sprintf("dequeue: returned value cannot be cast to int")
		error(t, context, msg)
	}
	if oint != e {
		msg := fmt.Sprintf("dequeue: expected %d, got %d", e, oint)
		error(t, context, msg)
	}
}

// DequeueErrEmpty checks that calling the Dequeue method on the given queue
// returns an ErrEmpty error.  If not, the test is failed and an error
// message, based on the context string, is reported to the testing library.
func DequeueErrEmpty(t *testing.T, q queue.Queue, context string) {
	_, err := q.Dequeue()
	if err == nil {
		error(t, context, "dequeue returned a nil error")
	}
	if err != queue.ErrEmpty {
		msg := fmt.Sprintf(
			"dequeue returned an error different than ErrEmpty: %q", err)
		error(t, context, msg)
	}
}

func EnqueueAll(t *testing.T, q queue.Queue, data []int, context string) {
	for i, e := range data {
		innerCtx := fmt.Sprintf("%s: element %d", context, i)
		Enqueue(t, q, e, innerCtx)
	}
}

func DequeueAll(t *testing.T, q queue.Queue, expected []int, context string) {
	for i, e := range expected {
		innerCtx := fmt.Sprintf("%s: element %d", context, i)
		Dequeue(t, q, e, innerCtx)
	}
}

// Seq returns a slice with the natural numbers from min(a,b) (included)
// to max(a,b) (not included), sorted in ascending order.
func Seq(a, b int) []int {
	if b < a {
		b, a = a, b
	}
	count := b - a
	ret := make([]int, count)
	for i := 0; i < count; i++ {
		ret[i] = a + i
	}
	return ret
}
