package bitset

import (
	"math/bits"
	"sync/atomic"
)

// BitSet is a thread-safe bitmap structure
type BitSet struct {
	segments []atomic.Uint32
}

// New create a BitSet that supports up to n bits
func New(n int) *BitSet {
	segmentCount := (n + 31) / 32
	return &BitSet{
		segments: make([]atomic.Uint32, segmentCount),
	}
}

// Set the i-th bit to 1
func (bs *BitSet) Set(i int) {
	segIdx := i / 32
	bitIdx := i % 32
	bs.segments[segIdx].Or(1 << bitIdx)
}

// Unset the i-th bit to 0
func (bs *BitSet) Unset(i int) {
	segIdx := i / 32
	bitIdx := i % 32
	for {
		_old := bs.segments[segIdx].Load()
		_new := _old &^ (1 << bitIdx)
		if bs.segments[segIdx].CompareAndSwap(_old, _new) {
			break
		}
	}
}

// IsSet check if the i-th position is 1
func (bs *BitSet) IsSet(i int) bool {
	segIdx := i / 32
	bitIdx := i % 32
	return (bs.segments[segIdx].Load() & (1 << bitIdx)) != 0
}

// Clear all bits
func (bs *BitSet) Clear() {
	for i := range bs.segments {
		bs.segments[i].Store(0)
	}
}

// Len returns the total number of supported digits
func (bs *BitSet) Len() int {
	return len(bs.segments) * 32
}

// Count Returns the number of all bits set to 1
func (bs *BitSet) Count() int {
	total := 0
	for i := range bs.segments {
		val := bs.segments[i].Load()
		total += bits.OnesCount32(val)
	}
	return total
}

// ToSlice returns the index of all bits set to 1
func (bs *BitSet) ToSlice() []int {
	var result []int
	for i := range bs.segments {
		val := bs.segments[i].Load()
		if val == 0 {
			continue
		}
		for j := 0; j < 32; j++ {
			if val&(1<<j) != 0 {
				result = append(result, i*32+j)
			}
		}
	}
	return result
}
