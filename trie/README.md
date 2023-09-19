# Trie

Trie is a Go package that implements a trie data structure. A trie, also known as a prefix tree, is an efficient data structure for storing and searching strings.

## Installation

To use the Trie package, you need to have Go installed and set up on your machine. Then, you can install the package using the `go get` command:

```shell
go get github.com/paranoid.xc/go-util/trie
```

## Usage

Import the Trie package in your Go code:

```go
import "github.com/paranoid.xc/go-util/trie"
```

### Creating a Trie

To create a new Trie instance, use the `NewTrie` function:

```go
t := trie.NewTrie()
```

### Inserting Words

To insert a word into the Trie, use the `Insert` method:

```go
t.Insert("apple")
t.Insert("banana")
t.Insert("orange")
```

### Searching for Words

To search for a word in the Trie, use the `Search` method:

```go
found := t.Search("apple")
if found {
    fmt.Println("Word found!")
} else {
    fmt.Println("Word not found!")
}
```

### Finding Words with a Prefix

To find all words in the Trie that start with a given prefix, use the `StartWith` method:

```go
words := t.StartWith("app")
for _, word := range words {
    fmt.Println(word)
}
```

## Functions

### NewTrie

```go
func NewTrie() *Trie
```

NewTrie creates a new Trie instance.

### Insert

```go
func (t *Trie) Insert(word string)
```

Insert inserts a word into the Trie.

### Search

```go
func (t *Trie) Search(word string) bool
```

Search checks if a word exists in the Trie.

### StartWith

```go
func (t *Trie) StartWith(word string) []string
```

StartWith returns all words in the Trie that start with a given prefix.

## Internal Functions

### findNode

```go
func (t *Trie) findNode(word string) (bool, *node)
```

findNode finds a node in the Trie based on a given word.

### allValidChildren

```go
func (t *node) allValidChildren() []string
```

allValidChildren returns all valid words from the current node and its children.

## Contributing

Contributions to the Trie package are welcome. If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.
