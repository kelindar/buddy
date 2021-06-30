// Copyright (c) Manoj Babu Katragadda and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package buddy

type Pool struct {
	size int
	tree []int
	hash map[int]int
}

func isPowerOfTwo(num int) bool {
	if (num & (num - 1)) == 0 {
		return true
	}
	return false
}

func fitPowerOfTwo(size int, bsize int) int {
	if size > bsize {
		return 0
	}
	for {
		if bsize/2 < size {
			return bsize
		} else {
			bsize /= 2
		}
	}
}

func leftNode(index int) int {
	return 2*index + 1
}

func rightNode(index int) int {
	return 2*index + 2
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func New(zsize int) *Pool {
	if zsize < 1 || !isPowerOfTwo(zsize) {
		return nil
	}
	nodesize := zsize

	p := new(Pool)
	p.size = zsize
	p.tree = make([]int, 2*zsize-1)
	p.tree[0] = zsize
	p.hash = make(map[int]int, zsize)
	for i := 1; i < 2*zsize-1; i++ {
		if isPowerOfTwo(i + 1) {
			nodesize /= 2
		}
		p.tree[i] = nodesize
	}
	return p
}

func (p *Pool) Alloc(hash32 int, zsize int) (int, bool) {
	println(hash32)
	if val, ok := p.hash[hash32]; ok {
		return val, true
	}

	if zsize < 1 {
		return -1, false
	}

	if !isPowerOfTwo(zsize) {
		zsize = fitPowerOfTwo(zsize, p.size)
	}

	index := 0
	nodesize := p.size
	if zsize > p.tree[index] {
		return -1, false
	}
	for ; nodesize != zsize; nodesize /= 2 {
		if p.tree[leftNode(index)] >= zsize {
			index = leftNode(index)
		} else {
			index = rightNode(index)
		}
	}
	p.tree[index] = 0
	offset := (index+1)*nodesize - p.size

	for index > 0 {
		index = (index - 1) / 2
		p.tree[index] = max(p.tree[rightNode(index)], p.tree[leftNode(index)])

	}
	p.hash[hash32] = offset
	return offset, true
}

func (p *Pool) Free(hash32 int) {
	if _, ok := p.hash[hash32]; !ok {
		return
	}
	offset := p.hash[hash32]
	nodesize := 1
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

		if p.tree[leftNode(index)] == nodesize && p.tree[rightNode(index)] == nodesize {
			p.tree[index] = nodesize * 2
		} else {
			p.tree[index] = max(p.tree[leftNode(index)], p.tree[rightNode(index)])
		}
	}
	delete(p.hash, hash32)
}

func (p *Pool) Find(hash32 int) (int, bool) {
	if offset, ok := p.hash[hash32]; ok {
		return offset, true
	}
	return -1, false
}
