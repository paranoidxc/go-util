# Queue

This package provides a generic implementation of a queue data structure in Go.

**Note: The Queue implementation in this package is thread-safe.**

## Usage

To use this package, import it into your Go code:

```go
import "github.com/paranoidxc/go-util/queue"
```

## Example

```go
// init
q := queue.New[string](1024, 32)
var q queue.Queue[string]

// put
q.Push("World")
q.Unshift("Hello")

// get
q.Pop()

// len
q.Len()

// fmt.Printf
q.Pringf()
```

## Functions

### New

```go
func New[T any](optBufSize int, optMinSize int) *Queue[T]
```

New creates a new instance of the Queue data structure with the specified buffer size and minimum size. If the buffer size is less than 0, it will be set to the default buffer size. If the minimum size is less than 0, it will also be set to the default buffer size. The buffer size will be adjusted to the next power of 2 greater than or equal to the minimum size. The returned value is a pointer to the created Queue.

### Len

```go
func (q *Queue[T]) Len() int
```

Len returns the number of elements in the queue. If the queue is nil, it returns 0.

### Printf

```go
func (q *Queue[T]) Printf()
```

Printf prints the content of the buffer in the queue. It prints each element on a new line. If the queue is empty, it prints an empty line.

### Push

```go
func (q *Queue[T]) Push(item T)
```

Push adds an item to the end of the queue. If the buffer is full, it expands the buffer before adding the item.

### Pop

```go
func (q *Queue[T]) Pop() (T, error)
```

Pop removes and returns the item at the front of the queue. If the queue is empty, it returns an error.

### Peek

```go
func (q *Queue[T]) Peek() (T, error)
```

Peek returns the item at the front of the queue without removing it. If the queue is empty, it returns an error.

### UnShift

```go
func (q *Queue[T]) UnShift(item T)
```

UnShift adds an item to the front of the queue. If the buffer is full, it expands the buffer before adding the item.

## Internal Functions

These functions are used internally by the Queue data structure and are not intended to be used directly.

- `expandBuffIfNeed`: Expands the buffer if it is full.
- `shrink`: Shrinks the buffer if it is more than 4 times larger than the number of elements.
- `resize`: Resizes the buffer to accommodate more elements.
- `prev`: Returns the previous index in the buffer.
- `next`: Returns the next index in the buffer.
