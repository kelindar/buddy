package buddy

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

// The page size on which the pool should operate.
const pageSize = 1 << 20 // 1MB slabs

// Allocator represents a slab allocator implementation that allows the pool
// to request pages of memory from somewhere.
type Allocator interface {
	Allocate(size int) []byte
}

// --------------------------- Go Allocator ----------------------------

// goAlloc is a default allocator that uses Go runtime to allocate.
type goAlloc struct{}

// Allocate allocates a slab of memory.
func (goAlloc) Allocate(size int) []byte {
	return make([]byte, size)
}

// --------------------------- Buddy Pool ----------------------------

// Pool represents a buddy memory sub-allocator that allows to sub-partition
// memory slabs provided to it and retrieve them.
type Pool struct {
	alloc  Allocator
	size   int
	fill   []uint32
	cache  map[uint32]uint32
	memory []byte // The slabs of data
}

// New creates a new buddy pool
func New(alloc Allocator, size int) *Pool {
	if alloc == nil {
		alloc = goAlloc{}
	}
	return &Pool{
		alloc: alloc,
		size:  size,
		fill:  make([]uint32, 1, 8),
		cache: make(map[uint32]uint32, 64),
	}
}

// Store stores a value in the pool and returns an offset to it. If the value
// already exists, it returns an offset to an existing value instead.
func (p *Pool) Store(value []byte) (uint32, bool) {
	hash32 := crc32.ChecksumIEEE(value)
	if val, ok := p.cache[hash32]; ok {
		return val, true
	}

	// prepend value with its length(32 bit, max 1<<32 char)
	buffer := make([]byte, len(value)+4, len(value)+4)
	binary.BigEndian.PutUint32(buffer[0:4], uint32(len(value)))
	copy(buffer[4:], value[:])

	blockSize := capacityFor(uint32(len(buffer)))
	if int(blockSize) > p.size {
		//fmt.Errorf("Requested block size: %d is greater than pool size: %d", blockSize, p.size)
		return 0, false
	}
	i := uint32(0)
	levelSize := uint32(p.size)
	offset := uint32(0)
	bitOffset := uint32(0)
	for {
		// check for the level given by capacityFor
		if blockSize == levelSize {
			// check if the level is completely occupied
			if p.fill[i] == 1<<(i+1)-1 {
				//fmt.Errorf("Requested block size: %d cannot be accommodated in the memory pool", blockSize)
				return offset, false
			}
			// Once the level was found, look for empty blocks
			j := uint32(1 << i) // j equals to number of bits used to represent the ith level
			for {
				if p.fill[i]&(1<<(j-1)) == 0 {
					bitOffset = (1 << i) - j
					offset = bitOffset * levelSize
					break
				} else {
					j -= 1
				}
			}
			break
		} else {
			i += 1
			levelSize = levelSize >> 1
			// check if the last level is full and return OOM if true
			if int(i) == len(p.fill) && p.fill[i-1] == 1<<(i)-1 {
				return offset, false
			} else if int(i) == len(p.fill) {
				p.fill = append(p.fill, makeBuddies(p.fill[i-1], 1<<(i-1)))
			} else {
				continue
			}
		}
	}
	// store the data
	copy(p.memory[offset:], buffer)

	// set bits for the offset in the parent levels starting from the block level found.
	for {
		p.fill[i] |= 1 << ((1 << i) - 1 - bitOffset)
		i -= 1
		bitOffset /= 2
		if i == (1<<32 - 1) {
			break
		}
	}
	p.cache[hash32] = offset
	fmt.Println("................................")
	return offset, true
}

// LoadAt loads a value at a specified offset and returns the data and
// whether it exists or not.
func (p *Pool) Load(offset uint32) ([]byte, bool) {
	byteLength := binary.BigEndian.Uint32(p.memory[offset : offset+4])
	return p.memory[offset+4 : offset+4+byteLength], true
}

// Delete removes the entry and frees up the space used by it.
func (p *Pool) Delete(offset uint32) bool {
	panic("not implemented")
}

// ----------------------- Funcs ---------------------------------------

// capacityFor returns a power of 2 just greater than or equals given int
func capacityFor(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

// makeBuddies returns a new level bitmap int with already filled blocks in parent level
func makeBuddies(v uint32, level uint32) uint32 {
	// if level has 1010, return value will be 11001100(every 1,0 gets repeated 2 times as buddies)
	i := uint32(0)
	r := uint32(0)
	for {
		if v&(1<<i) == (1 << i) {
			r |= (1 << (2 * i))
			r |= (1 << ((2 * i) + 1))
		}
		i += 1
		if i == level+1 {
			break
		}
	}
	return r
}
