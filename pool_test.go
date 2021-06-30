package buddy

import (
	crand "crypto/rand"
	"fmt"
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
