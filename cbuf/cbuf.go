/*
Package cbuf implements an in-memory bounded queue with constant
worst-case time complexity on all its methods.

Internally, it uses a fixed size circular buffer allocated upon construction.
This implementation is not thread-safe.
*/
package cbuf

import (
	"fmt"

	"github.com/alcortesm/queue"
)

// CBuf values are circular buffers of constant size.  They are not
// zero-value safe, use the New function below to instantiate them.
type CBuf struct {
	cap   int
	len   int
	elems []interface{}
	head  int
}

// New returns a new circular buffer of the given capacity and a nil
// error on success.  If the capacity is negative it returns nil and an
// error.  Zero capacity buffers are allowed; they will be empty and
// full at the same time.
func New(capacity int) (*CBuf, error) {
	if capacity < 0 {
		return nil, fmt.Errorf("negative capacity (%d)",
			capacity)
	}
	return &CBuf{
		cap:   capacity,
		elems: make([]interface{}, capacity),
	}, nil
}

func (c *CBuf) next(n int) int {
	return (n + 1) % c.cap
}

func (c *CBuf) tail() int {
	return (c.head + c.len) % c.cap
}

// Insert implements queue.Queue.  It fails and returns queue.ErrFull if
// there is not enough capacity in the circular buffer to accommodate
// the given element.
func (c *CBuf) Insert(e interface{}) error {
	if c.len == c.cap {
		return queue.ErrFull
	}
	c.elems[c.tail()] = e
	c.len++
	return nil
}

// Extract implements queue.Queue.
func (c *CBuf) Extract() (interface{}, error) {
	if c.len == 0 {
		return nil, queue.ErrEmpty
	}
	ret := c.elems[c.head]
	if c.len > 1 {
		c.head = c.next(c.head)
	}
	c.len--
	return ret, nil
}
