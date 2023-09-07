package queue

import (
	"errors"
	"fmt"
	"sync"
)

const defBufferSize = 256

// Queue represents a generic queue data structure.
type Queue[T any] struct {
	buffer  []T
	cnt     int
	tail    int
	head    int
	minSize int
	mutex   sync.Mutex // 添加互斥锁
}

// New creates a new instance of the Queue data structure with the specified buffer size and minimum size.
// If the buffer size is less than 0, it will be set to the default buffer size.
// If the minimum size is less than 0, it will also be set to the default buffer size.
// The buffer size will be adjusted to the next power of 2 greater than or equal to the minimum size.
// The returned value is a pointer to the created Queue.
func New[T any](optBufSize int, optMinSize int) *Queue[T] {
	if optBufSize < 0 {
		optBufSize = defBufferSize
	}
	if optMinSize < 0 {
		optMinSize = defBufferSize
	}

	minSize := defBufferSize
	for minSize < optMinSize {
		minSize <<= 1
	}

	var buf []T
	if optBufSize != 0 {
		pow2ToBufSize := minSize
		for pow2ToBufSize < optBufSize {
			pow2ToBufSize <<= 1
		}
		buf = make([]T, pow2ToBufSize)
	}

	return &Queue[T]{
		buffer:  buf,
		minSize: minSize,
	}
}

// Len returns the number of elements in the queue.
// If the queue is nil, it returns 0.
func (q *Queue[T]) Len() int {
	if nil == q {
		return 0
	}

	return q.cnt
}

// Printf prints the content of the buffer in the queue.
// It prints each element on a new line.
// If the queue is empty, it prints an empty line.
func (q *Queue[T]) Printf() {
	fmt.Println("buffer content [")
	if q.cnt > 0 {
		start := q.head
		for i := 0; i < q.cnt; i++ {
			fmt.Printf("%v\n", q.buffer[start])
			start = q.next(start)
		}
	}
	fmt.Println("]")
}

// Push adds an item to the end of the queue.
// If the buffer is full, it expands the buffer before adding the item.
func (q *Queue[T]) Push(item T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.expandBuffIfNeed()

	q.buffer[q.tail] = item
	q.tail = q.next(q.tail)

	q.cnt++
}

// Pop removes and returns the item at the front of the queue.
// If the queue is empty, it returns an error.
func (q *Queue[T]) Pop() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.cnt <= 0 {
		var empty T
		return empty, errors.New("Pop() Err: queue is empty")
	}
	ret := q.buffer[q.head]

	var empty T
	q.buffer[q.head] = empty
	q.head = q.next(q.head)
	q.cnt--

	q.shrink()
	return ret, nil
}

// Peek returns the item at the front of the queue without removing it.
// If the queue is empty, it returns an error.
func (q *Queue[T]) Peek() (T, error) {
	if q.cnt <= 0 {
		var empty T
		return empty, errors.New("Peek() Err: queue is empty")
	}
	ret := q.buffer[q.head]
	return ret, nil
}

// UnShift adds an item to the front of the queue.
// If the buffer is full, it expands the buffer before adding the item.
func (q *Queue[T]) UnShift(item T) {
	q.expandBuffIfNeed()

	q.head = q.prev(q.head)
	q.buffer[q.head] = item
	q.cnt++
}

// expandBuffIfNeed expands the buffer if it is full.
func (q *Queue[T]) expandBuffIfNeed() {
	if q.cnt != len(q.buffer) {
		return
	}
	if len(q.buffer) == 0 {
		if 0 == q.minSize {
			q.minSize = defBufferSize
		}
		q.buffer = make([]T, q.minSize)
		return
	}

	q.resize()
}

// shrink shrinks the buffer if it is more than 4 times larger than the number of elements.
func (q *Queue[T]) shrink() {
	if len(q.buffer) > q.minSize && (q.cnt<<2) == len(q.buffer) {
		q.resize()
	}
}

// resize resizes the buffer to accommodate more elements.
func (q *Queue[T]) resize() {
	newBuf := make([]T, q.cnt<<1)
	if q.tail > q.head {
		copy(newBuf, q.buffer[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buffer[q.head:])
		copy(newBuf[n:], q.buffer[:q.tail])
	}

	q.head = 0
	q.tail = q.cnt
	q.buffer = newBuf
}

// prev returns the previous index in the buffer.
func (q *Queue[T]) prev(i int) int {
	return (i - 1) & (len(q.buffer) - 1)
}

// next returns the next index in the buffer.
func (q *Queue[T]) next(i int) int {
	return (i + 1) & (len(q.buffer) - 1)
}
