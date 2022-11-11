package quinemccluskey

import (
	"strconv"
)

// implicant table is list of implicantColumns representing iterations of
// combinations.
type implicantTable struct {
	nOutputs int
	columns  []implicantColumn
}

// init sets nOutputs, and initializes a column with enough space in each group
// for the max number of possible set bits to be added.
func (table *implicantTable) init() {
	table.nOutputs = 0
	table.columns = []implicantColumn{make([][]implicant, 0, 64)}
}

// addOutput takes a list of minterms and don't cares, and sorts them into
// groups for the calling implicantTable corresponding to the number of set
// bits in the term. For each output added, the corresponding tag bit is set
// to '1', and nOutputs is incremented.
func (table *implicantTable) addOutput(minterms []uint64, dontCares []uint64, printoutsEnabled bool) int {
	// maximum number of outputs exceeded
	if table.nOutputs == 64 {
		return -1
	}

	// copy all terms into a single new list
	terms := []uint64{}
	terms = append(terms, minterms...)
	terms = append(terms, dontCares...)

	// add minterms to groups based on the number of set bits
NEXT_TERM:
	for _, term := range terms {
		setBits := bitCount(term)

		// grow the slice to accomodate setBits as index
		if setBits >= len(table.columns[0]) {
			table.columns[0] = table.columns[0][:setBits+1]
		}

		// if the implicant already exists, set its tag bit for this output
		for i, im := range table.columns[0][setBits] {
			if term == im.literals {
				table.columns[0][setBits][i].tag |= 1 << table.nOutputs
				continue NEXT_TERM
			}
		}

		// otherwise build and add an implicant with its tag bit set for this output
		im := implicant{term, 0, (1 << table.nOutputs), false}
		table.columns[0][setBits] = append(table.columns[0][setBits], im)
	}

	table.nOutputs++

	visualizeLogicFunction(minterms, dontCares, printoutsEnabled)

	return 0
}

// reduce will solve the table by iterating columns until no new combinations
// can be made. After this process, the terms in each column that are still
// unchecked are prime implicants, which are placed in a list and returned.
func (table *implicantTable) reduce(implicantDisplayWidth int, mintermDisplayWidth int, printoutsEnabled bool) []implicant {
	visualizeBanner("TABLE REDUCTION FOR PRIME IMPLICANTS", printoutsEnabled)

	// iterate lists until no more combinations can be made
	visualizeHeading("TABLE: 0", printoutsEnabled)
	table.visualize(implicantDisplayWidth, printoutsEnabled)
	nextColumn := table.columns[len(table.columns)-1].iterate()
	for iter := 1; len(nextColumn) > 0; iter++ {
		table.columns = append(table.columns, nextColumn)
		visualizeHeading("TABLE: "+strconv.FormatInt(int64(iter), 10), printoutsEnabled)
		table.visualize(implicantDisplayWidth, printoutsEnabled)
		nextColumn = table.columns[len(table.columns)-1].iterate()
	}

	// add unchecked implicants from all lists to a new list
	// these are our prime implicants
	primes := []implicant{}
	for _, list := range table.columns {
		primes = append(primes, list.primes()...)
	}

	visualizePrimeImplicantList(primes, implicantDisplayWidth, printoutsEnabled)

	return primes
}
