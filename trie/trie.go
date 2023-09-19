package trie

import (
	"sync"
)

type node struct {
	children map[int]*node
	isEnd    bool   // 标记当前节点是否为单词的结尾
	word     string // 当前节点对应的单词
}

type Trie struct {
	root  *node      // Trie的根节点
	mutex sync.Mutex // 添加互斥锁
}

// NewTrie 创建一个新的Trie实例
func NewTrie() *Trie {
	ret := &Trie{
		root: &node{
			children: make(map[int]*node),
		},
	}

	return ret
}

// Insert 将一个单词插入到Trie中
func (t *Trie) Insert(word string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	currentNode := t.root
	for _, char := range word {
		idx := int(char)
		if _, ok := currentNode.children[idx]; !ok {
			currentNode.children[idx] = &node{
				children: make(map[int]*node),
			} // 如果子节点不存在，则创建一个新的节点
		}
		currentNode = currentNode.children[idx] // 移动到下一个节点
	}
	currentNode.isEnd = true // 标记当前节点为单词的结尾
	currentNode.word = word  // 记录当前节点对应的单词
}

// Delete 将一个单词从Trie中删除
func (t *Trie) Delete(word string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if find, node := t.findNode(word); find {
		node.isEnd = false
		node.word = ""
	}
}

// findNode 根据给定的单词在Trie中查找对应的节点
func (t *Trie) findNode(word string) (bool, *node) {
	currentNode := t.root
	var ok bool

	for _, char := range word {
		idx := int(char)
		currentNode, ok = currentNode.children[idx]
		//fmt.Println("char:", string(char), " idx:", idx, " ok:", ok)
		if !ok {
			return false, nil // 如果子节点不存在，则单词不存在于Trie中
		}
	}

	return true, currentNode // 返回单词是否存在以及对应的节点
}

// Search 检查一个单词是否存在于Trie中
func (t *Trie) Search(word string) bool {
	if find, node := t.findNode(word); find {
		return node.isEnd // 如果单词存在
	}
	return false
}

// StartWith 返回Trie中以给定前缀开头的所有单词
func (t *Trie) StartWith(word string) (r []string) {
	if find, node := t.findNode(word); find {
		return node.allValidChildren() // 返回以给定前缀开头的所有单词
	}
	return
}

// allValidChildren 返回当前节点及其子节点中的所有有效单词
func (t *node) allValidChildren() (r []string) {
	if t.isEnd {
		r = append(r, t.word) // 如果当前节点为单词的结尾，则将其加入结果列表
	}

	for _, child := range t.children {
		if nil != child {
			r = append(r, child.allValidChildren()...) // 递归获取子节点中的所有有效单词
		}
	}

	return r
}
