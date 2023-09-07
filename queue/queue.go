package queue

import (
	"errors"
	"fmt"
)

const defBufferSize = 256

type Queue[T any] struct {
	buffer  []T
	cnt     int
	tail    int
	head    int
	minSize int
}

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

func (q *Queue[T]) Len() int {
	if nil == q {
		return 0
	}

	return q.cnt
}

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

func (q *Queue[T]) Push(item T) {
	q.expandBuffIfNeed()

	q.buffer[q.tail] = item
	q.tail = q.next(q.tail)

	q.cnt++
}

func (q *Queue[T]) Pop() (T, error) {
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

func (q *Queue[T]) Peek() (T, error) {
	if q.cnt <= 0 {
		var empty T
		return empty, errors.New("Peek() Err: queue is empty")
	}
	ret := q.buffer[q.head]
	return ret, nil
}

func (q *Queue[T]) UnShift(item T) {
	q.expandBuffIfNeed()

	q.head = q.prev(q.head)
	q.buffer[q.head] = item
	q.cnt++
}

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

func (q *Queue[T]) shrink() {
	if len(q.buffer) > q.minSize && (q.cnt<<2) == len(q.buffer) {
		q.resize()
	}
}

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

func (q *Queue[T]) prev(i int) int {
	return (i - 1) & (len(q.buffer) - 1)
}

func (q *Queue[T]) next(i int) int {
	return (i + 1) & (len(q.buffer) - 1)
}
