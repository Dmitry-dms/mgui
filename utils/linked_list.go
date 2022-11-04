package utils

import (
	"fmt"
	"sync"
)

// Node a single node that composes the list
type Node[T any] struct {
	Content T
	Next    *Node[T]
	Prev    *Node[T]
}

// LinkedList the linked list of Items
type LinkedList[T any] struct {
	head *Node[T]
	size int
	lock sync.RWMutex
}

func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

// Append adds an Item to the end of the linked list
func (ll *LinkedList[T]) Append(t T) {
	ll.lock.Lock()
	node := Node[T]{t, nil, nil}
	if ll.head == nil {
		ll.head = &node
	} else {
		last := ll.head
		for {
			if last.Next == nil {
				node.Prev = last
				break
			}
			last = last.Next
		}
		last.Next = &node
	}
	ll.size++
	ll.lock.Unlock()
}

// Insert adds an Item at position i
func (ll *LinkedList[T]) Insert(i int, t T) error {
	ll.lock.Lock()
	defer ll.lock.Unlock()
	if i < 0 || i > ll.size {
		return fmt.Errorf("Index out of bounds")
	}
	addNode := Node[T]{t, nil, nil}
	if i == 0 {
		addNode.Next = ll.head
		ll.head = &addNode
		return nil
	}
	node := ll.head
	j := 0
	for j < i-2 {
		j++
		node = node.Next
	}
	addNode.Next = node.Next
	node.Next = &addNode
	ll.size++
	return nil
}

// RemoveAt removes a node at position i
func (ll *LinkedList[T]) RemoveAt(i int) (*T, error) {
	ll.lock.Lock()
	defer ll.lock.Unlock()
	if i < 0 || i > ll.size {
		return nil, fmt.Errorf("Index out of bounds")
	}
	node := ll.head
	j := 0
	for j < i-1 {
		j++
		node = node.Next
	}
	remove := node.Next
	node.Next = remove.Next
	ll.size--
	return &remove.Content, nil
}

// IndexOf returns the position of the Item t
// func (ll *LinkedList[T]) IndexOf(t T) int {
// 	ll.lock.RLock()
// 	defer ll.lock.RUnlock()
// 	node := ll.head
// 	j := 0
// 	for {
// 		if node.content == t {
// 			return j
// 		}
// 		if node.next == nil {
// 			return -1
// 		}
// 		node = node.next
// 		j++
// 	}
// }

// IsEmpty returns true if the list is empty
func (ll *LinkedList[T]) IsEmpty() bool {
	ll.lock.RLock()
	defer ll.lock.RUnlock()
	if ll.head == nil {
		return true
	}
	return false
}

// Size returns the linked list size
func (ll *LinkedList[T]) Size() int {
	ll.lock.RLock()
	defer ll.lock.RUnlock()
	size := 1
	last := ll.head
	for {
		if last == nil || last.Next == nil {
			break
		}
		last = last.Next
		size++
	}
	return size
}

// Head returns a pointer to the first node of the list
func (ll *LinkedList[T]) Head() *Node[T] {
	ll.lock.RLock()
	defer ll.lock.RUnlock()
	return ll.head
}
