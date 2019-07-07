package goqueue

import (
	"math"
	"sync"
	"sync/atomic"
)

const (
	upScale     = 2
	downScale   = 2
	downLimit   = 1024
	defaultSize = 30
)

// A ring-buffer queue is a FIFO(first-in-first-out) queue which is based
// on a resizable ring-buffer.
type ringBufferQueue struct {
	sync.Mutex
	buf      []interface{}
	head     uint64
	tail     uint64
	length   uint64
	capacity uint64
}

// NewRingBufferQueue creates a new ring-buffer queue.
//
// The param 'initSize' must be in the interval of [10, MaxUint16], otherwise it will
// be adjusted to a default value.
func NewRingBufferQueue(initSize uint64) Queue {
	if initSize < 10 || initSize > math.MaxUint16 {
		initSize = defaultSize
	}

	return &ringBufferQueue{
		buf:      make([]interface{}, initSize),
		capacity: initSize,
	}
}

func (q *ringBufferQueue) Enqueue(item interface{}) {
	q.Lock()
	if q.length == q.capacity {
		// The ring-buffer is full, resize it.
		newCapacity := q.capacity * upScale
		newBuf := make([]interface{}, newCapacity)
		for i := uint64(0); i < q.length; i++ {
			index := (q.tail + i) % q.capacity
			newBuf[i] = q.buf[index]
		}
		q.tail = 0
		q.head = q.length
		q.buf = newBuf
		atomic.StoreUint64(&q.capacity, newCapacity)
	}

	q.buf[q.head] = item
	q.head = (q.head + 1) % q.capacity
	atomic.AddUint64(&q.length, 1)
	q.Unlock()
}

func (q *ringBufferQueue) Dequeue() (interface{}, bool) {
	q.Lock()
	if q.length == 0 {
		return nil, false
	}

	item := q.buf[q.tail]
	q.tail = (q.tail + 1) % q.capacity
	atomic.AddUint64(&q.length, ^uint64(0))

	// The
	if q.capacity > downLimit && q.length < q.capacity/downScale {
		newCapacity := q.capacity / downScale
		newBuf := make([]interface{}, newCapacity)
		for i := uint64(0); i < q.length; i++ {
			index := (q.tail + i) % q.capacity
			newBuf[i] = q.buf[index]
		}
		q.tail = 0
		q.head = q.length
		q.buf = newBuf
		atomic.StoreUint64(&q.capacity, newCapacity)
	}

	q.Unlock()
	return item, true
}

func (q *ringBufferQueue) Peek() (interface{}, bool) {
	q.Lock()
	if q.length == 0 {
		return nil, false
	}
	item := q.buf[q.tail]
	q.Unlock()
	return item, true
}

func (q *ringBufferQueue) Length() uint64 {
	return atomic.LoadUint64(&q.length)
}

func (q *ringBufferQueue) Capacity() uint64 {
	return atomic.LoadUint64(&q.capacity)
}

func (q *ringBufferQueue) IsEmpty() bool {
	return q.Length() == 0
}
