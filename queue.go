/*
Package queue contains the interface and definitions for the queues
implemented in this project.
*/
package queue

// A Queue represents a collection of values of type interface{} with a
// queue access semantics.
//
// Elements can be inserted at the back of the queue using the Enqueue
// method.  The first element in the queue can be accessed through the
// Head method or extracted using the Dequeue method.
//
// Queues can be bounded (fixed capacity) or infinite.  Bounded queues
// will have a constant size representing the maximum number of elements
// they can store at any given time.
//
// If a bounded queue has size 0, all enqueue calls will fail and the
// queue will be empty and full at the same time.
//
// Infinite queues can still fail to enqueue elements due to memory
// exhaustion.
//
// This interface makes no assumtions about the thread safety of its
// implementations.
type Queue interface {
	// Bounded returns true if the queue has fixed capacity and false if
	// it is infinite.
	Bounded() bool
	// Size returns the maximum number of elements allowed in the queue
	// at any given time and a nil error if the queue is bounded.  On
	// infinite queues it return 0 and ErrInfinite.  The size of bounded
	// queues will be 0 or a positve integer.
	Size() (int, error)
	// Len returns the number of elements currently in the queue.
	Len() int
	// Empty returns if the queue is empty, this is, if Len is 0.
	Empty() bool
	// Full returns if the queue is full, this is, if no more elements
	// can fit in the queue.  Infinite queues are never full, yet
	// enqueuing can fail due to memory exhaustion.
	Full() bool
	// Enqueue try to add the given element at the back of the queue and
	// returns a nil error on success.  Otherwise, it returns an error.
	// In particular, for bouned queues, ErrFull is returned if the call
	// failed because the queue was already full.
	Enqueue(interface{}) error
	// Head returns the first element in the queue and a nil error if
	// the queue is not empty.  Otherwise It returns nil and ErrEmpty.
	Head() (interface{}, error)
	// Dequeue extracts and returns the first element of the queue and a
	// nil error if the queue is not empty.  Otherwise it returns nil
	// and ErrEmpty.
	Dequeue() (interface{}, error)
}
