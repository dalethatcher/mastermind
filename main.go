package main

import (
	"fmt"
	"github.com/golang-collections/go-datastructures/bitarray"
	"log"
	"math"
)

type Score struct {
	rightValueAndPosition   int
	rightValueWrongPosition int
}

type CodeScore struct {
	guess []int
	score Score
}

func CalculateScore(code []int, guess []int) Score {
	if len(code) == 0 || len(code) != len(guess) {
		panic(fmt.Sprint("Passed invalid code or guess", code, "and", guess))
	}

	codeMatched := make([]bool, len(code))
	guessMatched := make([]bool, len(guess))
	result := Score{}

	for i, c := range code {
		if guess[i] == c {
			codeMatched[i] = true
			guessMatched[i] = true
			result.rightValueAndPosition++
		}
	}

	for i, c := range code {
		if guessMatched[i] {
			continue
		}

		for j, g := range guess {
			if i != j && c == g && !codeMatched[j] {
				codeMatched[j] = true
				result.rightValueWrongPosition++
			}
		}
	}

	return result
}

func IndexToCode(numberOfHoles int, numberOfColours int, index int) []int {
	if index > NumberOfCombinations(numberOfHoles, numberOfColours)-1 {
		panic(fmt.Sprint("index ", index, " is larger than the number of combinations for ", numberOfHoles,
			" holes and ", numberOfColours, " colours!"))
	}

	result := make([]int, numberOfHoles)
	for i := numberOfHoles - 1; i >= 0 && index > 0; i-- {
		colour := index % numberOfColours

		result[i] = colour
		index = index / numberOfColours
	}

	return result
}

func CodeToIndex(numberOfHoles int, numberOfColours int, code []int) int {
	var result int

	for i := 0; i < numberOfHoles; i++ {
		result = (result * numberOfColours) + code[i]
	}

	return result
}

func SetBits(ba *bitarray.BitArray, bits []int) error {
	for _, bit := range bits {
		if err := (*ba).SetBit(uint64(bit)); err != nil {
			return err
		}
	}

	return nil
}

func NumberOfCombinations(numberOfHoles int, numberOfColours int) int {
	return int(math.Pow(float64(numberOfColours), float64(numberOfHoles)))
}

func GuessIsPossible(facts []CodeScore, guess []int) bool {
	for _, fact := range facts {
		score := CalculateScore(guess, fact.guess)

		if score != fact.score {
			return false
		}
	}

	return true
}

func FindPossibleCodesIndicies(numberOfHoles int, numberOfColours int, facts []CodeScore) bitarray.BitArray {
	result := bitarray.NewBitArray(uint64(NumberOfCombinations(numberOfHoles, numberOfColours)))

	combinations := NumberOfCombinations(numberOfHoles, numberOfColours)
	for i := 0; i < combinations; i++ {
		guess := IndexToCode(numberOfHoles, numberOfColours, i)

		if GuessIsPossible(facts, guess) {
			if result.SetBit(uint64(i)) != nil {
				log.Panicln("Failed to set bit ", i)
			}
		}
	}

	return result
}

func PossibleScores(numberOfHoles int) []Score {
	result := []Score{}

	for correct := 0; correct <= numberOfHoles; correct++ {
		for wrongPosition := 0; correct+wrongPosition <= numberOfHoles; wrongPosition++ {
			if !(correct == numberOfHoles-1 && wrongPosition == 1) {
				result = append(result, Score{rightValueAndPosition: correct, rightValueWrongPosition: wrongPosition})
			}
		}
	}

	return result
}

func main() {
	numberOfHoles := 3
	numberOfColours := 4
	numberOfCombinations := NumberOfCombinations(numberOfHoles, numberOfColours)
	foundScores := make(map[Score]bool)

	for i := 0; i < numberOfCombinations; i++ {
		code := IndexToCode(numberOfHoles, numberOfColours, i)

		for j := 0; j < numberOfCombinations; j++ {
			guess := IndexToCode(numberOfHoles, numberOfCombinations, j)

			score := CalculateScore(code, guess)

			if _, ok := foundScores[score]; !ok {
				fmt.Println("Found new score:", score)
				foundScores[score] = true
			}
		}
	}
}
