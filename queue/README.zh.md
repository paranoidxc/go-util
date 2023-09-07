# 队列

该包提供了一个通用的队列数据结构的Go实现。

**注意：该队列实现是线程安全的。**

## 使用方法

要使用该包，请将其导入到您的Go代码中：

```go
import "github.com/paranoidxc/go-util/queue
```

## 例子

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

## 函数

### New

```go
func New[T any](optBufSize int, optMinSize int) *Queue[T]
```

New函数创建一个具有指定缓冲区大小和最小大小的Queue数据结构的新实例。如果缓冲区大小小于0，则将其设置为默认缓冲区大小。如果最小大小小于0，则也将其设置为默认缓冲区大小。缓冲区大小将调整为大于或等于最小大小的下一个2的幂。返回值是指向创建的Queue的指针。

### Len

```go
func (q *Queue[T]) Len() int
```

Len函数返回队列中的元素数量。如果队列为nil，则返回0。

### Printf

```go
func (q *Queue[T]) Printf()
```

Printf函数打印队列中缓冲区的内容。它每行打印一个元素。如果队列为空，则打印一个空行。

### Push

```go
func (q *Queue[T]) Push(item T)
```

Push函数将一个元素添加到队列的末尾。如果缓冲区已满，则在添加元素之前扩展缓冲区。

### Pop

```go
func (q *Queue[T]) Pop() (T, error)
```

Pop函数移除并返回队列前端的元素。如果队列为空，则返回一个错误。

### Peek

```go
func (q *Queue[T]) Peek() (T, error)
```

Peek函数返回队列前端的元素，但不将其移除。如果队列为空，则返回一个错误。

### UnShift

```go
func (q *Queue[T]) UnShift(item T)
```

UnShift函数将一个元素添加到队列的前端。如果缓冲区已满，则在添加元素之前扩展缓冲区。

## 内部函数

这些函数是Queue数据结构内部使用的，不建议直接使用。

- `expandBuffIfNeed`：如果缓冲区已满，则扩展缓冲区。
- `shrink`：如果缓冲区的大小超过元素数量的4倍，则缩小缓冲区。
- `resize`：调整缓冲区的大小以容纳更多元素。
- `prev`：返回缓冲区中的前一个索引。
- `next`：返回缓冲区中的下一个索引。
