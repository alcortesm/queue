package cbuf_test

import (
	"fmt"
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
		t.Errorf("new: negative capacity: expected error, got nil")
	}
}

func TestInitiallyEmpty(t *testing.T) {
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
		check.Len(t, q, 0, test.context)
		check.HeadErrEmpty(t, q, test.context)
		check.DequeueErrEmpty(t, q, test.context)
	}
}

func TestOneElementNotFull(t *testing.T) {
	for _, test := range []struct {
		context  string
		capacity int
	}{
		{"two", 2},
		{"ten", 10},
	} {
		q := mustNew(t, test.capacity)
		check.Enqueue(t, q, 12, test.context)
		check.IsEmpty(t, q, false, test.context)
		check.Len(t, q, 1, test.context)
		check.Dequeue(t, q, 12, test.context)
	}
}

func TestFull(t *testing.T) {
	for _, test := range []struct {
		context  string
		capacity int
	}{
		{"one", 1},
		{"two", 2},
		{"ten", 10},
	} {
		q := mustNew(t, test.capacity)
		check.EnqueueAll(t, q, check.Seq(0, test.capacity), test.context)
		check.Len(t, q, test.capacity, test.context)
		check.IsEmpty(t, q, false, test.context)
		check.EnqueueErrFull(t, q, test.context)
	}
}

func TestFillUpThenEmpty(t *testing.T) {
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
		// fill half of it and empty it
		check.EnqueueAll(t, q, check.Seq(0, test.capacity/2), test.context)
		check.DequeueAll(t, q, check.Seq(0, test.capacity/2), test.context)
		check.IsEmpty(t, q, true, test.context)
		check.Len(t, q, 0, test.context)
		check.HeadErrEmpty(t, q, test.context)
		check.DequeueErrEmpty(t, q, test.context)
		// now fill it up and empty it
		check.EnqueueAll(t, q, check.Seq(0, test.capacity), test.context)
		check.DequeueAll(t, q, check.Seq(0, test.capacity), test.context)
		check.IsEmpty(t, q, true, test.context)
		check.Len(t, q, 0, test.context)
		check.HeadErrEmpty(t, q, test.context)
		check.DequeueErrEmpty(t, q, test.context)
	}
}

func TestAllWhileFillingUpAndDepleting(t *testing.T) {
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
		// fill it up with numbers
		for i := 0; i < test.capacity; i++ {
			context := fmt.Sprintf("%s: enqueuing %d", test.context, i)
			check.Len(t, q, i, context)
			check.Enqueue(t, q, i, context)
			check.IsEmpty(t, q, false, context)
			check.Head(t, q, 0, context)
			check.Head(t, q, 0, context) // head should be idempotent
		}
		context := test.context + ": filled up"
		check.Len(t, q, test.capacity, context)
		check.EnqueueErrFull(t, q, context)
		// extract all numbers
		for i := 0; i < test.capacity; i++ {
			context := fmt.Sprintf("%s: dequeuing %d", test.context, i)
			check.IsEmpty(t, q, false, context)
			check.Len(t, q, test.capacity-i, context)
			check.Head(t, q, i, context)
			check.Head(t, q, i, context) // head should be idempotent
			check.Dequeue(t, q, i, context)
		}
		context = test.context + ": depleted"
		check.IsEmpty(t, q, true, test.context)
		check.Len(t, q, 0, test.context)
		check.HeadErrEmpty(t, q, test.context)
		check.DequeueErrEmpty(t, q, test.context)
	}
}
