// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package buddy

type Buddy struct {
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

func NewBuddySystem(zsize int) *Buddy {
	if zsize < 1 || !isPowerOfTwo(zsize) {
		return nil
	}
	nodesize := zsize

	b := new(Buddy)
	b.size = zsize
	b.tree = make([]int, 2*zsize-1)
	b.tree[0] = zsize
	b.hash = make(map[int]int, zsize)
	for i := 1; i < 2*zsize-1; i++ {
		if isPowerOfTwo(i + 1) {
			nodesize /= 2
		}
		b.tree[i] = nodesize
	}
	return b
}

func (b *Buddy) Alloc(hash32 int, zsize int) (int, bool) {
	println(hash32)
	if val, ok := b.hash[hash32]; ok {
		return val, true
	}

	if zsize < 1 {
		return -1, false
	}

	if !isPowerOfTwo(zsize) {
		zsize = fitPowerOfTwo(zsize, b.size)
	}

	index := 0
	nodesize := b.size
	if zsize > b.tree[index] {
		return -1, false
	}
	for ; nodesize != zsize; nodesize /= 2 {
		if b.tree[leftNode(index)] >= zsize {
			index = leftNode(index)
		} else {
			index = rightNode(index)
		}
	}
	b.tree[index] = 0
	offset := (index+1)*nodesize - b.size

	for index > 0 {
		index = (index - 1) / 2
		b.tree[index] = max(b.tree[rightNode(index)], b.tree[leftNode(index)])

	}
	b.hash[hash32] = offset
	return offset, true
}

func (b *Buddy) Free(hash32 int) {
	if _, ok := b.hash[hash32]; !ok {
		return
	}
	offset := b.hash[hash32]
	nodesize := 1
	index := offset + b.size - 1
	for ; b.tree[index] > 0; index = (index - 1) / 2 {
		nodesize *= 2
		if index == 0 {
			return
		}
	}
	b.tree[index] = nodesize

	for ; index > 0; nodesize *= 2 {
		index = (index - 1) / 2

		if b.tree[leftNode(index)] == nodesize && b.tree[rightNode(index)] == nodesize {
			b.tree[index] = nodesize * 2
		} else {
			b.tree[index] = max(b.tree[leftNode(index)], b.tree[rightNode(index)])
		}
	}
	delete(b.hash, hash32)
}

func (b *Buddy) Find(hash32 int) (int, bool) {
	if offset, ok := b.hash[hash32]; ok {
		return offset, true
	}
	return -1, false
}
