package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Set value in cache.
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	listItemValue := &cacheItem{
		key:   key,
		value: value,
	}

	if item, exist := c.items[key]; exist {
		c.queue.MoveToFront(item)
		item.Value = listItemValue
		return true
	}
	if c.capacity == c.queue.Len() {
		toRemoveItem := c.queue.Back()
		c.queue.Remove(toRemoveItem)
		keyName := toRemoveItem.Value.(*cacheItem).key
		delete(c.items, keyName)
	}

	listItem := c.queue.PushFront(listItemValue)
	c.items[key] = listItem

	return false
}

// Get value from cache.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, exist := c.items[key]; exist {
		c.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}

	return nil, false
}

// Clear all value in cache.
func (c *lruCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.queue = NewList()

	newItems := make(map[Key]*ListItem, c.capacity)
	c.items = newItems
}
