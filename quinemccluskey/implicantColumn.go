package quinemccluskey

import (
	"fmt"
	"sync"
)

// implicantColumn is a list of implicants divided into groups.
type implicantColumn []map[implicant]bool

func processGroup(group int, column *implicantColumn, newColumn *implicantColumn) {
	list0 := []implicant{}
	for im := range (*column)[group] {
		list0 = append(list0, im)
	}

	list1 := []implicant{}
	for im := range (*column)[group+1] {
		list1 = append(list1, im)
	}

	if (*newColumn)[group] == nil {
		(*newColumn)[group] = map[implicant]bool{}
	}

	// try to combine implicants and add the result to the new group
	for i := range list0 {
		for j := range list1 {
			if list0[i].tag&list1[j].tag != 0 {
				delete((*column)[group], list0[i])
				delete((*column)[group+1], list1[j])

				newImplicant := list0[i].combine(&list1[j])
				if newImplicant.tag != 0 {
					(*newColumn)[group][newImplicant] = true
				}

				(*column)[group][list0[i]] = true
				(*column)[group+1][list1[j]] = true
			}
		}
	}
	fmt.Printf("%d/%d\n", group, len(*column)-1)
}

// iterate attempts to combine implicants with those in consecutive groups.
// Sucessful combinations are added to a new implicantColumn with a group for
// each pair of consecutive groups in the previous implicantColumn. The
// resulting new implicantColumn is returned.
func (column *implicantColumn) iterate() implicantColumn {
	var newColumn implicantColumn = make([]map[implicant]bool, len(*column)-1)

	var wg sync.WaitGroup

	// iterate over every even pair of consecutive groups
	for group := range (*column)[:len(*column)-1] {
		if group%2 == 1 {
			continue
		}
		wg.Add(1)
		group := group
		go func() {
			defer wg.Done()
			processGroup(group, column, &newColumn)
		}()
	}

	wg.Wait()

	// iterate over every odd pair of consecutive groups
	for group := range (*column)[:len(*column)-1] {
		if group%2 == 0 {
			continue
		}
		wg.Add(1)
		group := group
		go func() {
			defer wg.Done()
			processGroup(group, column, &newColumn)
		}()
	}

	wg.Wait()

	// trim excess empty groups
	for len(newColumn) > 0 && len(newColumn[len(newColumn)-1]) == 0 {
		newColumn = newColumn[:len(newColumn)-1]
	}

	return newColumn
}

// primes returns a list of all unchecked implicants in an implicantColumn
func (column implicantColumn) primes() []implicant {
	primeList := []implicant{}

	for _, group := range column {
		for term := range group {
			if !term.checked {
				primeList = append(primeList, term)
			}
		}
	}

	return primeList
}
