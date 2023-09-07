# go-util


# Queue Usage

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
