package quinemccluskey

import (
	"math"

	"golang.org/x/exp/slices"
)

// LogicFunction represents a group of outputs to determing a minumum cost
// cover for.
type LogicFunction struct {
	printoutsEnabled      bool
	largestTerm           uint64
	implicantDisplayWidth int
	mintermDisplayWidth   int
	minterms              [][]uint64
	dontCares             [][]uint64
	m_implicantTable      implicantTable
	m_coverTable          coverTable
}

// Init zeroes all members of LogicFunction.
func (solver *LogicFunction) Init(enablePrintouts bool) {
	solver.printoutsEnabled = enablePrintouts
	solver.largestTerm = 0
	solver.implicantDisplayWidth = 0
	solver.mintermDisplayWidth = 0
	solver.minterms = [][]uint64{}
	solver.dontCares = [][]uint64{}
	solver.m_implicantTable.init()
}

// AddOutput will add an output to the LogicFunction to be included in the
// minimum cost cover.
func (solver *LogicFunction) AddOutput(minterms []uint64, dontCares []uint64) int {
	mintermSet := []uint64{}
	dontCareSet := []uint64{}

	// copy minterms into a sorted set
	for _, minterm := range minterms {
		mintermSet = insert(mintermSet, minterm, func(i int) bool {
			return mintermSet[i] >= minterm
		})
	}

	// copy dontCares into a sorted set
	for _, dontCare := range dontCares {
		// exclude don't cares that appear in minterms
		if slices.Contains(mintermSet, dontCare) {
			dontCareSet = insert(dontCareSet, dontCare, func(i int) bool {
				return dontCareSet[i] >= dontCare
			})
		}
	}

	// attempt to add the output to the implicant table
	status := solver.m_implicantTable.addOutput(mintermSet, dontCareSet, solver.printoutsEnabled)
	if status == 0 {
		// save the added outputs
		solver.minterms = append(solver.minterms, mintermSet)
		solver.dontCares = append(solver.dontCares, dontCareSet)

		// update saved display widths for printouts
		for _, minterm := range mintermSet {
			solver.largestTerm = uint64(math.Max(float64(solver.largestTerm), float64(minterm)))
		}
		for _, dontCare := range dontCareSet {
			solver.largestTerm = uint64(math.Max(float64(solver.largestTerm), float64(dontCare)))
		}
		solver.implicantDisplayWidth = msbPos(uint64(solver.largestTerm))
		solver.mintermDisplayWidth = 1 + int(math.Ceil(math.Log10(float64(solver.largestTerm))))
	}

	return status
}

// testOutput will return the outptut for a given input for a cover defined by
// the list of passed implicants.
func testOutput(implicants []implicant, output int, input uint64) bool {
	for _, im := range implicants {
		if (im.tag>>output)&1 == 1 && im.covers(input) {
			return true
		}
	}

	return false
}

// verify cover will test a list of passed implicants against the defined
// outputs for a LogicFunction and return a boolean for whether the list is a
// correct cover of the function.
func (solver *LogicFunction) verifyCover(cover []implicant) bool {
	testPassed := true

	for output := 0; output < solver.m_implicantTable.nOutputs; output++ {
		for input := 0; uint64(input) < (bitLimit(uint64(solver.largestTerm)) + 1); input++ {
			// do not test don't cares
			if slices.Contains(solver.dontCares[output], uint64(input)) {
				continue
			}

			// determine expected output from minterms
			expectedOutput := false
			if slices.Contains(solver.minterms[output], uint64(input)) {
				expectedOutput = true
			}

			// compare output of cover against expected output
			if testOutput(cover, output, uint64(input)) != expectedOutput {
				testPassed = false
			}
		}
	}

	return testPassed
}

// stringifyLogicFunction returns a string representation of a logic function
// descibed by the passed list of implicants. Outputs and input bits are
// printed using the labels described in InputLabels and OutputLabels.
func (solver LogicFunction) stringifyLogicFunction(implicants []implicant, inLabels InputLabels, outLabels OutputLabels) string {
	outputEquations := ""

	for output := 0; output < solver.m_implicantTable.nOutputs; output++ {
		logicalEquation := outLabels.Str(output) + " = "

		for i, prime := range implicants {
			if (prime.tag>>output)&1 == 0 {
				continue
			}

			for j := 0; j < solver.implicantDisplayWidth; j++ {
				bit := solver.implicantDisplayWidth - j - 1
				if (prime.xMask>>bit)&1 == 0 {
					logicalEquation += inLabels.Str(bit)
					if (prime.literals >> (bit) & 1) == 0 {
						logicalEquation += "'"
					}

					if j != solver.implicantDisplayWidth-1 {
						logicalEquation += "."
					}
				}
			}

			if logicalEquation[len(logicalEquation)-1:] == "." {
				logicalEquation = logicalEquation[:len(logicalEquation)-1]
			}

			if i != len(implicants)-1 {
				logicalEquation += " + "
			}
		}

		if logicalEquation[len(logicalEquation)-3:] == " + " {
			logicalEquation = logicalEquation[:len(logicalEquation)-3]
		}

		outputEquations += logicalEquation + "\n"
	}

	return outputEquations
}

// GetMinimumCostCover will solve the LogicFunction for a minimum cost cover
// and return a string representation of the solution with outputs and input
// bits printed using the labels described in InputLabels and OutputLabels.
func (solver *LogicFunction) GetMinimumCostCover(inLabels InputLabels, outLabels OutputLabels) string {
	// reduce the implicant table to identify prime implicants
	primeImplicants := solver.m_implicantTable.reduce(solver.implicantDisplayWidth, solver.mintermDisplayWidth, solver.printoutsEnabled)

	// solve the cover table of prime implicants for a minimum cost cover
	solver.m_coverTable.build(solver.minterms, primeImplicants)
	minimumCostCover := solver.m_coverTable.getMinimumCostCover(solver.implicantDisplayWidth, solver.mintermDisplayWidth, solver.printoutsEnabled)

	// verify that the found minimum cost cover is a correct solution
	if !solver.verifyCover(minimumCostCover) {
		return "failed to yeild a correct solution"
	}

	// return a string representation of yhe minimum cost cover
	return solver.stringifyLogicFunction(minimumCostCover, inLabels, outLabels)
}
