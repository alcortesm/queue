// Package assert provides convenience methods to test queue
// implementations.
package assert

import (
	"fmt"

	"github.com/alcortesm/queue"
)

// Test represents a test that can be set to fail with a formated
// message using the Errorf method, like *testing.T from the
// standard library.
type Test interface {
	Errorf(format string, a ...interface{})
}

// Assert represent an association between a queue.Queue and a Test.
// Its methods can be used to assert that the queue behaves in
// particular ways.  If not, the test will be notified of the error with
// a sensible message, prefixed with the value of the Prefix field.
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
// arguments.  The error will be prefixed with the receiver's Prefix
// value.
func (a *Assert) errorf(format string, v ...interface{}) {
	a.t.Errorf("%s%s", a.Prefix, fmt.Sprintf(format, v...))
}

// Insert asserts that enqueuing the given numbers on the associated
// queue success.
func (a *Assert) Insert(numbers ...int) {
	for i, n := range numbers {
		if err := a.q.Insert(n); err != nil {
			a.errorf("wrong #%d Insert(%d): unexpected error: %q",
				i, n, err)
		}
	}
}

// InsertErrFull asserts that trying to insert any element on the
// associated queue fails with queue.ErrFull.
func (a *Assert) InsertErrFull() {
	if err := a.q.Insert(42); err != queue.ErrFull {
		a.errorf("wrong Insert: should have failed with queue.ErrFull, "+
			"got: %v", err)
	}
}

// Extract asserts that dequeuing len(expected) times from the
// associated queue success and return the expected values.
func (a *Assert) Extract(expected ...int) {
	for i, e := range expected {
		obtained, err := a.q.Extract()
		if err != nil {
			a.errorf("wrong #%d Extract: unexpected error: %q", i, err)
			return
		}
		n, ok := obtained.(int)
		if !ok {
			a.errorf("wrong #%d Extract: obtained (%v) cannot be cast to int",
				i, obtained)
			return
		}
		if n != e {
			a.errorf("wrong #%d Extract: expected %d, got %d", i, e, n)
		}
	}
}

// ExtractErrEmpty asserts that calling Extract on the associated queue
// fails with the error queue.ErrEmpty.
func (a *Assert) ExtractErrEmpty() {
	obtained, err := a.q.Extract()
	if err != queue.ErrEmpty {
		a.errorf("wrong Extract: should have failed with queue.ErrEmpty, "+
			"got: (obtained) %v and (error) %v", obtained, err)
	}
}

// Seq returns a slice with the natural numbers from min(a,b) (included)
// to max(a,b) (not included), sorted in ascending order.  This method
// can be used to generate test and expected data for Assert.Insert and
// Assert.Extract.
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
