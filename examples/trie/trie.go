package main

import (
	"fmt"

	"github.com/paranoidxc/go-util/trie"
)

func main() {
	testTrie := trie.NewTrie()
	toAdd := []string{
		"aragorn",
		"arat",
		"aragon",
		"argon",
		"eragon",
		"oregon",
		"oregano",
		"oreo",
	}
	for _, word := range toAdd {
		testTrie.Insert(word)
	}

	fmt.Println(testTrie.Search("eragon"))
	testTrie.Delete("eragon")
	fmt.Println(testTrie.Search("eragon"))
	fmt.Println(testTrie.Search("wizard"))

	findRes := testTrie.StartWith("ara")
	fmt.Println("find len:", len(findRes), "find Res:", findRes)

	findRes = testTrie.StartWith("wizard")
	fmt.Println("find len:", len(findRes), "find Res:", findRes)
}
