package quinemccluskey

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

// visualizeHeading conditionally prints a line of text stdout
func visualizeHeading(text string, show bool) {
	if show {
		fmt.Println(text)
	}
}

// ----------------------------------------------------------------
//   IMPLICANT TABLE VISUALIZATION
// ----------------------------------------------------------------

// visualizeLogicFunctionList prints a comma seperated list wrapped with
// parentheses to stdout for use by visualizeLogicFunction.
func visualizeLogicFunctionList(list []uint64) {
	fmt.Printf("(")
	for i, item := range list {
		fmt.Printf("%d", item)
		if i != len(list)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Println(")")
}

// visualizeLogicFunction conditionally prints a formal representation of a
// logic function to stdout.
func visualizeLogicFunction(minterms []uint64, dontCares []uint64, show bool) {
	if show {
		fmt.Printf("FUNCTION_ADDED: S")
		visualizeLogicFunctionList(minterms)
		if len(dontCares) > 0 {
			fmt.Printf("                D")
			visualizeLogicFunctionList(dontCares)
		}
		fmt.Printf("\n")
	}
}

// visualizeImplicantTableHorizontalBar prints a horizontal bar to stdout with
// the correct dimentions for the calling implicant table.
func visualizeImplicantTableHorizontalBar(columns int, headingWidth1 int, headingWidth2 int) {
	for i := 0; i < columns; i++ {
		fmt.Printf("+%s%s-------------",
			strings.Repeat("-", int(math.Max(4, float64(headingWidth1)))),
			strings.Repeat("-", int(math.Max(4, float64(headingWidth2)))))
	}
	fmt.Println("+")
}

// visualizeImplicantTableHeader prints a table header to stdout with the
// correct dimentions for the calling implicant table.
func visualizeImplicantTableHeader(columns int, headingWidth1 int, headingWidth2 int) {
	visualizeImplicantTableHorizontalBar(columns, headingWidth1, headingWidth2)
	for i := 0; i < columns; i++ {
		fmt.Printf("| term  %stags  %schecked ",
			strings.Repeat(" ", int(math.Max(0, float64(headingWidth1-4)))),
			strings.Repeat(" ", int(math.Max(0, float64(headingWidth2-4)))))
	}
	fmt.Println("|")
	visualizeImplicantTableHorizontalBar(columns, headingWidth1, headingWidth2)
}

// visualizeImplicantTableNilEntry prints an empty column to stdout with the
// correct dimentions for the calling implicant table.
func visualizeImplicantTableNilEntry(headingWidth1 int, headingWidth2 int) {
	fmt.Printf("| %s  %s          ",
		strings.Repeat(" ", int(math.Max(4, float64(headingWidth1)))),
		strings.Repeat(" ", int(math.Max(4, float64(headingWidth2)))))
}

// ImplicantTable.visualize conditionally prints a tabular representation of the
// calling implicant table to stdout.
func (table implicantTable) visualize(implicantDisplayWidth int, show bool) {
	if show {
		visualizeImplicantTableHeader(len(table.columns), implicantDisplayWidth, table.nOutputs)
		for group := 0; group < len(table.columns[0]); group++ {
			groupEntries := 0
			for _, column := range table.columns {
				if len(column) > group {
					groupEntries = int(math.Max(float64(groupEntries), float64(len(column[group]))))
				}
			}
			for entry := 0; entry < groupEntries; entry++ {
				for _, column := range table.columns {
					if len(column) > group {
						if len(column[group]) > entry {
							f := "| %s%s  %0" + strconv.FormatInt(int64(table.nOutputs), 10) + "b%s  %t   "
							fmt.Printf(f,
								column[group][entry].stringify(implicantDisplayWidth),
								strings.Repeat(" ", int(math.Max(0, float64(4-implicantDisplayWidth)))),
								column[group][entry].tag,
								strings.Repeat(" ", int(math.Max(0, float64(4-table.nOutputs)))),
								column[group][entry].checked)
							if column[group][entry].checked {
								fmt.Printf(" ")
							}
						} else {
							visualizeImplicantTableNilEntry(implicantDisplayWidth, table.nOutputs)
						}
					} else {
						visualizeImplicantTableNilEntry(implicantDisplayWidth, table.nOutputs)
					}
				}
				fmt.Printf("|\n")
			}
			visualizeImplicantTableHorizontalBar(len(table.columns), implicantDisplayWidth, table.nOutputs)
		}
		fmt.Printf("\n")
	}
}

// visualizePrimeImplicantList conditionally prints a list of implicants to
// stdout.
func visualizePrimeImplicantList(primes []implicant, implicantDisplayWidth int, show bool) {
	if show {
		visualizeHeading("PRIME IMPLICANTS:", true)
		for _, prime := range primes {
			fmt.Println("  " + prime.stringify(implicantDisplayWidth))
		}
		fmt.Printf("\n")
	}
}

// ----------------------------------------------------------------
//   COVER TABLE VISUALIZATION
// ----------------------------------------------------------------

// visualizeCoverTableHorizontalBar prints a horizontal bar to stdout with the
// correct dimentions for the calling cover table.
func visualizeCoverTableHorizontalBar(columns int, implicantDisplayWidth int, mintermDisplayWidth int) {
	fmt.Printf("+-%s-+-", strings.Repeat("-", implicantDisplayWidth))
	fmt.Printf("%s+\n", strings.Repeat("-", mintermDisplayWidth*columns))
}

// visualizeCoverTableHeader prints a table header to stdout with the correct
// dimentions for the calling cover table.
func visualizeCoverTableHeader(minterms [][]uint64, columns int, implicantDisplayWidth int, mintermDisplayWidth int) {
	visualizeCoverTableHorizontalBar(columns, implicantDisplayWidth, mintermDisplayWidth)
	fmt.Printf("| %s | ", strings.Repeat(" ", implicantDisplayWidth))
	for _, output := range minterms {
		for _, minterm := range output {
			s := fmt.Sprintf("%d", minterm)
			fmt.Printf("%s%s", s, strings.Repeat(" ", mintermDisplayWidth-len(s)))
		}
	}
	fmt.Println("|")
	visualizeCoverTableHorizontalBar(columns, implicantDisplayWidth, mintermDisplayWidth)
}

// ImplicantTable.visualize conditionally prints a tabular representation of
// the calling cover table to stdout.
func (table coverTable) visualize(implicantDisplayWidth int, mintermDisplayWidth int, show bool) {
	if show {
		totalMinterms := 0
		for _, output := range table.remainingMinterms {
			totalMinterms += len(output)
		}

		visualizeCoverTableHeader(table.remainingMinterms, totalMinterms, implicantDisplayWidth, mintermDisplayWidth)
		// print prime implicant covers
		for i, prime := range table.primes {
			fmt.Printf("| %s | ", prime.stringify(implicantDisplayWidth))
			for j, output := range table.remainingMinterms {
				for _, minterm := range output {
					if slices.Contains(table.covers[i][j], minterm) {
						fmt.Printf("x %s", strings.Repeat(" ", mintermDisplayWidth-2))
					} else {
						fmt.Printf("%s", strings.Repeat(" ", mintermDisplayWidth))
					}
				}
			}
			fmt.Println("|")
		}
		visualizeCoverTableHorizontalBar(totalMinterms, implicantDisplayWidth, mintermDisplayWidth)
		fmt.Printf("\n")
	}
}
