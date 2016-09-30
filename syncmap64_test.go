package syncmap

import (
	"testing"
)

func Test_New64(t *testing.T) {
	m1 := New64()
	if m1 == nil {
		t.Error("New64(): map is nil")
	}
	if m1.shardCount != defaultShardCount {
		t.Error("New(): map's shard count is wrong")
	}
	if m1.Size() != 0 {
		t.Error("New(): new map should be empty")
	}

	var shardCount uint8 = 64
	m2 := NewWithShard64(shardCount)
	if m2 == nil {
		t.Error("NewWithShard64(): map is nil")
	}
	if m2.shardCount != shardCount {
		t.Error("NewWithShard64(): map's shard count is wrong")
	}
	if m2.Size() != 0 {
		t.Error("New64(): new map should be empty")
	}
}

func Test_Set64(t *testing.T) {
	m := New64()
	m.Set(1, 1)
	m.Set(2, 2)
	if m.Size() != 2 {
		t.Error("map should have 2 items.")
	}
}

func Test_Get64(t *testing.T) {
	m := New64()
	v1, ok := m.Get(7788414)
	if ok {
		t.Error("ok should be false when key is missing")
	}
	if v1 != nil {
		t.Error("value should be nil for missing key")
	}

	m.Set(1, 1)

	v2, ok := m.Get(1)
	if !ok {
		t.Error("ok should be true when key exists")
	}
	if 1 != v2.(int) {
		t.Error("value should be an integer of value 1")
	}
}

func Test_Has64(t *testing.T) {
	m := New64()
	if m.Has(1) {
		t.Error("Has should return False for missing key")
	}

	m.Set(1, 1)
	if !m.Has(1) {
		t.Error("Has should return True for existing key")
	}
}

func Test_Delete64(t *testing.T) {
	m := New64()
	m.Set(1, 1)
	m.Delete(1)
	if m.Has(1) {
		t.Error("Delete shoudl remove the given key from map")
	}
}

func Test_Size64(t *testing.T) {
	m := New64()
	for i := 0; i < 42; i++ {
		m.Set(uint64(i), i)
	}
	if m.Size() != 42 {
		t.Error("Size doesn't return the right number of items")
	}
}

func Test_Flush64(t *testing.T) {
	var shardCount uint8 = 64
	m := NewWithShard64(shardCount)
	for i := 0; i < 42; i++ {
		m.Set(uint64(i), i)
	}
	count := m.Flush()
	if count != 42 {
		t.Error("Flush should return the size before removing")
	}
	if m.Size() != 0 {
		t.Error("Flush should remove all items from map", m.Size())
	}
	if m.shardCount != shardCount {
		t.Error("map should have the same shardCount after Flush")
	}
}

/*
func Test_IterKeys(t *testing.T) {
	loop := 100
	expectedKeys := make([]uint64, loop)

	m := New64()
	for i := 0; i < loop; i++ {
		key := uint64(i)
		expectedKeys[i] = key
		m.Set(key, i)
	}

	keys := make([]uint64, 0)
	for key := range m.IterKeys() {
		keys = append(keys, key)
	}

	if len(keys) != len(expectedKeys) {
		t.Error("IterKeys doesn't loop the right times")
	}

}
*/

func Test_Pop64(t *testing.T) {
	m := New64()
	// m.Pop()

	m.Set(uint64(1), 1)

	k, v := m.Pop()
	if k != uint64(1) && v.(int) != 1 {
		t.Error("Pop should returns the only item")
	}
	if m.Size() != 0 {
		t.Error("Size should be 0 after pop the only item")
	}
}
