package lru

import (
	"sync"
)

type Key string

type ListManipulator interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type LruCache struct {
	sync.RWMutex
	capacity int
	queue    *List
	items    map[Key]*ListItem
}

func NewCache(capacity int) *LruCache {
	return &LruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

type pair struct {
	key   Key
	value interface{}
}

// Add a value to the cache by key,
// if it already exists, replace it with a new value.
func (c *LruCache) Set(key Key, value interface{}) bool {
	c.Lock()
	defer c.Unlock()

	node, ok := c.items[key]

	if ok {
		node.Value = pair{key: key, value: value}
		c.queue.MoveToFront(node)
		return true
	}

	kv := pair{key: key, value: value}

	if c.queue.length == c.capacity {
		valueTodelete := c.queue.tail
		c.queue.Remove(valueTodelete)

		val, ok := valueTodelete.Value.(pair)
		if !ok {
			panic("Unexpected val type")
		}

		delete(c.items, val.key)
	}

	li := c.queue.PushFront(kv)
	c.items[key] = li

	return false
}

// Get a value from the cache by key.
func (c *LruCache) Get(key Key) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	node, ok := c.items[key]

	if ok {
		c.queue.MoveToFront(node)
		val, ok := node.Value.(pair)
		if !ok {
			panic("Unexpected value type")
		}
		return val.value, true
	}
	return nil, false
}

func (c *LruCache) Clear() {
	c.Lock()
	defer c.Unlock()

	c.items = map[Key]*ListItem{}
	c.queue = NewList()
}

func (c *LruCache) PrintList() {
	c.queue.printList()
}
