package buddy

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
	memory [][]byte // The slabs of data
}

// New creates a new buddy pool
func New(alloc Allocator) *Pool {
	if alloc == nil {
		alloc = goAlloc{}
	}
	return &Pool{
		alloc: alloc,
	}
}

// Store stores a value in the pool and returns an offset to it. If the value
// already exists, it returns an offset to an existing value instead.
func (p *Pool) Store(value []byte) uint32 {
	panic("not implemented")
}

// LoadAt loads a value at a specified offset and returns the data and
// whether it exists or not.
func (p *Pool) Load(offset uint32) ([]byte, bool) {
	panic("not implemented")
}
