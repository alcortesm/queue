/*
Package queue contains an interface and some implementations for queues
of elements of the empty interface type.
*/
package queue

import "errors"

// A Queue interface represents a queue of elements of the empty
// interface type.
//
// This interface makes no assumtions about the thread safety of its
// implementations.
type Queue interface {
	// Enqueue tries to add the given element at the back of the queue
	// and returns a nil error on success or an error on failure.  The
	// error will be ErrFull if the queue ran out of resources to
	// accommodate the given element.
	Enqueue(interface{}) error
	// Dequeue extracts and returns the first element in the queue and a
	// nil error on success.  On failure, it returns nil and an error.
	// The error will be ErrEmpty if the queue was originally empty.
	Dequeue() (interface{}, error)
}

var (
	// ErrEmpty is returned by Dequeue if the queue is empty.
	ErrEmpty = errors.New("empty queue")
	// ErrFull is returned by the Enqueue method when there is not
	// enough resources to accommodate the given element.
	ErrFull = errors.New("full queue")
)
