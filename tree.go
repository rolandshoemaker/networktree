package main

import (
	"net"
	"sync"
)

// Tree is a binary tree that stores network prefixes
type Tree struct {
	mu *sync.RWMutex // only used at the root node

	children [2]*Tree

	value interface{}
}

func popcount(i uint8) int {
	i -= ((i >> 1) & 0x55)
	i = (i & 0x33) + ((i >> 2) & 0x33)
	return int((i + (i >> 4)) & 0x0F)
}

// New returns an initialized Tree
func New() *Tree {
	return &Tree{mu: new(sync.RWMutex)}
}

// Insert adds a network and the associated value to the tree. The passed
// net.IPNet is assumed to be valid.
func (t *Tree) Insert(n *net.IPNet, value interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	maskLen := 0
	for _, b := range n.Mask {
		maskLen += popcount(b)
	}

	node := t
	for i := 0; i < maskLen; i++ {
		bit := (n.IP[i/8] >> uint8(7-i%8)) & 0x01
		c := node.children[bit]
		if c == nil {
			c = &Tree{}
			node.children[bit] = c
		}
		node = c
	}

	node.value = value
	return
}

// Lookup looks up the value associated with the most specific network that
// contains the passed net.IP. If no applicable network is found nil is
// returned.
func (t *Tree) Lookup(ip net.IP) interface{} {
	if ip == nil {
		return nil
	}
	if x := ip.To4(); x != nil {
		ip = x
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	node := t
	var longestPrefix *Tree
	for i := 0; i < len(ip)*8; i++ {
		bit := (ip[i/8] >> uint8(7-i%8)) & 0x01
		child := node.children[bit]

		if child == nil {
			break
		}
		node = child
		if node.value != nil {
			longestPrefix = node
		}
	}

	if longestPrefix == nil {
		return nil
	}
	return longestPrefix.value
}

// path compression would probably be nice...
// func (t *Tree) Compact() {
// }
