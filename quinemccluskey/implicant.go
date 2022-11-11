package quinemccluskey

import "strconv"

// implicant represents a cover of one or more minterms.
type implicant struct {
	// literals is the input bits that are strictly '0' or '1'.
	literals uint64
	// xMask is the mask of input bits that may be either '0' or '1'.
	xMask uint64
	// tag decribes the outputs for which the implicant applies.
	tag uint64
	// checked is set to true when the implicant is successfully combined
	// with another implicant, and the resulting implicant's tag is equal
	// to the input implicant. Remaining unchecked implicants are prime.
	checked bool
}

// combine tests if the calling and passed implicant may be combined. if so,
// it merges them into a single implicant which covers the union of minterms
// represented by the input implicants and returns the result. If a combination
// is not possible, an empty implicant is returned.
func (im1 *implicant) combine(im2 *implicant) implicant {
	// xMasks must be equal for combination
	if im1.xMask != im2.xMask {
		return implicant{}
	}

	// implicants may be combined if they differ by exactly one bit
	literalsDelta := im1.literals ^ im2.literals
	if bitCount(literalsDelta) != 1 {
		return implicant{}
	}

	// build combination implicant
	im := implicant{
		literals: im1.literals &^ literalsDelta, // set differing bit to '0' in literal
		xMask:    im1.xMask | literalsDelta,     // add differing bit to xMask
		tag:      im1.tag & im2.tag,             // tag becomes the logical product of previous tags
		checked:  false,                         // checked is false for new implicants
	}

	// mark input implicants as checked
	im1.checked = (im1.tag == im.tag)
	im2.checked = (im2.tag == im.tag)

	return im
}

// equals tests whether two implicants are equal in terms of literals and xMask
func (im1 implicant) equals(im2 implicant) bool {
	return (im1.literals == im2.literals) && (im1.xMask == im2.xMask)
}

// covers tests whether the calling implicant covers the passed minterm
func (im implicant) covers(minterm uint64) bool {
	return ((im.xMask & minterm) | im.literals) == minterm
}

// outputList returns a list of all outputs the calling implicant applies to in
// integer form
func (im implicant) outputList() []int {
	outputs := []int{}

	for i := 0; i <= msbPos(im.tag); i++ {
		if (im.tag>>i)&1 == 1 {
			outputs = append(outputs, i)
		}
	}

	return outputs
}

// stringify returns a string representation of the calling implicant using
// '1's, '0's, and 'x's as determined from its `term` and `xMask` members
func (im implicant) stringify(bits int) string {
	s := ""
	for i := 0; i < bits; i++ {
		if ((im.xMask >> (bits - 1 - i)) & 1) == 1 {
			s += "x"
			continue
		}

		s += strconv.FormatUint((im.literals>>(bits-1-i))&1, 10)
	}

	return s
}
