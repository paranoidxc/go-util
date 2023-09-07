# go-util

# Queue

This package provides a generic implementation of a queue data structure in Go.

For more detail see [Queue README](./queue/README.md)

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
