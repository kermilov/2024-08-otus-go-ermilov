package hw04lrucache

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
}

type lruCacheValue struct {
	key   Key
	value interface{}
}

func (l *lruCache) Clear() {
	l.items = make(map[Key]*ListItem, l.capacity)
	l.queue = NewList()
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	val, isExist := l.items[key]
	if isExist {
		l.queue.MoveToFront(val)
		return val.Value.(lruCacheValue).value, true
	}
	return nil, false
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	existValue, isExist := l.items[key]
	if isExist && existValue.Value == value {
		l.queue.MoveToFront(existValue)
		return true
	} else if isExist {
		l.queue.Remove(existValue)
	}
	if !isExist && l.queue.Len() == l.capacity {
		lastInQueue := l.queue.Back()
		delete(l.items, lastInQueue.Value.(lruCacheValue).key)
		l.queue.Remove(lastInQueue)
	}
	l.items[key] = l.queue.PushFront(lruCacheValue{key, value})
	return isExist
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
