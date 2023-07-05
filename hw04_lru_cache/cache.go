package main

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
	mu       sync.RWMutex
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

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	newCacheItem := cacheItem{
		key:   key,
		value: value,
	}

	listItem, exist := lru.items[key]
	if exist {
		listItem.Value = newCacheItem
		lru.queue.MoveToFront(listItem)

		return true
	}

	listItem = lru.queue.PushFront(newCacheItem)
	lru.items[key] = listItem

	if lru.capacity < lru.queue.Len() {
		listItem := lru.queue.Back()

		lru.queue.Remove(listItem)
		delete(lru.items, listItem.Value.(cacheItem).key)
	}
	return false
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.mu.RLock()
	defer lru.mu.RUnlock()

	listItem, exist := lru.items[key]
	if exist {
		lru.queue.MoveToFront(listItem)

		return listItem.Value.(cacheItem).value, true
	}
	return nil, false
}

func (lru *lruCache) Clear() {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	lru.queue = NewList()
	lru.items = make(map[Key]*ListItem, lru.capacity)
}
