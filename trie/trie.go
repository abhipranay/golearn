// Package main provides ...
package main

import "fmt"

type TrieNode struct {
	Char     byte
	Children map[byte]*TrieNode
	IsEnd    bool
}

func NewTrieNode(ch byte) *TrieNode {
	return &TrieNode{
		Char:     ch,
		Children: make(map[byte]*TrieNode, 0),
		IsEnd:    false,
	}
}

type Trie struct {
	Root *TrieNode
}

func NewTrie() Trie {
	return Trie{
		Root: NewTrieNode(byte('*')),
	}
}

func (t Trie) AddWord(word string) {
	if len(word) == 0 {
		return
	}

	t.add(t.Root, []byte(word), 0)
}

func (t Trie) add(root *TrieNode, word []byte, pos int) {
	if pos >= len(word) {
		root.IsEnd = true
		return
	}
	ch := word[pos]
	if node, ok := root.Children[ch]; ok {
		t.add(node, word, pos+1)
	} else {
		root.Children[ch] = NewTrieNode(ch)
		t.add(root.Children[ch], word, pos+1)
	}
}

func (t Trie) Search(word string) bool {
	return t.searchDotPattern(t.Root, []byte(word), 0)
}

func (t Trie) searchDotPattern(root *TrieNode, word []byte, pos int) bool {
	if pos == len(word) {
		return root.IsEnd
	}
	ch := word[pos]
	dot := byte('.')
	if ch != dot {
		if node, ok := root.Children[ch]; ok {
			return t.searchDotPattern(node, word, pos+1)
		}
		return false
	}
	for _, node := range root.Children {
		r := t.searchDotPattern(node, word, pos+1)
		if r {
			return true
		}
	}
	return false
}

func main() {
	trie := NewTrie()

	var words []string = []string{"bad", "bat", "battle", "apple", "cat"}

	for i := 0; i < len(words); i++ {
		trie.AddWord(words[i])
	}
	fmt.Println(trie.Search("bat"))
	fmt.Println(trie.Search("batman"))
	fmt.Println(trie.Search("battle"))
	fmt.Println(trie.Search("ba.t.e"))
	fmt.Println(trie.Search("ba.t."))
}
