package cbuf_test

import (
	"testing"

	"github.com/alcortesm/queue"
	"github.com/alcortesm/queue/cbuf"
	"github.com/alcortesm/queue/check"
)

func mustNew(t *testing.T, n int) queue.Queue {
	q, err := cbuf.New(n)
	if err != nil {
		t.Errorf("creating cbuf of size %d: %s", n, err)
	}
	var i interface{} = q
	if _, ok := i.(queue.Queue); !ok {
		t.Errorf("cbuf does not implement queue.Queue")
	}
	return q
}

func TestNewFailsWithNegativeCapacity(t *testing.T) {
	_, err := cbuf.New(-1)
	if err == nil {
		t.Errorf("new: negative capacity: expected error, got nil error")
	}
}

func TestIsBounded(t *testing.T) {
	for _, test := range []struct {
		context  string
		capacity int
		expected bool
	}{
		{"zero", 0, true},
		{"one", 1, true},
		{"two", 2, true},
		{"ten", 10, true},
	} {
		check.IsBounded(t,
			mustNew(t, test.capacity),
			test.expected,
			test.context)
	}
}

func TestCap(t *testing.T) {
	for _, test := range []struct {
		context  string
		capacity int
	}{
		{"zero", 0},
		{"one", 1},
		{"two", 2},
		{"ten", 10},
	} {
		expected := test.capacity
		check.CapBounded(t,
			mustNew(t, test.capacity),
			expected,
			test.context)
	}
}

func TestEmptyIsCorrectlyDetected(t *testing.T) {
	for _, test := range []struct {
		context  string
		capacity int
	}{
		{"one", 0},
		{"one", 1},
		{"two", 2},
		{"ten", 10},
	} {
		q := mustNew(t, test.capacity)
		// it must be empty initially...
		check.IsEmpty(t, q, true, test.context)
		check.Len(t, q, 0, test.context)
		check.HeadErrEmpty(t, q, test.context)
		check.DequeueErrEmpty(t, q, test.context)
		// and also if filled and depleted...
		check.FillEmptyWithNumbers(t, q, test.context)
		check.DepleteFullExpectingNumbers(t, q, test.context)
		check.IsEmpty(t, q, true, test.context)
		check.Len(t, q, 0, test.context)
		check.HeadErrEmpty(t, q, test.context)
		check.DequeueErrEmpty(t, q, test.context)
		// and even if filled and depleted a second time...
		check.FillEmptyWithNumbers(t, q, test.context)
		check.DepleteFullExpectingNumbers(t, q, test.context)
		check.IsEmpty(t, q, true, test.context)
		check.Len(t, q, 0, test.context)
		check.HeadErrEmpty(t, q, test.context)
		check.DequeueErrEmpty(t, q, test.context)
	}
}

func TestInitiallyAreNotFull(t *testing.T) {
	for _, test := range []struct {
		context  string
		capacity int
		expected bool
	}{
		{"zero", 0, true}, // zero capacity bounded queues are always full
		{"one", 1, false},
		{"two", 2, false},
		{"ten", 10, false},
	} {
		check.IsFull(t,
			mustNew(t, test.capacity),
			test.expected,
			test.context)
	}
}

func TestFull(t *testing.T) {
	for _, test := range []struct {
		context  string
		capacity int
	}{
		{"zero", 0},
		{"one", 1},
		{"two", 2},
		{"ten", 10},
	} {
		q := mustNew(t, test.capacity)
		check.FillEmptyWithNumbers(t, q, test.context)
		check.IsFull(t, q, true, test.context)
		check.EnqueueErrFull(t, q, test.context)
	}
}

func TestHeadOKWhileFillingUpAndDepleting(t *testing.T) {
	for _, test := range []struct {
		context  string
		capacity int
	}{
		{"zero", 0},
		{"one", 1},
		{"two", 2},
		{"ten", 10},
	} {
		q := mustNew(t, test.capacity)
		check.IsEmpty(t, q, true, test.context)
		check.HeadErrEmpty(t, q, test.context)
		// fill it up with numbers
		for i := 0; i < test.capacity; i++ {
			check.Enqueue(t, q, i, test.context)
			check.Head(t, q, 0, test.context)
			check.Head(t, q, 0, test.context) // head should not extract
		}
		check.IsFull(t, q, true, test.context)
		// extract all numbers
		for i := 0; i < test.capacity; i++ {
			check.Head(t, q, i, test.context)
			check.Head(t, q, i, test.context) // head should not extract
			check.Dequeue(t, q, i, test.context)
		}
		check.IsEmpty(t, q, true, test.context)
		check.HeadErrEmpty(t, q, test.context)
	}
}
