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

func TestBounded(t *testing.T) {
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
		check.Bounded(t,
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

func TestInitiallyLenIsZero(t *testing.T) {
	for _, test := range []struct {
		context  string
		capacity int
	}{
		{"zero", 0},
		{"one", 1},
		{"two", 2},
		{"ten", 10},
	} {
		check.Len(t,
			mustNew(t, test.capacity),
			0,
			test.context)
	}
}

func TestInitiallyIsEmpty(t *testing.T) {
	for _, test := range []struct {
		context  string
		capacity int
		expected bool
	}{
		{"one", 0, true},
		{"one", 1, true},
		{"two", 2, true},
		{"ten", 10, true},
	} {
		check.IsEmpty(t,
			mustNew(t, test.capacity),
			test.expected,
			test.context)
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

func TestErrorWhenCapIsReached(t *testing.T) {
	for _, test := range []struct {
		context  string
		capacity int
	}{
		{"zero", 0},
		{"one", 1},
		{"two", 2},
		{"ten", 10},
	} {
		check.ErrorWhenCapIsReached(t,
			mustNew(t, test.capacity),
			test.context)
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
		check.HeadOKWhileFillingUpAndDepleting(t,
			mustNew(t, test.capacity),
			test.context)
	}

}
