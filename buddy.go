package buddy

// Buddy interface to assign slots and free them.
type Buddy interface {
	Fill(width uint32) uint32
	Free(offset uint32) error
}

// type Binary, Weighted, Fibonacci, Tertiary
// should be of type Page which implements
// Grant and Free.

// Fill returns a slot index to allocate.
func (p *Page) Fill(width uint32) uint32 {
	panic("not implemented")
}

// Free updates the available slots list and allocated
func (p *Page) Free(offset uint32) error {
	panic("not implemented")
}
