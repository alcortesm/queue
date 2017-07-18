/*
Package queue contains an interface and implementations for queues of
elements of the empty interface type.
*/
package queue

import "errors"

// A Queue represents an interface for a collection of elements of the
// empty interface type with queue access semantics.
//
// This interface makes no assumtions about the thread safety of its
// implementations.
type Queue interface {
	// Len returns the number of elements currently in the queue.
	Len() int
	// IsEmpty returns if the queue has no elements.
	IsEmpty() bool
	// Enqueue tries to add the given element at the back of the queue
	// and returns a nil error on success or an error on failure.
	Enqueue(interface{}) error
	// Head returns the first element in the queue and a nil error on
	// success.  It returns nil and an error on failure.  The error will
	// be ErrEmpty if the queue was originally empty.
	Head() (interface{}, error)
	// Dequeue extracts and returns the first element in the queue and a
	// nil error on success.  It returns nil and an error on failure.
	// The error will be ErrEmpty if the queue was originally empty.
	Dequeue() (interface{}, error)
}

var (
	// ErrEmpty is returned by the Head and Dequeue methods if the queue
	// was originally empty.
	ErrEmpty = errors.New("empty queue")
	// ErrFull is returned by Enqueue when there in not enough resources
	// to accommodate the given element.
	ErrFull = errors.New("full queue")
)
