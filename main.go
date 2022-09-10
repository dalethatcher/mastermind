package main

import (
	"fmt"
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

func main() {
	for i := 0; i < NumberOfCombinations(4, 4); i++ {
		fmt.Println(IndexToCode(4, 4, i))
	}
}
