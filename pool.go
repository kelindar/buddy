// Copyright (c) Manoj Babu Katragadda and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package buddy

type Pool struct {
	size uint32
	tree []uint32
	hash map[uint32]uint32
}

// expTwo returns true if the given number is a power of Two
func expTwo(num uint32) bool {
	if (num & (num - 1)) == 0 {
		return true
	}
	return false
}

//
func fitBlock(v uint32, bsize uint32) uint32 {
	if v > bsize {
		return 0
	}
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

func left(index uint32) uint32 {
	return 2*index + 1
}

func right(index uint32) uint32 {
	return 2*index + 2
}

func max(a uint32, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func New(zsize uint32) *Pool {
	if zsize < 1 || !expTwo(zsize) {
		return nil
	}
	nodesize := zsize

	p := new(Pool)
	p.size = zsize
	p.tree = make([]uint32, 2*zsize-1)
	p.tree[0] = zsize
	p.hash = make(map[uint32]uint32, zsize)
	for i := uint32(1); i < 2*zsize-1; i++ {
		if expTwo(i + 1) {
			nodesize /= 2
		}
		p.tree[i] = nodesize
	}
	return p
}

func (p *Pool) Allocate(hash32 uint32, zsize uint32) (uint32, bool) {
	if val, ok := p.hash[hash32]; ok {
		return val, true
	}

	if zsize < 1 {
		return 0, false
	}

	if !expTwo(zsize) {
		zsize = fitBlock(zsize, p.size)
	}

	index := uint32(0)
	nodesize := p.size
	if zsize > p.tree[index] {
		return 0, false
	}
	for ; nodesize != zsize; nodesize /= 2 {
		if p.tree[left(index)] >= zsize {
			index = left(index)
		} else {
			index = right(index)
		}
	}
	p.tree[index] = 0
	offset := (index+1)*nodesize - p.size

	for index > 0 {
		index = (index - 1) / 2
		p.tree[index] = max(p.tree[right(index)], p.tree[left(index)])

	}
	p.hash[hash32] = offset
	return offset, true
}

func (p *Pool) Release(hash32 uint32) {
	if _, ok := p.hash[hash32]; !ok {
		return
	}
	offset := p.hash[hash32]
	nodesize := uint32(1)
	index := offset + p.size - 1
	for ; p.tree[index] > 0; index = (index - 1) / 2 {
		nodesize *= 2
		if index == 0 {
			return
		}
	}
	p.tree[index] = nodesize

	for ; index > 0; nodesize *= 2 {
		index = (index - 1) / 2

		if p.tree[left(index)] == nodesize && p.tree[right(index)] == nodesize {
			p.tree[index] = nodesize * 2
		} else {
			p.tree[index] = max(p.tree[left(index)], p.tree[right(index)])
		}
	}
	delete(p.hash, hash32)
}

func (p *Pool) Find(hash32 uint32) (uint32, bool) {
	if offset, ok := p.hash[hash32]; ok {
		return offset, true
	}
	return 0, false
}
