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

type cbuf struct {
	size  int
	len   int
	elems []interface{}
	head  int
}

func New(size int) (queue.Queue, error) {
	if size < 0 {
		return nil, fmt.Errorf("size must be 0 or positive, was %d", size)
	}
	return &cbuf{size: size}, nil
}

func (c *cbuf) lazyElems() []interface{} {
	if c.elems == nil {
		c.elems = make([]interface{}, c.size)
	}
	return c.elems
}

func (c *cbuf) Bounded() bool {
	return true
}

func (c *cbuf) Size() (int, error) {
	return c.size, nil
}

func (c *cbuf) Len() int {
	return c.len
}

func (c *cbuf) Empty() bool {
	return c.len == 0
}

func (c *cbuf) Full() bool {
	return c.len == c.size
}

func (c *cbuf) next(n int) int {
	return (n + 1) % c.size
}

func (c *cbuf) tail() int {
	return (c.head + c.len - 1) % c.size
}

func (c *cbuf) Enqueue(e interface{}) error {
	if c.len == c.size {
		return queue.ErrFull
	}
	elems := c.lazyElems()
	elems[c.tail()] = e
	c.len++
	return nil
}

func (c *cbuf) Head() (interface{}, error) {
	if c.len == 0 {
		return nil, queue.ErrEmpty
	}
	return c.elems[c.head], nil
}

func (c *cbuf) Dequeue() (interface{}, error) {
	if c.len == 0 {
		return nil, queue.ErrEmpty
	}
	ret := c.elems[c.head]
	c.elems[c.head] = nil
	c.len--
	return ret, nil
}
