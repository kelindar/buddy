package buddy

import (
	crand "crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"hash/crc32"
	"math/rand"
	"strings"
	"testing"
)

func BenchmarkPool(b *testing.B) {
	buffer := make([]byte, 1<<20)
	crand.Read(buffer)

	// Baseline implementaiton
	noop := &noop{}
	dict := &dict{
		data: make(map[uint32][]byte, 1024),
	}

	// Run benchmark suite
	for _, chance := range []int32{10, 50, 90} {
		run(b, chance, buffer, noop)
		run(b, chance, buffer, dict)
	}
}

func run(b *testing.B, chance int32, buffer []byte, impl pooler) {
	typ := strings.ReplaceAll(fmt.Sprintf("%T", impl), "*buddy.", "")
	name := fmt.Sprintf("%v-%v-%v", 100-chance, chance, typ)
	busy := make([]uint32, 1024)
	b.Run(name, func(b *testing.B) {
		rand.Seed(0)
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			v := rand.Int31n(100)
			switch {
			case v < chance:
				offset := rand.Int31n(int32(len(buffer)) - 5000)
				value := buffer[offset : offset+rand.Int31n(1000)+1]
				busy[rand.Int31n(int32(len(busy)))] = impl.Store(value)
			default:
				offset := busy[rand.Int31n(int32(len(busy)))]
				impl.Load(offset)
			}
		}
	})
}

type pooler interface {
	Store(value []byte) uint32
	Load(offset uint32) ([]byte, bool)
	Delete(offset uint32) bool
}

// --------------------------- Baseline ----------------------------

// Naive hashmap implementation for baseline benchmarks
type dict struct {
	data map[uint32][]byte
}

// Store stores a value in the pool and returns an offset to it. If the value
// already exists, it returns an offset to an existing value instead.
func (p *dict) Store(value []byte) uint32 {
	hash := crc32.ChecksumIEEE(value)
	if _, ok := p.data[hash]; !ok {
		p.data[hash] = value
	}
	return hash
}

// LoadAt loads a value at a specified offset and returns the data and
// whether it exists or not.
func (p *dict) Load(offset uint32) ([]byte, bool) {
	v, ok := p.data[offset]
	return v, ok
}

// Delete removes the entry and frees up the space used by it.
func (p *dict) Delete(offset uint32) bool {
	delete(p.data, offset)
	return true
}

// --------------------------- Noop ----------------------------

// Doesn't do anything, for tracking
type noop struct {
}

// Store stores a value in the pool and returns an offset to it. If the value
// already exists, it returns an offset to an existing value instead.
func (p *noop) Store(value []byte) uint32 {
	return 0
}

// LoadAt loads a value at a specified offset and returns the data and
// whether it exists or not.
func (p *noop) Load(offset uint32) ([]byte, bool) {
	return nil, true
}

// Delete removes the entry and frees up the space used by it.
func (p *noop) Delete(offset uint32) bool {
	return true
}

// --------------------------- Pool -------------------------------

// TestPoolAllocate tests the pool allocator
func TestPoolAllocate(t *testing.T) {
	pool := New(nil, 1024)
	pool.memory = pool.alloc.Allocate(1024)
	assert.Equal(t, pool.size, len(pool.memory))
}

// TestStore tests if the value is stored in pool memory.
func TestStore(t *testing.T) {
	// test invalid size
	value := make([]byte, 2041)
	pool := New(nil, 1024)
	pool.memory = pool.alloc.Allocate(1024)
	offset, ok := pool.Store(value)
	assert.Equal(t, uint32(0), offset)
	assert.False(t, ok)
	// try to accommodate 514B, results offset 0.
	value = make([]byte, 258)
	offset, ok = pool.Store(value)
	assert.Equal(t, uint32(0), offset)
	assert.True(t, ok)
	// try to accommodate 514B, results offset 0.
	value = make([]byte, 258)
	offset, ok = pool.Store(value)
	assert.Equal(t, uint32(0), offset)
	assert.True(t, ok)
	// try to accommodate 14B, results offset 0 and false OOM
	value = make([]byte, 45)
	offset, ok = pool.Store(value)
	assert.Equal(t, uint32(512), offset)
	assert.True(t, ok)
	// try to accommodate 14B, results offset 0 and false OOM
	value = make([]byte, 31)
	offset, ok = pool.Store(value)
	assert.Equal(t, uint32(576), offset)
	assert.True(t, ok)
	// try to accommodate 14B, results offset 0 and false OOM
	value = make([]byte, 231)
	offset, ok = pool.Store(value)
	assert.Equal(t, uint32(768), offset)
	assert.True(t, ok)
}

