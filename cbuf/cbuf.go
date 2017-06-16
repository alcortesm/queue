/*
Package cbuf implements a bounded Queue interface using a fixed
capacity, in memory, circular buffer.

This implementation is not thread-safe.
*/
package cbuf

import (
	"fmt"

	"github.com/alcortesm/queue"
)

type CBuf struct {
	size  int
	len   int
	elems []interface{}
	head  int
}

// New returns a new CBuf implementing a bounded queue.Queue of the
// given size and a nil error.  It returns nil and an error if size is
// negative.
//
// The allocation of the memory to store the elements is delayed until
// the first enqueue operation.
func New(size int) (queue.Queue, error) {
	if size < 0 {
		return nil, fmt.Errorf("size must be 0 or positive, was %d", size)
	}
	return &CBuf{size: size}, nil
}

func (c *CBuf) lazyElems() []interface{} {
	if c.elems == nil {
		c.elems = make([]interface{}, c.size)
	}
	return c.elems
}

// Implements queue.Queue.  CBuf queues are always bounded, which means this
// method always return true.
func (c *CBuf) Bounded() bool {
	return true
}

// Implements queue.Queue.  CBuf queues are always bounded, which means
// this method always return the size and a nil error.
func (c *CBuf) Size() (int, error) {
	return c.size, nil
}

// Implements queue.Queue.
func (c *CBuf) Len() int {
	return c.len
}

// Implements queue.Queue.
func (c *CBuf) Empty() bool {
	return c.len == 0
}

// Implements queue.Queue.
func (c *CBuf) Full() bool {
	return c.len == c.size
}

func (c *CBuf) next(n int) int {
	return (n + 1) % c.size
}

func (c *CBuf) tail() int {
	return (c.head + c.len - 1) % c.size
}

// Implements queue.Queue.
func (c *CBuf) Enqueue(e interface{}) error {
	if c.len == c.size {
		return queue.ErrFull
	}
	elems := c.lazyElems()
	elems[c.tail()] = e
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
	c.elems[c.head] = nil
	c.len--
	return ret, nil
}
