/*
Package cbuf implements a bounded Queue interface using a fixed
capacity, in memory, circular buffer that gets allocated upon
construction.

This implementation is not thread-safe.
*/
package cbuf

import (
	"fmt"

	"github.com/alcortesm/queue"
)

type CBuf struct {
	cap   int
	len   int
	elems []interface{}
	head  int
}

// New returns a new CBuf implementing a bounded queue.Queue of the
// given capacity and a nil error.  It returns nil and an error if the
// given capacity is a negative integer.
func New(capacity int) (queue.Queue, error) {
	if capacity < 0 {
		return nil, fmt.Errorf("capacity must be 0 or positive, was %d",
			capacity)
	}
	return &CBuf{
		cap:   capacity,
		elems: make([]interface{}, capacity),
	}, nil
}

// Implements queue.Queue.  CBuf queues are always bounded, which means this
// method always return true.
func (c *CBuf) IsBounded() bool {
	return true
}

// Implements queue.Queue.
func (c *CBuf) Cap() (int, error) {
	return c.cap, nil
}

// Implements queue.Queue.
func (c *CBuf) Len() int {
	return c.len
}

// Implements queue.Queue.
func (c *CBuf) IsEmpty() bool {
	return c.len == 0
}

// Implements queue.Queue.
func (c *CBuf) IsFull() bool {
	return c.len == c.cap
}

func (c *CBuf) next(n int) int {
	return (n + 1) % c.cap
}

func (c *CBuf) tail() int {
	return (c.head + c.len) % c.cap
}

// Implements queue.Queue.
func (c *CBuf) Enqueue(e interface{}) error {
	if c.len == c.cap {
		return queue.ErrFull
	}
	c.elems[c.tail()] = e
	c.len++
	return nil
}

// Implements queue.Queue.
func (c *CBuf) Head() (interface{}, error) {
	if c.len == 0 {
		return nil, queue.ErrEmpty
	}
	return c.elems[c.head], nil
}

// Implements queue.Queue.
func (c *CBuf) Dequeue() (interface{}, error) {
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
