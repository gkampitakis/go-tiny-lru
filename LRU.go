package LRU

import (
	"errors"
	"time"
)

type LRU struct {
	max   int
	size  int
	ttl   int64
	items map[string]*node
	first *node
	last  *node
}

type node struct {
	value  interface{}
	key    string
	expiry int64
	next   *node
	prev   *node
}

func New(max int, ttl int64) (error, *LRU) {
	if max < 0 {
		return errors.New("Invalid max value provided"), nil
	}

	if ttl < 0 {
		return errors.New("Invalid ttl value provided"), nil
	}

	return nil, &LRU{
		size:  0,
		max:   max,
		ttl:   ttl,
		first: nil,
		last:  nil,
		items: make(map[string]*node),
	}
}

func (cache *LRU) has(key string) (exists bool) {
	_, exists = cache.items[key]

	return
}

func (cache *LRU) Keys() []string {
	keys := make([]string, 0, len(cache.items))

	for k := range cache.items {

		keys = append(keys, k)

	}

	return keys

	// keys := make([]string, len(cache.items)) //TODO: test in benchmarks
	// i := 0

	// for k := range cache.items {

	// 	keys[i] = k
	// 	i++

	// }

	// return keys
}

func (cache *LRU) Clear() *LRU {
	cache.size = 0
	cache.first = nil
	cache.last = nil
	cache.items = make(map[string]*node)

	return cache
}

func (cache *LRU) evict() {
	if cache.first == nil {
		return
	}

	item := cache.first
	delete(cache.items, item.key)

	cache.first = item.next
	if cache.first != nil {
		cache.first.prev = nil
	}
	cache.size--
}

func (cache *LRU) Get(key string) interface{} {
	if cache.has(key) {
		item := cache.items[key]

		if cache.ttl > 0 && item.expiry < dateNow()+cache.ttl {
			cache.Delete(key)
		} else {

			cache._set(key, item.value, true)
			return item.value
		}
	}

	return nil
}

func (cache *LRU) Delete(key string) *LRU {
	if cache.has(key) {
		item := cache.items[key]

		delete(cache.items, key)
		cache.size--

		if item.prev != nil {
			item.prev.next = item.next
		}

		if item.next != nil {
			item.next.prev = item.prev
		}

		if cache.first == item {
			cache.first = item.next
		}

		if cache.last == item {
			cache.last = item.prev
		}
	}

	return cache
}

func (cache *LRU) _set(key string, value interface{}, bypass bool) {
	var item *node

	if cache.has(key) {
		item = cache.items[key]
		item.value = value

		if !bypass {
			var expiry int64 = cache.ttl

			if cache.ttl > 0 {
				expiry = cache.ttl + dateNow()
			}

			item.expiry = expiry
		}

		if cache.last != item {
			last := cache.last
			next := item.next
			prev := item.prev

			if cache.first == item {
				cache.first = item.next
			}

			item.next = nil
			item.prev = cache.last
			last.next = item

			if prev != nil {
				prev.next = next
			}

			if next != nil {
				next.prev = prev
			}
		}
	} else {
		if cache.max > 0 && cache.size >= cache.max {
			cache.evict()
		}

		var expiry int64 = cache.ttl

		if cache.ttl > 0 {
			expiry = cache.ttl + dateNow()
		}

		item = newNode(
			key,
			value,
			expiry,
			nil,
			cache.last,
		)
		cache.items[key] = item

		if cache.size++; cache.size == 1 {
			cache.first = item
		} else {
			cache.last.next = item
		}
	}

	cache.last = item
}

func (cache *LRU) Set(key string, value interface{}) *LRU {
	cache._set(key, value, false)
	return cache
}

func newNode(key string, value interface{}, expiry int64, next *node, prev *node) *node {

	return &node{value, key, expiry, next, prev}

}

func dateNow() int64 {

	return time.Now().Unix()

}
