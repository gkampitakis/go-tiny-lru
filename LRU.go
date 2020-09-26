// package LRU
package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type LRU struct {
	max   int
	size  int
	ttl   int64
	Items map[string]interface{}
	first interface{}
	last  interface{}
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
		Items: make(map[string]interface{}),
	}
}

func (cache *LRU) has(key string) (exists bool) {
	_, exists = cache.Items[key]

	return
}

func (cache *LRU) Keys() []string {

	keys := make([]string, 0, len(cache.Items))

	for k := range cache.Items {

		keys = append(keys, k)

	}

	return keys

	// keys := make([]string, len(cache.Items)) //TODO: test in benchmarks
	// i := 0

	// for k := range cache.Items {

	// 	keys[i] = k
	// 	i++

	// }

	// return keys
}

func (*LRU) Delete(key string) {
	return
}

func (cache *LRU) Set(key string, value interface{}) {

	cache.Items[key] = value

	return
}

func newNode(key string, value interface{}, expiry int64) *node {

	return &node{value, key, expiry, nil, nil}

}

func dateNow() int64 {

	return time.Now().Unix()

}

func main() {

	err, cache := New(10, 10)

	if err != nil {

		log.Fatal(err)

	}

	fmt.Println(dateNow())

	cache.Set("george", "name")
	cache.Set("george1", "name")

	fmt.Println(cache.has("george"))
	fmt.Println(cache.Keys())

	// cache.items.PushBack(newNode('['))

}

// func main() {
// 	mymap := make(map[int]string)
// 	keys := make([]int, 0, len(mymap))
// 	for k := range mymap {
// 			keys = append(keys, k)
// 	}
// }
