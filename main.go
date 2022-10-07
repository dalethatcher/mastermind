package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"

	"github.com/golang-collections/go-datastructures/bitarray"
)

type Rules struct {
	holes        int
	colours      int
	combinations int
}

type Score struct {
	correct   int
	misplaced int
}

type CodeScore struct {
	guess []int
	score Score
}

func NewRules(holes int, colours int) Rules {
	return Rules{
		holes:        holes,
		colours:      colours,
		combinations: int(math.Pow(float64(colours), float64(holes))),
	}
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

func IndexToCode(rules Rules, index int, code []int) {
	if index > rules.combinations-1 {
		panic(fmt.Sprint("index ", index, " is larger than the number of combinations for ", rules.holes,
			" holes and ", rules.colours, " colours!"))
	}

	for i := rules.holes - 1; i >= 0; i-- {
		if index <= 0 {
			code[i] = 0
		} else {
			colour := index % rules.colours

			code[i] = colour
			index = index / rules.colours
		}
	}
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

func GuessIsPossible(facts []CodeScore, guess []int) bool {
	for _, fact := range facts {
		score := CalculateScore(guess, fact.guess)

		if score != fact.score {
			return false
		}
	}

	return true
}

func FindPossibleCodes(rules Rules, facts []CodeScore) (int, bitarray.BitArray) {
	count := 0
	codes := bitarray.NewBitArray(uint64(rules.combinations))

	guess := make([]int, rules.holes)
	for i := 0; i < rules.combinations; i++ {
		IndexToCode(rules, i, guess)

		if GuessIsPossible(facts, guess) {
			count++
			if codes.SetBit(uint64(i)) != nil {
				log.Panicln("Failed to set bit ", i)
			}
		}
	}

	return count, codes
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

func FindMaxPossibleCountForGuess(rules Rules, facts []CodeScore, guess []int) int {
	result := 0
	possibleScores := PossibleScores(rules)

	facts = append(facts, CodeScore{guess: guess})
	lastFactIndex := len(facts) - 1
	for _, score := range possibleScores {
		facts[lastFactIndex].score = score
		count, _ := FindPossibleCodes(rules, facts)

		if count > result {
			result = count
		}
	}

	return result
}

func FindBestGuess(rules Rules, facts []CodeScore) []int {
	result := make([]int, rules.holes)

	remainingCount, remainingCandidates := FindPossibleCodes(rules, facts)
	if remainingCount == 1 || remainingCount == 2 {
		index := int(remainingCandidates.ToNums()[0])

		IndexToCode(rules, index, result)
		return result
	}

	lowestCount := rules.combinations
	guess := make([]int, rules.holes)

	for i := 0; i < rules.combinations; i++ {
		IndexToCode(rules, i, guess)
		count := FindMaxPossibleCountForGuess(rules, facts, guess)

		if count < lowestCount {
			lowestCount = count
			for h := 0; h < rules.holes; h++ {
				result[h] = guess[h]
			}
		}
	}

	return result
}

const profile = false

func main() {
	if profile {
		f, err := os.Create("profile.cpu")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	rules := NewRules(4, 6)
	code := []int{2, 5, 2, 1}
	facts := []CodeScore{}

	for {
		fmt.Println("Thinking...")
		guess := FindBestGuess(rules, facts)

		fmt.Println("Guessing", guess)
		score := CalculateScore(code, guess)
		fmt.Println("    score", score)

		if score.correct == 4 {
			fmt.Println("Code found")
			break
		}

		facts = append(facts, CodeScore{guess: guess, score: score})
		count, _ := FindPossibleCodes(rules, facts)
		fmt.Println("    remaining possibilities", count)
	}
}
