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
		assert.DequeueErrEmpty()
	}
}

func TestInsertExtractOneElement(t *testing.T) {
	for _, cap := range []int{1, 2, 10} {
		q := mustNew(t, cap)
		assert := queueAssert.New(t, q)
		assert.Prefix = fmt.Sprintf("capacity %d", cap)
		assert.Enqueue(12)
		assert.Dequeue(12)
	}
}

func TestZeroCapacityIsEmptyAndFull(t *testing.T) {
	q := mustNew(t, 0)
	assert := queueAssert.New(t, q)
	assert.DequeueErrEmpty()
	assert.EnqueueErrFull()
}

func TestFull(t *testing.T) {
	for _, cap := range []int{1, 2, 10} {
		q := mustNew(t, cap)
		assert := queueAssert.New(t, q)
		assert.Prefix = fmt.Sprintf("capacity %d", cap)
		data := queueAssert.Seq(0, cap)
		assert.Enqueue(data...)
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
		// empty it entirely
		assert.Dequeue(half...)
		assert.DequeueErrEmpty()
		// fill it completely
		assert.Enqueue(full...)
		assert.EnqueueErrFull()
		// empty it entirely
		assert.Dequeue(full...)
		assert.DequeueErrEmpty()
	}
}
