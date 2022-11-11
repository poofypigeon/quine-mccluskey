package quinemccluskey

// implicantColumn is a list of implicants divided into groups.
type implicantColumn [][]implicant

// insert add an implicant to the last group in the calling implicantColumn,
// avoiding duplicates.
func (column *implicantColumn) insert(im implicant) {
	if im.tag == 0 {
		return
	}

	insert((*column)[len(*column)-1], im, func(i int) bool {
		return (*column)[len(*column)-1][i].literals <= im.literals || bitCount((*column)[len(*column)-1][i].xMask) <= bitCount(im.xMask)
	})
}

// iterate attempts to combine implicants with those in consecutive groups.
// Sucessful combinations are added to a new implicantColumn with a group for
// each pair of consecutive groups in the previous implicantColumn. The
// resulting new implicantColumn is returned.
func (column implicantColumn) iterate() implicantColumn {
	newColumn := implicantColumn{}

	// iterate over every pair of consecutive groups
	for group, _ := range column[:len(column)-1] {
		// add a group to newColumn for current group pair
		newColumn = append(newColumn, []implicant{})

		// try to combine implicants and add the result to the new group
		for _, im1 := range column[group] {
			for _, im2 := range column[group+1] {
				newColumn.insert(im1.combine(&im2))
			}
		}
	}

	// remove any empty groups that have been appended
	for len(newColumn) > 0 && len(newColumn[len(newColumn)-1]) == 0 {
		newColumn = newColumn[:len(newColumn)-1]
	}

	return newColumn
}

// primes returns a list of all unchecked implicants in an implicantColumn
func (column implicantColumn) primes() []implicant {
	primeList := []implicant{}

	for _, group := range column {
		for _, term := range group {
			if !term.checked {
				primeList = append(primeList, term)
			}
		}
	}

	return primeList
}
