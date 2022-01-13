package buddy

// The page size on which the pool should operate.
const pageSize = 1 << 20 // 1MB slabs

// Allocator represents a slab allocator implementation that allows the pool
// to request pages of memory from somewhere.
type Allocator interface {
	Allocate(size int) []byte
	Grow(buffer []byte, size int) []byte
}

// --------------------------- Go Allocator ----------------------------

// goAlloc is a default allocator that uses Go runtime to allocate.
type goAlloc struct{} // Allocate allocates a slab of memory.
func (goAlloc) Allocate(size int) []byte {
	return make([]byte, size)
}

// Grow increases the capacity of the memory.
func (goAlloc) Grow(buffer []byte, size int) []byte {
	return append(buffer, make([]byte, size)...)
}

// --------------------------- Page (implements Buddy) ----------------------------

// Page represents a data structure to manage dynamic memory slot allocator using buddy system
type Page struct {
	size      uint32   // size of the memory block.
	available []uint32 // list to manage available slots.
	allocated uint32   // total memory used up for statistics.(includes internal fragments)
}

// --------------------------- Buddy Pool ----------------------------

// Pool represents a buddy memory sub-allocator that allows to sub-partition
// memory slabs provided to it and retrieve them.
type Pool struct {
	alloc     Allocator
	cache     map[uint32]uint32 // maintains offsets to values stored in pages
	pages     []Page            // List of logical pages to manage memory
	allocated uint32            // total allocated memory incluing internal fragments
	memory    []byte            // The slabs of data
}

// New creates a new buddy pool
func New(alloc Allocator) *Pool {
	if alloc == nil {
		alloc = goAlloc{}
	}
	return &Pool{
		alloc: alloc,
		pages: make([]Page, 0, 10),
	}
}

// Store stores a value in the pool and returns an offset to it. If the value
// already exists, it returns an offset to an existing value instead.
func (p *Pool) Store(value []byte) uint32 {
	// TODO: Look up cache if the value is already present and return offset
	// If the value is not present, fit the value in one of the pages and
	// update the page slots
	// If pool memory is insufficient grow the memory slab and create new page.
	// return the offset
	panic("not implemented")
}

// LoadAt loads a value at a specified offset and returns the data and
// whether it exists or not.
func (p *Pool) Load(offset uint32) ([]byte, bool) {
	// TODO: Look up cache and serve the actual value.
	panic("not implemented")
}

// Delete removes the entry and frees up the space used by it.
func (p *Pool) Delete(offset uint32) bool {
	// TODO: clean up the memory at offset and
	// free the corresponding page slots.
	// Divide the offset with pagesize to determine the page index.
	panic("not implemented")
}
