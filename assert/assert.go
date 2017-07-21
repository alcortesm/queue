// Package assert provides convenience ways to test queues
// implementations.
package assert

import (
	"fmt"

	"github.com/alcortesm/queue"
)

// Test represents a test that can be set to fail with a formated
// message using the Errorf method, for instance, *testing.T from the
// standard library.
type Test interface {
	Errorf(format string, a ...interface{})
}

// Assert represent an association between a queue.Queue and a Test.
// The methods from this type can be used to manipulate a queue and
// assert that it behaves in particular ways.  If not, the test will
// be notified of the error with a sensible message, prefixed
// with the value of the Prefix field.
type Assert struct {
	// The prefix to use in error messages when an assertion fails.
	Prefix string
	q      queue.Queue
	t      Test
}

// New returns a new association between the given queue and the given
// test.  The Prefix field is initialized to the empty string.
func New(t Test, q queue.Queue) *Assert {
	return &Assert{
		q:      q,
		t:      t,
		Prefix: "",
	}
}

// Calls the Errorf method on the test with the given formated
// arguments, prefixed with the receiver's Prefix value.
func (a *Assert) errorf(format string, v ...interface{}) {
	a.t.Errorf("%s%s", a.Prefix, fmt.Sprintf(format, v...))
}

// Enqueue asserts that enqueuing the given numbers on the associated
// queue success.
func (a *Assert) Enqueue(numbers ...int) {
	for i, n := range numbers {
		if err := a.q.Enqueue(n); err != nil {
			a.errorf("wrong #%d Enqueue(%d): unexpected error: %q",
				i, n, err)
		}
	}
}

// EnqueueErrFull asserts that trying to enqueue any element fails with
// queue.ErrFull.
func (a *Assert) EnqueueErrFull() {
	if err := a.q.Enqueue(42); err != queue.ErrFull {
		a.errorf("wrong Enqueue: should have failed with ErrFull, "+
			"but got %q instead", err)
	}
}

// Dequeue asserts that dequeuing len(expected) times from the
// associated queue success and returns the expected values.
func (a *Assert) Dequeue(expected ...int) {
	for i, e := range expected {
		obtained, err := a.q.Dequeue()
		if err != nil {
			a.errorf("wrong #%d Dequeue: unexpected error: %q", i, err)
			return
		}
		n, ok := obtained.(int)
		if !ok {
			a.errorf("wrong #%d Dequeue: obtained (%v) cannot be cast to int",
				i, obtained)
			return
		}
		if n != e {
			a.errorf("wrong #%d Dequeue: expected %d, got %d", i, e, n)
		}
	}
}

// DequeueErrEmpty asserts that calling Dequeue on the associated queue
// fails with the error ErrEmpty.
func (a *Assert) DequeueErrEmpty() {
	obtained, err := a.q.Dequeue()
	if err != queue.ErrEmpty {
		a.errorf("wrong Dequeue: should have failed with ErrEmpty, "+
			"but got %v and an %q error instead", obtained, err)
	}
}

// Seq returns a slice with the natural numbers from min(a,b) (included)
// to max(a,b) (not included), sorted in ascending order.  This method
// helps generating data for Assert.Enqueue and expected values for
// Assert.Dequeue.
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
