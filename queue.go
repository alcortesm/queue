/*
Package queue contains an interface and some implementations for a FIFO
collection of empty interfaces.
*/
package queue

import "errors"

// A Queue interface represents a First-In-First-Out (FIFO) collection
// of the empty interfaces.
//
// This interface makes no assumtions about the thread safety of its
// implementations.
type Queue interface {
	// Insert appends the given element at the end of the queue and
	// returns nil on success or an error on failure.
	Insert(interface{}) error
	// Extract removes and returns the first element in the queue and
	// nil on success.  On failure, it returns nil and an error.
	Extract() (interface{}, error)
}

var (
	// ErrEmpty is returned by Extract if the queue is empty.
	ErrEmpty = errors.New("empty queue")
	// ErrFull is returned by Insert when there are not
	// enough resources to accommodate the given element.
	ErrFull = errors.New("full queue")
)
