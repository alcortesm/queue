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

func (q *zeroCap) Insert(_ interface{}) error    { return queue.ErrFull }
func (q *zeroCap) Extract() (interface{}, error) { return nil, queue.ErrEmpty }

// a mock queue that always returns 42 when dequeue or head are called.
type canDequeue struct {
	zeroCap
}

func (q *canDequeue) Extract() (interface{}, error) { return 42, nil }

// a mock queue that always returns the string "42" when dequeue or head
// are called.
type canDequeueAString struct {
	zeroCap
}

func (q *canDequeueAString) Extract() (interface{}, error) { return "42", nil }

// a mock queue where Enqueue never fails
type neverFull struct {
	zeroCap
}

func (q *neverFull) Insert(_ interface{}) error { return nil }

func TestNewHasEmptyPrefix(t *testing.T) {
	a := assert.New(nil, nil)
	if a.Prefix != "" {
		t.Errorf("bad prefix: expected \"\", got %q", a.Prefix)
	}
}

func TestEnqueueTruePositive(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(neverFull))
	a.Insert(42)
	if mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestEnqueueTrueNegative(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(zeroCap))
	a.Insert(42)
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestEnqueueErrFullTruePositive(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(zeroCap))
	a.InsertErrFull()
	if mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestEnqueueErrFullTrueNegative(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(neverFull))
	a.InsertErrFull()
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestDequeueTruePositive(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(canDequeue))
	a.Extract(42)
	if mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestDequeueOnEmpty(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(zeroCap))
	a.Extract(42)
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestDequeueCannotCast(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(canDequeueAString))
	a.Extract(42)
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestDequeueWrongExpected(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(canDequeue))
	a.Extract(43)
	if !mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestDequeueErrEmptyTruePositive(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(zeroCap))
	a.ExtractErrEmpty()
	if mt.ErrorHasBeenCalled() {
		t.Error()
	}
}

func TestDequeueErrEmptyTrueNegative(t *testing.T) {
	mt := new(mockTest)
	a := assert.New(mt, new(canDequeue))
	a.ExtractErrEmpty()
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
