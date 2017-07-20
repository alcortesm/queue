package cbuf_test

import (
	"fmt"
	"testing"

	"github.com/alcortesm/queue"
	queueAssert "github.com/alcortesm/queue/assert"
	"github.com/alcortesm/queue/cbuf"
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
	for _, cap := range []int{0, 1, 2, 10} {
		q := mustNew(t, cap)
		assert := queueAssert.New(t, q)
		assert.Prefix = fmt.Sprintf("capacity %d", cap)
		assert.Len(0)
		assert.IsEmpty(true)
		assert.HeadErrEmpty()
		assert.DequeueErrEmpty()
	}
}

func TestOneElementNotFull(t *testing.T) {
	for _, cap := range []int{2, 10} {
		q := mustNew(t, cap)
		assert := queueAssert.New(t, q)
		assert.Prefix = fmt.Sprintf("capacity %d", cap)
		assert.Enqueue(12)
		assert.IsEmpty(false)
		assert.Len(1)
		assert.Dequeue(12)
	}
}

func TestZeroCapacityIsEmptyAndFull(t *testing.T) {
	q := mustNew(t, 0)
	assert := queueAssert.New(t, q)
	assert.IsEmpty(true)
	assert.EnqueueErrFull()
}

func TestFull(t *testing.T) {
	for _, cap := range []int{1, 2, 10} {
		q := mustNew(t, cap)
		assert := queueAssert.New(t, q)
		assert.Prefix = fmt.Sprintf("capacity %d", cap)
		data := queueAssert.Seq(0, cap)
		assert.Enqueue(data...)
		assert.Len(len(data))
		assert.IsEmpty(false)
		assert.EnqueueErrFull()
	}
}

func TestFillHalfThenEmptyFillFullThenEmpty(t *testing.T) {
	for _, cap := range []int{2, 3, 10} {
		q := mustNew(t, cap)
		assert := queueAssert.New(t, q)
		assert.Prefix = fmt.Sprintf("capacity %d", cap)
		// fill half of it
		full := queueAssert.Seq(0, cap)
		half := full[:len(full)/2]
		assert.Enqueue(half...)
		assert.Len(len(half))
		assert.IsEmpty(false)
		// empty it entirely
		assert.Dequeue(half...)
		assert.IsEmpty(true)
		assert.Len(0)
		assert.HeadErrEmpty()
		assert.DequeueErrEmpty()
		// fill it completely
		assert.Enqueue(full...)
		assert.Len(len(full))
		assert.IsEmpty(false)
		assert.EnqueueErrFull()
		// empty it entirely
		assert.Dequeue(full...)
		assert.IsEmpty(true)
		assert.Len(0)
		assert.HeadErrEmpty()
		assert.DequeueErrEmpty()
	}
}

func TestAllWhileEnqueueFullAndDequeueUntilEmpty(t *testing.T) {
	for _, cap := range []int{0, 1, 2, 3, 10} {
		q := mustNew(t, cap)
		assert := queueAssert.New(t, q)
		// enqueue number until full
		assert.Prefix = fmt.Sprintf("capacity %d", cap)
		for i := 0; i < cap; i++ {
			assert.Prefix = fmt.Sprintf("%s: enqueuing #%d", assert.Prefix, i)
			assert.Len(i)
			assert.Enqueue(i)
			assert.IsEmpty(false)
			assert.Head(0)
			assert.Head(0) // head should be idempotent
		}
		assert.Prefix = fmt.Sprintf("capacity %d: full", cap)
		assert.Len(cap)
		assert.EnqueueErrFull()
		// extract all numbers
		assert.Prefix = fmt.Sprintf("capacity %d", cap)
		for i := 0; i < cap; i++ {
			assert.Prefix = fmt.Sprintf("%s: dequeuing #%d", assert.Prefix, i)
			assert.IsEmpty(false)
			assert.Len(cap - i)
			assert.Head(i)
			assert.Head(i) // head should be idempotent
			assert.Dequeue(i)
		}
		assert.Prefix = fmt.Sprintf("capacity %d: empty", cap)
		assert.IsEmpty(true)
		assert.Len(0)
		assert.HeadErrEmpty()
		assert.DequeueErrEmpty()
	}
}
