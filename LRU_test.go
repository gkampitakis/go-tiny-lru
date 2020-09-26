package LRU

import (
	"bou.ke/monkey"
	"reflect"
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

	})

	t.Run("Should delete the only item", func(t *testing.T) {

	})

	t.Run("Should delete the last item", func(t *testing.T) {

	})
}

func TestSet(t *testing.T) {
	t.Run("Should create item if not existent", func(t *testing.T) {

	})

	t.Run("Should evict item if max capacity reached and add it at the end", func(t *testing.T) {

	})

	t.Run("Should evict item if max capacity reached", func(t *testing.T) {

	})

	t.Run("Should refresh expiry", func(t *testing.T) {

	})

	t.Run("Should not refresh expiry just move item as last accessed", func(t *testing.T) {

	})
}

func newMockStruct(valueStr string, valueInt int) *mockStruct {
	return &mockStruct{valueStr, valueInt}
}
