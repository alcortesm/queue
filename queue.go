/*
Package queue contains the interface and definitions for the queues
implemented in this project.
*/
package queue

import "errors"

// A Queue represents a collection of values of type interface{} with a
// queue access semantics.
//
// Elements can be inserted at the back of the queue using the Enqueue
// method.  The first element in the queue can be accessed through the
// Head method or extracted using the Dequeue method.
//
// There are two types of queues, bounded and infinite.  Bounded queues
// have a constant capacity representing the maximum number of elements
// they can store at any given time.  Infinite queues don't have the
// concept of capacity, you will be able to enqueue elements in them
// until they eventually run out of resources.
//
// If a bounded queue has capacity 0, this means it will be empty and
// full at the same time: all enqueue calls will fail because there is
// not enough capacity to store the given element.  At the same time,
// all dequeue calls will fail because there aren't any elements in the
// queue to retrieve.
//
// This interface makes no assumtions about the thread safety of its
// implementations.
type Queue interface {
	// IsBounded returns true if the queue has fixed capacity and false if
	// it is infinite.
	IsBounded() bool
	// Cap returns the capacity of the queue and nil for bounded queues.
	// On infinite queues it returns 0 and ErrInfinite.  The capacity of
	// bounded queues will be 0 or a positive integer.
	Cap() (int, error)
	// Len returns the number of elements currently in the queue.
	Len() int
	// IsEmpty returns if the queue is empty, this is, if Len is 0.
	IsEmpty() bool
	// IsFull returns if there is not enough capacity in the queue to
	// store more elements.  Infinite queues always return false.
	IsFull() bool
	// Enqueue try to add the given element at the back of the queue and
	// returns a nil error on success.  It returns ErrFull on a full
	// bounded queue.  Infinite queues should not panic if they cannot enqueue
	// due to lack of resources, instead a sensible error from the underliying
	// medium should be returned.
	Enqueue(interface{}) error
	// Head returns the first element in the queue and a nil error if
	// the queue is not empty.  Otherwise It returns nil and ErrEmpty.
	Head() (interface{}, error)
	// Dequeue extracts and returns the first element from the queue and a
	// nil error if the queue is not empty.  Otherwise it returns nil
	// and ErrEmpty.
	Dequeue() (interface{}, error)
}

var (
	// ErrInfinite is returned by Cap if the queue is infinite.
	ErrInfinite = errors.New("infinite queue")
	// ErrFull is returned when trying to enqueue to a full bounded queue.
	ErrFull = errors.New("full queue")
	// ErrEmpty is returned by Head and Dequeue if the queue is empty.
	ErrEmpty = errors.New("empty queue")
)
