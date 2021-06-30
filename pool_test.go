// Copyright (c) Manoj Babu Katragadda and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package buddy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpTwo(t *testing.T) {
	assert.Equal(t, true, expTwo(4))
	assert.Equal(t, false, expTwo(5))
}

func TestFitBlock(t *testing.T) {
	assert.Equal(t, uint32(4), fitBlock(3, 32))
	assert.Equal(t, uint32(32), fitBlock(17, 32))
}

func TestNew(t *testing.T) {
	pool := New(128)
	assert.Equal(t, uint32(128), pool.size)
}

func TestBuddy(t *testing.T) {
	pool := New(128)
	offset, _ := pool.Allocate(31221, 12)
	assert.Equal(t, uint32(0), offset)
	offset, _ = pool.Allocate(31222, 12)
	assert.Equal(t, uint32(16), offset)
	offset, ok := pool.Allocate(31223, 65)
	assert.Equal(t, false, ok)
	pool.Release(31221)
	pool.Release(31222)
	offset, ok = pool.Allocate(31223, 65)
	assert.Equal(t, uint32(0), offset)
}
