// Package assert provides convenience ways to test queue.Queue
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

// Assert represent an association between a queue.Queue and a test.
// The methods from this type can be used form inside a test to
// manipulate a queue and assert that it behaves in particular ways.  If
// not, a test error will be automatically issued with a sensible
// failure message, prefixed with the value of the Prefix field.
type Assert struct {
	q queue.Queue
	t Test
	// The prefix to use in the error messages when an assertion fails.
	Prefix string
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

// Len asserts that the length of the associated queue is the expected value.
func (a *Assert) Len(expected int) {
	obtained := a.q.Len()
	if obtained != expected {
		a.errorf("wrong length: expected %d, got %d", expected, obtained)
	}
}

// IsEmpty asserts that calling IsEmpty on the associated queue returns the
// expected value.
func (a *Assert) IsEmpty(expected bool) {
	obtained := a.q.IsEmpty()
	if obtained != expected {
		a.errorf("wrong IsEmpty: expected %t, got %t", expected, obtained)
	}
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

// Head asserts that calling Head on the associated queue success and returns
// the expected number.
func (a *Assert) Head(expected int) {
	obtained, err := a.q.Head()
	if err != nil {
		a.errorf("wrong Head: unexpected error: %s", err)
		return
	}
	n, ok := obtained.(int)
	if !ok {
		a.errorf("wrong Head: obtained (%v) cannot be cast to int", obtained)
		return
	}
	if expected != n {
		a.errorf("wrong Head: expected %d, got %d", expected, n)
	}
}

// HeadErrEmpty asserts that calling Head on the associated queue fails with
// the error ErrEmpty.
func (a *Assert) HeadErrEmpty() {
	obtained, err := a.q.Head()
	if err != queue.ErrEmpty {
		a.errorf("wrong Head: should have failed with ErrEmpty, "+
			"but got %v and an %q error instead", obtained, err)
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