// TestLoad tests if the value is returned at an offset
func TestLoad(t *testing.T) {
	// test invalid size
	value := make([]byte, 2041)
	pool := New(nil, 1024)
	pool.memory = pool.alloc.Allocate(1024)
	offset, ok := pool.Store(value)
	assert.Equal(t, uint32(0), offset)
	assert.False(t, ok)
	// try to accommodate 514B, results offset 0.
	value = make([]byte, 258)
	offset, ok = pool.Store(value)
	assert.Equal(t, uint32(0), offset)
	assert.True(t, ok)
	// try to accommodate 14B, results offset 0 and false OOM
	value = make([]byte, 45)
	offset, ok = pool.Store(value)
	assert.Equal(t, uint32(512), offset)
	assert.True(t, ok)
	// try to load data at offset and verify checksum for input value and output value
	var loadedVal []byte
	loadedVal, ok = pool.Load(offset)
	assert.Equal(t, crc32.ChecksumIEEE(value), crc32.ChecksumIEEE(loadedVal))
	assert.True(t, ok)
}

// TestDelete tests if the value is returned at an offset
func TestDelete(t *testing.T) {
	// test invalid size
	value := make([]byte, 2041)
	pool := New(nil, 1024)
	pool.memory = pool.alloc.Allocate(1024)
	offset, ok := pool.Store(value)
	assert.Equal(t, uint32(0), offset)
	assert.False(t, ok)
	// try to accommodate 514B, results offset 0.
	value = make([]byte, 258)
	offset, ok = pool.Store(value)
	assert.Equal(t, uint32(0), offset)
	assert.True(t, ok)
	// try to accommodate 514B, results offset 0.
	value = make([]byte, 49)
	offset, ok = pool.Store(value)
	assert.Equal(t, uint32(512), offset)
	assert.True(t, ok)
	// try to accommodate 514B, results offset 0.
	value = make([]byte, 189)
	offset, ok = pool.Store(value)
	assert.Equal(t, uint32(768), offset)
	assert.True(t, ok)
	// try to delete using value checksum
	value = make([]byte, 189)
	checksum := crc32.ChecksumIEEE(value)
	assert.True(t, pool.Delete(checksum))
	fmt.Println(pool.fill)
	// try to delete using value checksum
	value = make([]byte, 258)
	checksum = crc32.ChecksumIEEE(value)
	assert.True(t, pool.Delete(checksum))
	fmt.Println(pool.fill)
	// try to delete using value checksum
	value = make([]byte, 49)
	checksum = crc32.ChecksumIEEE(value)
	assert.True(t, pool.Delete(checksum))
	fmt.Println(pool.fill)
	// TODO: this shld be [0,0,0,0,0] eventually
}

// TestCapacityFor tests if the return value is just the next power of 2 greater/equal to given int
func TestCapacityFor(t *testing.T) {
	assert.Equal(t, uint32(4), capacityFor(3))
	assert.NotEqual(t, uint32(8), capacityFor(3))
	assert.Equal(t, uint32(256), capacityFor(256))
	assert.Equal(t, uint32(256), capacityFor(253))
}

// TestMakeBuddies tests if the new bitmap level contains the bits set if the parent bit is set
func TestMakeBuddies(t *testing.T) {
	assert.Equal(t, uint32(15), makeBuddies(uint32(3), uint32(1)))
	assert.Equal(t, uint32(51), makeBuddies(uint32(5), uint32(3)))
	assert.Equal(t, uint32(12), makeBuddies(uint32(2), uint32(1<<(1-1))))
	assert.Equal(t, uint32(65280), makeBuddies(uint32(240), uint32(1<<(4-1))))
}

// TestExpTwo tests if the return value is exp of 2 for the given int
func TestExpTwo(t *testing.T) {
	assert.Equal(t, uint32(4), expTwo(uint32(16)))
	assert.Equal(t, uint32(5), expTwo(uint32(32)))
}

// TestFindBuddy checks for the buddy bitOffset of given block
func TestFindBuddy(t *testing.T) {
	assert.Equal(t, uint32(5), findBuddy(uint32(4)))
	assert.Equal(t, uint32(6), findBuddy(uint32(7)))
}
