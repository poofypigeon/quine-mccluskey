package quinemccluskey

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type coverTable struct {
	primes            []implicant
	covers            [][][]uint64
	remainingMinterms [][]uint64
}

// build initialized cover table by copying the passed minterms and primes list
// into table.remainingMinterms and table.primes respectively, and populating
// covers with lists of of all the provided minterms covered by the implicant
// at the same index in primes.
func (table *coverTable) build(minterms [][]uint64, primes []implicant) {
	// deep copy minterms into table.remainingMinterms
	for i, output := range minterms {
		table.remainingMinterms = append(table.remainingMinterms, []uint64{})
		for _, minterm := range output {
			table.remainingMinterms[i] = append(table.remainingMinterms[i], minterm)
		}
	}

	// copy primes into table.primes
	table.primes = make([]implicant, len(primes))
	copy(table.primes, primes)

	// build a list of all of the minterms that each prime covers
	for p, prime := range table.primes {
		table.covers = append(table.covers, [][]uint64{})
		for o, output := range table.remainingMinterms {
			table.covers[p] = append(table.covers[p], []uint64{})
			for _, minterm := range output {
				if (prime.tag>>o)&1 == 1 && prime.covers(minterm) {
					table.covers[p][o] = append(table.covers[p][o], minterm)
				}
			}
		}
	}
}

// removePrimeAndCovers removes the passed implicant as well as all minterms
// that the implicant covers from all of the calling coverTable's stored lists.
func (table *coverTable) removePrimeAndCovers(prime implicant) {
	primeIndex := slices.Index(table.primes, prime)
	if primeIndex == -1 {
		return
	}

	for _, output := range prime.outputList() {
		nonVolatileCovers := make([]uint64, len(table.covers[primeIndex][output]))
		copy(nonVolatileCovers, table.covers[primeIndex][output])

		for _, coveredMinterm := range nonVolatileCovers {
			// remove minterm from all covers
			for _, cover := range table.covers {
				i := slices.Index(cover[output], coveredMinterm)
				if i >= 0 {
					cover[output] = append(cover[output][:i], cover[output][i+1:]...)
				}
			}

			// remove minterm from remainingMinterms for the current output
			for output, _ := range table.remainingMinterms {
				if (table.primes[primeIndex].tag>>output)&1 == 0 {
					continue
				}

				i := slices.Index(table.remainingMinterms[output], coveredMinterm)
				if i >= 0 {
					table.remainingMinterms[output] = append(table.remainingMinterms[output][:i], table.remainingMinterms[output][i+1:]...)
				}
			}
		}
	}

	// remove essential prime from primes and remove the coresponding entry in covers
	table.primes = append(table.primes[:primeIndex], table.primes[primeIndex+1:]...)
	table.covers = append(table.covers[:primeIndex], table.covers[primeIndex+1:]...)
}

func (table *coverTable) removeEssentialPrimes() []implicant {
	essentialPrimes := []implicant{}

	for out, output := range table.remainingMinterms {
	NEXT_MINTERM:
		for _, minterm := range output {
			// continue if minterm occurs in more than one cover
			essentialPrimeIndex := -1
			for i, cover := range table.covers {
				if len(cover[out]) > 0 {
					if slices.Contains(cover[out], minterm) {
						if essentialPrimeIndex >= 0 {
							continue NEXT_MINTERM
						}
						essentialPrimeIndex = i
					}
				}
			}

			// zero occurances found
			if essentialPrimeIndex == -1 {
				continue NEXT_MINTERM
			}

			// add prime to essential primes
			if !slices.Contains(essentialPrimes, table.primes[essentialPrimeIndex]) {
				essentialPrimes = append(essentialPrimes, table.primes[essentialPrimeIndex])
			}
		}
	}

	// remove essential primes from table
	for _, essentialPrime := range essentialPrimes {
		table.removePrimeAndCovers(essentialPrime)
	}

	return essentialPrimes
}

// getMinimumCostCover will solve for a minimum cost cover for all of the
// implicants contained in the calling coverTable, returning a list of
// implicants that make up the minimum cost cover.
func (table *coverTable) getMinimumCostCover(implicantDisplayWidth int, mintermDisplayWidth int, printoutsEnabled bool) []implicant {
	table.visualize(implicantDisplayWidth, mintermDisplayWidth, printoutsEnabled)

	// capture and remove essential prime implicants from the coverTable
	minimumCover := table.removeEssentialPrimes()
	visualizeHeading("ESSENTIAL PRIMES REMOVED", printoutsEnabled)

	// get the total number of remaining minterms across all outputs
	totalRemainingMinterms := 0
	for _, output := range table.remainingMinterms {
		totalRemainingMinterms += len(output)
	}

	table.visualize(implicantDisplayWidth, mintermDisplayWidth, printoutsEnabled)

	// capture and remove primes iteratively until all minterms are covered
	for totalRemainingMinterms > 0 {
		// pL is the set of remaining primes which cover the greatest number
		// of minterms
		pLIndices := []int{}
		mostCoveredMinterms := 0
		for i, cover := range table.covers {
			totalMintermsCovered := 0
			for _, output := range cover {
				totalMintermsCovered += len(output)
			}
			if totalMintermsCovered > mostCoveredMinterms {
				mostCoveredMinterms = totalMintermsCovered
				pLIndices = []int{i}
			} else if totalMintermsCovered == mostCoveredMinterms {
				pLIndices = append(pLIndices, i)
			}
		}

		// pI is the first occurance from the set of primes with the fewest
		// literals (most xMask bits) in pL
		pI := table.primes[pLIndices[0]]
		for _, index := range pLIndices[1:] {
			fmt.Println("checking", table.primes[index])
			if bitCount(table.primes[index].xMask) > bitCount(pI.xMask) {
				fmt.Println("picked")
				pI = table.primes[index]
			}
		}

		// capture and remove all occurances of the selected prime implicant
		// and the minterms that it covers from the coverTable
		table.removePrimeAndCovers(pI)
		minimumCover = append(minimumCover, pI)

		// update the count of remaining minterms across all outputs
		totalRemainingMinterms = 0
		for _, output := range table.remainingMinterms {
			totalRemainingMinterms += len(output)
		}

		visualizeHeading(pI.stringify(implicantDisplayWidth)+" REMOVED", printoutsEnabled)
		table.visualize(implicantDisplayWidth, mintermDisplayWidth, printoutsEnabled)
	}

	return minimumCover
}
