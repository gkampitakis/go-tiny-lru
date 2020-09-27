package LRU

import (
	"bou.ke/monkey"
	"reflect"
	"strconv"
	"testing"
)

type mockStruct struct {
	valueStr string
	valueInt int
}

func TestNew(t *testing.T) {
	t.Run("Should return invalid max value", func(t *testing.T) {
		err, _ := New(-10, 10)

		if err.Error() != "Invalid max value provided" {
			t.Errorf("Wanted error invalid max value but got %v", err)
		}
	})

	t.Run("Should return invalid ttl value", func(t *testing.T) {
		err, _ := New(10, -10)

		if err.Error() != "Invalid ttl value provided" {
			t.Errorf("Wanted error invalid ttl value but got %v", err)
		}
	})

	t.Run("Should initialize LRU", func(t *testing.T) {
		_, cache := New(10, 10)

		init := cache.max != 10 ||
			cache.size != 0 ||
			cache.ttl != 10 ||
			cache.first != nil ||
			cache.last != nil

		if init {
			t.Errorf("LRU was not initialized correctly")
		}
	})
}

func TestKeys(t *testing.T) {
	t.Run("Should return empty keys array", func(t *testing.T) {
		_, cache := New(10, 10)

		if len(cache.Keys()) > 0 {
			t.Errorf("Found items inside LRU")
		}
	})

	t.Run("Should return keys array", func(t *testing.T) {
		_, cache := New(10, 10)

		cache.Set("mock0", newMockStruct("test0", 10))
		cache.Set("mock1", newMockStruct("test1", 10))
		wantedValues := []string{"mock0", "mock1"}
		got := cache.Keys()

		if !reflect.DeepEqual(got, wantedValues) {
			t.Errorf("Wanted values: %v but got %v", wantedValues, got)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("Should return nil if no item", func(t *testing.T) {
		_, cache := New(10, 10)

		value := cache.Get("mockValue")

		if value != nil {
			t.Errorf("Returned value for expired key")
		}
	})

	t.Run("Should return nil if item expired", func(t *testing.T) {
		_, cache := New(10, 100)

		var counter = 0

		monkey.Patch(dateNow, func() int64 {
			value := int64(1 * counter)
			counter += 1000000

			return value
		})

		cache.Set("mockValue", newMockStruct("test0", 10))

		value := cache.Get("mockValue")

		if value != nil {
			t.Errorf("Returned value for non existent key")
		}
	})

	t.Run("Should return item", func(t *testing.T) {
		_, cache := New(10, 0)

		item := newMockStruct("test0", 10)

		cache.Set("mockValue", item)

		value := cache.Get("mockValue")

		if value != item || value == nil {
			t.Errorf("Returned wrong value")
		}
	})
}

func TestEvict(t *testing.T) {
	t.Run("Should try to evict even if LRU is empty", func(t *testing.T) {
		_, cache := New(10, 10)

		cache.evict()
	})

	t.Run("Should remove 1st item", func(t *testing.T) {
		_, cache := New(10, 10)

		cache.Set("mock0", newMockStruct("test0", 10))
		cache.Set("mock1", newMockStruct("test1", 10))
		cache.Set("mock2", newMockStruct("test2", 10))

		cache.evict()
		_, exists := cache.items["mock0"]

		test := exists ||
			cache.size != 2

		if test {
			t.Errorf("LRU eviction was wrong")
		}
	})
}

func TestClear(t *testing.T) {
	t.Run("Should clear LRU", func(t *testing.T) {
		_, cache := New(10, 10)

		cache.Set("mock0", newMockStruct("test0", 10))
		cache.Set("mock1", newMockStruct("test1", 10))
		cache.Clear()

		_, exists := cache.items["mock0"]

		test := cache.size != 0 ||
			cache.first != nil ||
			cache.last != nil ||
			exists

		if test {
			t.Errorf("LRU not Cleared")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("Should try to delete item even if not existent", func(t *testing.T) {
		_, cache := New(10, 0)

		cache.Delete("mockValue")
	})

	t.Run("Should delete the only item", func(t *testing.T) {
		_, cache := New(10, 0)

		cache.Set("mockValue", newMockStruct("test0", 10))
		cache.Delete("mockValue")

		test := cache.Get("mockValue") != nil ||
			cache.first != cache.last ||
			cache.first != nil

		if test {
			t.Errorf("LRU Delete with only one item failed")
		}
	})

	t.Run("Should delete the correct item", func(t *testing.T) {
		_, cache := New(10, 0)

		cache.Set("mockValue0", newMockStruct("test0", 10))
		cache.Set("mockValue1", newMockStruct("test1", 10))
		cache.Set("mockValue2", newMockStruct("test2", 10))
		cache.Set("mockValue3", newMockStruct("test3", 10))

		cache.Delete("mockValue2")

		if cache.has("mockValue2") || cache.size != 3 {
			t.Errorf("Value was not deleted")
		}
	})
}

func TestSet(t *testing.T) {
	var counter = 0

	monkey.Patch(dateNow, func() int64 {
		value := int64(1 * counter)
		counter += 1000000

		return value
	})

	t.Run("Should create item if not existent", func(t *testing.T) {
		_, cache := New(10, 0)

		item0, item1, item2 := newMockStruct("test0", 10), newMockStruct("test1", 10), newMockStruct("test2", 10)

		cache._set("mockValue0", item0, false)
		cache._set("mockValue1", item1, false)
		cache._set("mockValue2", item2, false)
		cache._set("mockValue1", item1, false)

		test := cache.size != 3 ||
			cache.first.value != item0 ||
			cache.last.value != item1

		if test {
			t.Errorf("Value was not added correctly to LRU")
		}
	})

	t.Run("Should evict item if max capacity reached and add it at the end", func(t *testing.T) {
		_, cache := New(3, 0)

		item0, item1, item2, item3 := newMockStruct("test0", 10),
			newMockStruct("test1", 10),
			newMockStruct("test2", 10),
			newMockStruct("test3", 10)

		cache._set("mockValue0", item0, false)
		cache._set("mockValue1", item1, false)
		cache._set("mockValue2", item2, false)
		cache._set("mockValue3", item3, false)

		test := cache.size != 3 ||
			cache.first.value != item1 ||
			cache.last.value != item3

		if test {
			t.Fail()
		}
	})

	t.Run("Should refresh expiry", func(t *testing.T) {
		_, cache := New(10, 100)

		item0, item1, item2 := newMockStruct("test0", 10),
			newMockStruct("test1", 10),
			newMockStruct("test2", 10)

		cache._set("mockValue0", item0, false)
		cache._set("mockValue1", item1, false)

		expiry := cache.items["mockValue0"].expiry

		cache._set("mockValue0", item2, false)

		test := expiry == cache.items["mockValue0"].expiry ||
			cache.last.value != item2

		if test {
			t.Fail()
		}
	})

	t.Run("Should not refresh expiry just move item as last accessed", func(t *testing.T) {
		_, cache := New(10, 0)

		item0, item1, item2 := newMockStruct("test0", 10),
			newMockStruct("test1", 10),
			newMockStruct("test2", 10)

		cache._set("mockValue0", item0, true)
		cache._set("mockValue1", item1, true)

		expiry := cache.items["mockValue0"].expiry

		cache._set("mockValue0", item2, true)

		test := expiry != cache.items["mockValue0"].expiry ||
			cache.last.value != item2

		if test {
			t.Fail()
		}
	})
}

/**
Benchmarks
*/

func BenchmarkSet(b *testing.B) {
	_, cache := New(1000, 1000)
	for n := 0; n < b.N; n++ {
		cache.Set("mockValue"+strconv.Itoa(n), newMockStruct("test"+strconv.Itoa(n), 10))
	}
}

func BenchmarkGet(b *testing.B) {
	_, cache := New(1000, 1000)

	cache.Set("mockValue", newMockStruct("test0", 10))

	for n := 0; n < b.N; n++ {
		cache.Get("mockValue")
	}
}

func BenchmarkClear(b *testing.B) {
	_, cache := New(1000, 1000)

	for n := 0; n < b.N; n++ {
		cache.Clear()
	}
}

func BenchmarkDelete(b *testing.B) {
	_, cache := New(1000, 1000)

	for n := 0; n < b.N; n++ {
		cache.Set("mockValue", newMockStruct("test0", 10))
		cache.Delete("mockValue")
		cache.Delete("mockValue")
	}
}

func BenchmarkKeys(b *testing.B) {
	_, cache := New(1000, 1000)

	for i := 0; i < 100; i++ {
		cache.Set("mockValue"+strconv.Itoa(i), newMockStruct("test"+strconv.Itoa(1), 10))
	}

	for n := 0; n < b.N; n++ {
		cache.Keys()
	}
}

/**
Utils
*/
func newMockStruct(valueStr string, valueInt int) *mockStruct {
	return &mockStruct{valueStr, valueInt}
}
