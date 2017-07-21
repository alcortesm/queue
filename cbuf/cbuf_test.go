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
		assert.Prefix = fmt.Sprintf("capacity %d:\n", cap)
		assert.ExtractErrEmpty()
	}
}

func TestInsertExtractOneElement(t *testing.T) {
	for _, cap := range []int{1, 2, 10} {
		q := mustNew(t, cap)
		assert := queueAssert.New(t, q)
		assert.Prefix = fmt.Sprintf("capacity %d:\n", cap)
		assert.Insert(12)
		assert.Extract(12)
	}
}

func TestZeroCapacityIsEmptyAndFull(t *testing.T) {
	q := mustNew(t, 0)
	assert := queueAssert.New(t, q)
	assert.ExtractErrEmpty()
	assert.InsertErrFull()
}

func TestFull(t *testing.T) {
	for _, cap := range []int{1, 2, 10} {
		q := mustNew(t, cap)
		assert := queueAssert.New(t, q)
		assert.Prefix = fmt.Sprintf("capacity %d:\n", cap)
		data := queueAssert.Seq(0, cap)
		assert.Insert(data...)
		assert.InsertErrFull()
	}
}

func TestFillHalfThenEmptyFillFullThenEmpty(t *testing.T) {
	for _, cap := range []int{2, 3, 10} {
		q := mustNew(t, cap)
		assert := queueAssert.New(t, q)
		assert.Prefix = fmt.Sprintf("capacity %d:\n", cap)
		// fill half of it
		full := queueAssert.Seq(0, cap)
		half := full[:len(full)/2]
		assert.Insert(half...)
		// empty it entirely
		assert.Extract(half...)
		assert.ExtractErrEmpty()
		// fill it completely
		assert.Insert(full...)
		assert.InsertErrFull()
		// empty it entirely
		assert.Extract(full...)
		assert.ExtractErrEmpty()
	}
}
