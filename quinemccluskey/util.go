package quinemccluskey

import "sort"

// bitCount returns the number of set bits in a v.
func bitCount(v uint64) int {
	b := 0
	for v != 0 {
		v &= v - 1
		b++
	}
	return b
}

// msbPos returns the position of the most significant bit for v.
func msbPos(v uint64) int {
	b := 0
	for v != 0 {
		v >>= 1
		b++
	}
	return b
}

// bitLimit returns the max integer representable with the same number of bits
// required to represent v.
func bitLimit(v uint64) uint64 {
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	return v
}

// insert inserts v into a sorted set s.
func insert[T comparable](s []T, v T, f func(i int) bool) []T {
	if len(s) == 0 {
		s = append(s, v)
		return s
	}

	i := sort.Search(len(s), f)

	if i == len(s) || s[i] != v {
		s = append(s, *new(T))
		copy(s[i+1:], s[i:])
		s[i] = v
	}

	return s
}
