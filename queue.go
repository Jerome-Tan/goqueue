// Package goqueue implements some high performance queues.
package goqueue

// Queue is a common interface for queues.
type Queue interface {
	// Enqueue adds an obj into the queue.
	Enqueue(obj interface{})

	// Dequeue remove an obj from the queue, and return the obj.
	// If the queue is empty, then this method will return a nil
	// object and a false.
	Dequeue() (interface{}, bool)

	// Peek will take a look at the next item which will be dequeued
	// from the queue, but will not the item.
	// If the queue is emtpy, then this method will return a nil
	// object and a false.
	Peek() (interface{}, bool)

	// Length returns the number of items in the queue.
	Length() uint64

	// Capacity returns the capacity of the queue.
	Capacity() uint64

	// IsEmpty return if the queue is empty.
	IsEmpty() bool
}
