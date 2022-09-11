package main

import (
	"fmt"
	"github.com/golang-collections/go-datastructures/bitarray"
	"log"
	"math"
)

type Rules struct {
	holes   int
	colours int
}

type Score struct {
	correct   int
	misplaced int
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
			result.correct++
		}
	}

	for i, c := range code {
		if guessMatched[i] {
			continue
		}

		for j, g := range guess {
			if i != j && c == g && !codeMatched[j] {
				codeMatched[j] = true
				result.misplaced++
			}
		}
	}

	return result
}

func IndexToCode(rules Rules, index int) []int {
	if index > NumberOfCombinations(rules)-1 {
		panic(fmt.Sprint("index ", index, " is larger than the number of combinations for ", rules.holes,
			" holes and ", rules.colours, " colours!"))
	}

	result := make([]int, rules.holes)
	for i := rules.holes - 1; i >= 0 && index > 0; i-- {
		colour := index % rules.colours

		result[i] = colour
		index = index / rules.colours
	}

	return result
}

func CodeToIndex(rules Rules, code []int) int {
	var result int

	for i := 0; i < rules.holes; i++ {
		result = (result * rules.colours) + code[i]
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

func NumberOfCombinations(rules Rules) int {
	return int(math.Pow(float64(rules.colours), float64(rules.holes)))
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

func FindPossibleCodesIndicies(rules Rules, facts []CodeScore) bitarray.BitArray {
	result := bitarray.NewBitArray(uint64(NumberOfCombinations(rules)))

	combinations := NumberOfCombinations(rules)
	for i := 0; i < combinations; i++ {
		guess := IndexToCode(rules, i)

		if GuessIsPossible(facts, guess) {
			if result.SetBit(uint64(i)) != nil {
				log.Panicln("Failed to set bit ", i)
			}
		}
	}

	return result
}

func PossibleScores(rules Rules) []Score {
	result := []Score{}

	for correct := 0; correct <= rules.holes; correct++ {
		for wrongPosition := 0; correct+wrongPosition <= rules.holes; wrongPosition++ {
			if !(correct == rules.holes-1 && wrongPosition == 1) {
				result = append(result, Score{correct: correct, misplaced: wrongPosition})
			}
		}
	}

	return result
}

func main() {
	rules := Rules{3, 4}
	numberOfCombinations := NumberOfCombinations(rules)
	foundScores := make(map[Score]bool)

	for i := 0; i < numberOfCombinations; i++ {
		code := IndexToCode(rules, i)

		for j := 0; j < numberOfCombinations; j++ {
			guess := IndexToCode(rules, j)

			score := CalculateScore(code, guess)

			if _, ok := foundScores[score]; !ok {
				fmt.Println("Found new score:", score)
				foundScores[score] = true
			}
		}
	}
}
