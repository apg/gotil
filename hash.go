/** 
 * Implements a basic power of 2 HashTable with a version for
 * quick "clearing" for reuse of grown hashtable.
 *
 * TODO:
 *  1. determine if my Clear() is better than zeroing memory
 *  2. linear probing -> quadratic probing
 *  3. Can I make this lock free? Cliff Click has a Lock Free Java HashTable 
 *     that just relies on CAS (not sure whether or not go's CAS will work or not.
 */

package gotil

import (
	L "github.com/okcupid/logchan"
	"math"
	)

const HASH_INITIAL_SIZE int = 32
const HASH_LOAD_FACTOR float32 = 0.60

type Hashable interface {
	Hash() uint32
	Equal(interface{}) bool
}

type HashTable struct {
	version uint32
	elements []KeyValue
	capacity uint32
	count uint32
}

type KeyValue struct {
	hash uint32
	Key Hashable
	Value interface{}
	version uint32 // doubles as a tombstone (= 0)
}

func NewHashTable() *HashTable {
	ret := new(HashTable)
	ret.version = 1
	ret.elements = make([]KeyValue, HASH_INITIAL_SIZE) 
	ret.capacity = uint32(HASH_INITIAL_SIZE)
	return ret
}

func (ht *HashTable) Put(key Hashable, value interface{}) {
	// expand if necessary

	lf := float32(ht.count) / float32(ht.capacity)
//	L.Printf(LOG_LOAD, "Load factor = %v", lf)
	if lf > HASH_LOAD_FACTOR {
		ht.grow()
	}

	ht._put(key, value)
}

func (ht *HashTable) _put(key Hashable, value interface{}) {
	var i uint32
	hash := key.Hash()
	probe := (ht.capacity - 1) & hash
	step := uint32(1)

	// do quadratic probing...
	for i = 0; i < ht.capacity; i++ {
		curr := &(ht.elements[probe])
		if ht.elements[probe].version != ht.version {
			// just put it in here.
			curr.hash = hash
			curr.Key = key
			curr.Value = value
			curr.version = ht.version
			ht.count++
			return
		} else if hash == curr.hash && key.Equal(curr.Key) {
			// overwrite
			curr.Value = value
			return
		}

		probe = (probe + step) & (ht.capacity - 1)
		step++
	}

	L.Printf(L.LOG_ERROR, "Didn't find a spot for %v", key)
}

func (ht *HashTable) Get(key Hashable) (value interface{}, ok bool) {
	var i uint32
	hash := key.Hash()
	step := uint32(1)

	probe := (ht.capacity - 1) & hash

	for i = 0; i < ht.capacity; i++ {
		curr := &(ht.elements[probe])
		if ht.version == curr.version && curr.hash == hash {
			if key.Equal(curr.Key) {
				ok = true
				value = curr.Value
				break
			}
		} else if ht.version != curr.version {
			break // if we hit an empty on our chain, nothing was put here...
		}

		probe = (probe + step) & (ht.capacity - 1)
		step++
	}

	return
}

func (ht *HashTable) Delete(key Hashable) {
	var i uint32
	hash := key.Hash()
	probe := (ht.capacity - 1) & hash
	step := uint32(1)

	for i = 0; i < ht.capacity; i++ {
		curr := &(ht.elements[probe])
		if ht.version == curr.version && curr.hash == hash {
			if key.Equal(curr.Key) {
				curr.version = 0
				break
			}
		} else if ht.version != curr.version {
			break // if we hit an empty on our chain, nothing was put here...
		}

		probe = (probe + step) & (ht.capacity - 1)
		step++
	}

	return
}

func (ht *HashTable) Keys() []Hashable {
	keys := make([]Hashable, 0, ht.count)

	for _, kv := range ht.elements {
		if ht.version == kv.version {
			keys = append(keys, kv.Key)
		}
	}

	return keys
}

func (ht *HashTable) Values() []Hashable {
	values := make([]Hashable, 0, ht.count)

	for _, kv := range ht.elements {
		if ht.version == kv.version {
			values = append(values, kv.Key)
		}
	}

	return values
}

func (ht *HashTable) Len() int {
	c := int(ht.count)
	return c
}

func (ht *HashTable) Cap() int {
	c := int(ht.capacity)
	return c
}

func (ht *HashTable) Clear() {
	ht.version++
	ht.count = 0
}

func (ht *HashTable) Truncate(newSize uint32) {
	exp := math.Ceil(math.Log2(float64(newSize)))
	nextPow2 := uint32(math.Pow(2.0, exp))

	if ht.capacity > nextPow2 {
		ht.elements = make([]KeyValue, nextPow2)
		ht.capacity = nextPow2
		ht.count = 0
		ht.version = 1 // we can reset the version cause all these are empty now.
	} else {
		ht.Clear()
	}
}

func (ht *HashTable) grow() {
	// rehash and update the version.
	nextPow2 := ht.capacity << 1
	oldV := ht.version
	oldE := ht.elements

	ht.capacity = nextPow2
	ht.elements = make([]KeyValue, nextPow2)
	ht.count = 0 // will be updated on Put()

	// rehash everything
	for _, kv := range oldE {
		if kv.version == oldV {
			ht._put(kv.Key, kv.Value)
		}
	}
}
