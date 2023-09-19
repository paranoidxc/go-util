package trie

import (
	"fmt"
	"testing"
)

func TestTrie_Insert(t *testing.T) {
	trie := NewTrie()

	trie.Insert("apple")
	trie.Insert("banana")
	trie.Insert("orange")

	// Check if words are inserted correctly
	if !trie.Search("apple") {
		t.Error("Failed to insert word: apple")
	}
	if !trie.Search("banana") {
		t.Error("Failed to insert word: banana")
	}
	if !trie.Search("orange") {
		t.Error("Failed to insert word: orange")
	}
}

func TestTrie_Search(t *testing.T) {
	trie := NewTrie()

	trie.Insert("apple")
	trie.Insert("banana")
	trie.Insert("orange")

	// Check if existing words are found
	if !trie.Search("apple") {
		t.Error("Failed to find word: apple")
	}
	if !trie.Search("banana") {
		t.Error("Failed to find word: banana")
	}
	if !trie.Search("orange") {
		t.Error("Failed to find word: orange")
	}

	// Check if non-existing words are not found
	if trie.Search("grape") {
		t.Error("Incorrectly found word: grape")
	}
	if trie.Search("melon") {
		t.Error("Incorrectly found word: melon")
	}
}

func TestTrie_StartWith(t *testing.T) {
	trie := NewTrie()

	trie.Insert("apple")
	trie.Insert("banana")
	trie.Insert("orange")

	// Check if words starting with a prefix are returned correctly
	words := trie.StartWith("app")
	expected := []string{"apple"}
	if !equalSlice(words, expected) {
		t.Errorf("Incorrect words returned for prefix: app. Expected: %v, Got: %v", expected, words)
	}

	words = trie.StartWith("ban")
	expected = []string{"banana"}
	if !equalSlice(words, expected) {
		t.Errorf("Incorrect words returned for prefix: ban. Expected: %v, Got: %v", expected, words)
	}

	words = trie.StartWith("or")
	expected = []string{"orange"}
	if !equalSlice(words, expected) {
		t.Errorf("Incorrect words returned for prefix: or. Expected: %v, Got: %v", expected, words)
	}

	// Check if empty slice is returned for non-existing prefix
	words = trie.StartWith("gr")
	if len(words) != 0 {
		t.Errorf("Incorrect words returned for non-existing prefix: gr. Expected: [], Got: %v", words)
	}
}

func BenchmarkTrie_Insert(b *testing.B) {
	trie := NewTrie()

	for i := 0; i < b.N; i++ {
		word := fmt.Sprintf("word%d", i)
		trie.Insert(word)
	}
}

func BenchmarkTrie_Search(b *testing.B) {
	trie := NewTrie()

	for i := 0; i < b.N; i++ {
		word := fmt.Sprintf("word%d", i)
		trie.Insert(word)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		word := fmt.Sprintf("word%d", i)
		trie.Search(word)
	}
}

func equalSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
