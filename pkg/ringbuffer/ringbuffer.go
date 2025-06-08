package ringbuffer

// Queue is a generic fixed-size queue
type Queue[T any] struct {
	items []T
	size  int
	max   int
}

// NewQueue creates a new generic fixed-size queue with the given maximum capacity
func NewQueue[T any](maxSize int) *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0, maxSize),
		size:  0,
		max:   maxSize,
	}
}

// Add adds an item to the queue, removing the oldest if at capacity
func (q *Queue[T]) Add(item T) {
	if q.size == q.max {
		// Remove the oldest item (first in the slice)
		q.items = q.items[1:]
		q.size--
	}
	q.items = append(q.items, item)
	q.size++
}

func (q *Queue[T]) Len() int {
	return len(q.items)
}

// GetAll returns all items in the queue
func (q *Queue[T]) GetAll() []T {
	return q.items
}

func (q *Queue[T]) GetLastN(n int) []T {
	if n > q.size {
		n = q.size
	}
	return q.items[q.size-n : q.size]
}
