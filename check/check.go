package check

import (
	"fmt"
	"testing"

	"github.com/alcortesm/queue"
)

func error(t *testing.T, ctx string, msg string) {
	t.Errorf("context: %q\n %s", ctx, msg)
}

func IsBounded(t *testing.T, q queue.Queue, expected bool, context string) {
	obtained := q.IsBounded()
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
	IsBounded(t, q, true, context)
	Len(t, q, 0, context)
	capacity, err := q.Cap()
	if err != nil {
		msg := fmt.Sprintf("unexpected error calling Cap: %q", err)
		error(t, context, msg)
	}
	// fill up the queue
	for i := 0; i < capacity; i++ {
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

// HeadOKWhileFillingUpAndDepleting fills up the queue until full and empties
// it, while checking that the head method returns the correct value
// after every enqueue or dequeue operation.  Requires a bounded, empty queue.
func HeadOKWhileFillingUpAndDepleting(
	t *testing.T, q queue.Queue, context string) {
	IsBounded(t, q, true, context)
	IsEmpty(t, q, true, context)
	capacity, err := q.Cap()
	if err != nil {
		msg := fmt.Sprintf("unexpected error calling Cap: %q", err)
		error(t, context, msg)
	}
	HeadErrEmpty(t, q, context)
	// fill it up with numbers
	for i := 0; i < capacity; i++ {
		Enqueue(t, q, i, context)
		Head(t, q, 0, context)
		Head(t, q, 0, context) // head should not extract
	}
	IsFull(t, q, true, context)
	// extract all numbers
	for i := 0; i < capacity; i++ {
		Head(t, q, i, context)
		Head(t, q, i, context) // head should not extract
		Dequeue(t, q, i, context)
	}
	IsEmpty(t, q, true, context)
	HeadErrEmpty(t, q, context)
}

// HeadErrEmpty checks that calling the Head method on the given queue
// returns an ErrEmpty error.  If not, the test is failed and an error
// message, based on the context string, is reported to the testing library.
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

// FillWithNumbers enqueues consecutive numbers, starting from 0 into
// the given queue until it is full, checking that all enqueue
// operations are successful.  It expects a empty bounded queue.
func FillEmptyWithNumbers(t *testing.T, q queue.Queue, context string) {
	IsBounded(t, q, true, context)
	IsEmpty(t, q, true, context)
	cap, err := q.Cap()
	if err != nil {
		msg := fmt.Sprintf("unexpected error getting capacity: %q", err)
		error(t, context, msg)
	}
	for i := 0; i < cap; i++ {
		Enqueue(t, q, i, context)
	}
	IsFull(t, q, true, context)
}

// DepleteFullExpectingNumbers receives a bounded queue full of numbers,
// sorted from 0 at the head, to capacity-1 at the tail, this is,
// exactly as FillWithNumbers will do.  This function will dequeue all the
// numbers, checking that all operations are successful and that the numbers
// are extracted in the right order (0..capacity-1).
func DepleteFullExpectingNumbers(t *testing.T, q queue.Queue, context string) {
	IsBounded(t, q, true, context)
	IsFull(t, q, true, context)
	cap, err := q.Cap()
	if err != nil {
		msg := fmt.Sprintf("unexpected error getting capacity: %q", err)
		error(t, context, msg)
	}
	for i := 0; i < cap; i++ {
		Dequeue(t, q, i, context)
	}
	IsEmpty(t, q, true, context)
}
