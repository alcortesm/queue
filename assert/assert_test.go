package assert_test

import (
	"reflect"
	"testing"

	"github.com/alcortesm/queue"
	"github.com/alcortesm/queue/assert"
)

// a mock test that can be asked if the Errorf was called or not.
type mockTest bool // has the Errorf method been called?

func (t *mockTest) Errorf(format string, a ...interface{}) { *t = mockTest(true) }
func (t *mockTest) ErrorHasBeenCalled() bool               { return bool(*t) }

// a mock queue with 0 capacity
type zeroCap struct{}

func (q *zeroCap) Len() int                      { return 0 }
func (q *zeroCap) IsEmpty() bool                 { return true }
func (q *zeroCap) Enqueue(_ interface{}) error   { return queue.ErrFull }
func (q *zeroCap) Dequeue() (interface{}, error) { return nil, queue.ErrEmpty }
func (q *zeroCap) Head() (interface{}, error)    { return nil, queue.ErrEmpty }

// a mock queue that always returns 42 when dequeue or head are called.
type canDequeue struct {
	zeroCap
}

func (q *canDequeue) Dequeue() (interface{}, error) { return 42, nil }
func (q *canDequeue) Head() (interface{}, error)    { return 42, nil }

// a mock queue that always returns the string "42" when dequeue or head
// are called.
type canDequeueAString struct {
	zeroCap
}

func (q *canDequeueAString) Dequeue() (interface{}, error) { return "42", nil }
func (q *canDequeueAString) Head() (interface{}, error)    { return "42", nil }

// a mock queue where Enqueue never fails
type neverFull struct {
	zeroCap
}

func (q *neverFull) Enqueue(_ interface{}) error { return nil }

func TestNewHasEmptyPrefix(t *testing.T) {
	a := assert.New(nil, nil)
	if a.Prefix != "" {
		t.Errorf("bad prefix: expected \"\", got %q", a.Prefix)
	}
}

func TestLenTruePositive(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(zeroCap))
	a.Len(0)
	if mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestLenTrueNegative(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(zeroCap))
	a.Len(1)
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestIsEmptyTruePositive(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(zeroCap))
	a.IsEmpty(true)
	if mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestIsEmptyTrueNegative(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(zeroCap))
	a.IsEmpty(false)
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestEnqueueTruePositive(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(neverFull))
	a.Enqueue(42)
	if mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestEnqueueTrueNegative(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(zeroCap))
	a.Enqueue(42)
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestEnqueueErrFullTruePositive(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(zeroCap))
	a.EnqueueErrFull()
	if mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestEnqueueErrFullTrueNegative(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(neverFull))
	a.EnqueueErrFull()
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestHeadTruePositive(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(canDequeue))
	a.Head(42)
	if mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestHeadOnEmpty(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(zeroCap))
	a.Head(42)
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestHeadCannotCast(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(canDequeueAString))
	a.Head(42)
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestHeadWrongExpected(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(canDequeue))
	a.Head(43)
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestHeadErrEmptyTruePositive(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(zeroCap))
	a.HeadErrEmpty()
	if mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestHeadErrEmptyTrueNegative(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(canDequeue))
	a.HeadErrEmpty()
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestDequeueTruePositive(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(canDequeue))
	a.Dequeue(42)
	if mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestDequeueOnEmpty(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(zeroCap))
	a.Dequeue(42)
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestDequeueCannotCast(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(canDequeueAString))
	a.Dequeue(42)
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestDequeueWrongExpected(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(canDequeue))
	a.Dequeue(43)
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestDequeueErrEmptyTruePositive(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(zeroCap))
	a.DequeueErrEmpty()
	if mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestDequeueErrEmptyTrueNegative(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(canDequeue))
	a.DequeueErrEmpty()
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestSeq(t *testing.T) {
	for _, test := range []struct {
		from, to int
		expected []int
	}{
		{0, 0, []int{}},
		{0, 1, []int{0}},
		{1, 0, []int{0}},
		{0, 2, []int{0, 1}},
		{2, 0, []int{0, 1}},
		{-2, 2, []int{-2, -1, 0, 1}},
		{2, -2, []int{-2, -1, 0, 1}},
	} {
		obtained := assert.Seq(test.from, test.to)
		if !reflect.DeepEqual(obtained, test.expected) {
			t.Errorf("from %d to %d: expected %v, obtained %v",
				test.from, test.to, test.expected, obtained)
		}
	}
}
