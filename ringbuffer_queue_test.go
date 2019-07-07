package goqueue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRingBufferQueue_EnqueueDequeue(t *testing.T) {
	q := NewRingBufferQueue(10)

	q.Enqueue("data")
	item, ok := q.Dequeue()
	assert.Equal(t, "data", item)
	assert.True(t, ok)

	for i := 0; i < 10; i++ {
		q.Enqueue(i)
		assert.Equal(t, uint64(i+1), q.Length())
		assert.Equal(t, uint64(10), q.Capacity())
	}

	q.Enqueue(10)
	assert.Equal(t, uint64(11), q.Length())
	assert.Equal(t, uint64(20), q.Capacity())

	for i := 0; i < 11; i++ {
		assert.Equal(t, uint64(11-i), q.Length())

		item, ok = q.Dequeue()
		assert.True(t, ok)
		assert.Equal(t, i, item)
	}

	item, ok = q.Dequeue()
	assert.False(t, ok)
	assert.Nil(t, item)

	assert.Equal(t, uint64(0), q.Length())
	assert.Equal(t, uint64(20), q.Capacity())
}

func BenchmarkNewRingBufferQueue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewRingBufferQueue(defaultSize)
	}
}

func BenchmarkRingBufferQueue_Length(b *testing.B) {
	q := NewRingBufferQueue(defaultSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Length()
	}
}

func BenchmarkRingBufferQueue_IsEmpty_Yes(b *testing.B) {
	q := NewRingBufferQueue(defaultSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.IsEmpty()
	}
}

func BenchmarkRingBufferQueue_IsEmpty_No(b *testing.B) {
	q := NewRingBufferQueue(defaultSize)
	q.Enqueue("item")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.IsEmpty()
	}
}

func BenchmarkRingBufferQueue_Capacity(b *testing.B) {
	q := NewRingBufferQueue(defaultSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Capacity()
	}
}

func BenchmarkRingBufferQueue_Enqueue(b *testing.B) {
	q := NewRingBufferQueue(defaultSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
}

func BenchmarkRingBufferQueue_Dequeue(b *testing.B) {
	q := NewRingBufferQueue(defaultSize)
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Dequeue()
	}
}

func BenchmarkRingBufferQueue_Peek(b *testing.B) {
	q := NewRingBufferQueue(defaultSize)
	q.Enqueue("item")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Peek()
	}
}

func BenchmarkRingBufferQueue_EnqueueDequeue(b *testing.B) {
	q := NewRingBufferQueue(defaultSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
		q.Dequeue()
	}
}

func BenchmarkRingBufferQueue_EnqueueDequeue_Parallel(b *testing.B) {
	q := NewRingBufferQueue(defaultSize)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			q.Enqueue(i)
			q.Dequeue()
		}
	})
}
