package main

import (
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateScore(t *testing.T) {
	type TestScenario struct {
		code        []int
		guess       []int
		expectation Score
	}

	scenarios := []TestScenario{
		{code: []int{1, 2, 3, 4}, guess: []int{1, 1, 1, 1}, expectation: Score{correct: 1}},
		{code: []int{1, 2, 2, 2}, guess: []int{3, 1, 3, 3}, expectation: Score{misplaced: 1}},
		{code: []int{1, 2, 3, 4}, guess: []int{1, 1, 2, 5}, expectation: Score{correct: 1, misplaced: 1}},
		{code: []int{1, 1, 2, 2}, guess: []int{1, 2, 1, 2}, expectation: Score{correct: 2, misplaced: 2}},
	}

	for _, scenario := range scenarios {
		score := CalculateScore(scenario.code, scenario.guess)

		if score != scenario.expectation {
			t.Error("Expected score", scenario.expectation, "but got", score, "for code", scenario.code, "and guess", scenario.guess)
		}
	}
}

func TestIndexToCode(t *testing.T) {
	type TestScenario struct {
		colours     int
		index       int
		expectation []int
	}

	scenarios := []TestScenario{
		{colours: 4, index: 0, expectation: []int{0, 0, 0, 0}},
		{colours: 4, index: 1, expectation: []int{0, 0, 0, 1}},
		{colours: 4, index: 255, expectation: []int{3, 3, 3, 3}},
		{colours: 10, index: 4, expectation: []int{0, 0, 0, 4}},
	}

	for _, scenario := range scenarios {
		code := IndexToCode(Rules{len(scenario.expectation), scenario.colours}, scenario.index)

		if !reflect.DeepEqual(code, scenario.expectation) {
			t.Error("Expected code", scenario.expectation, "but got", code, "for index",
				scenario.index, "and", scenario.colours, "colours")
		}
	}
}

func TestGuessIsPossible(t *testing.T) {
	type TestScenario struct {
		facts []CodeScore
		guess []int
	}

	possible := []TestScenario{
		{facts: []CodeScore{{guess: []int{0, 1, 2, 3}, score: Score{correct: 1}}}, guess: []int{1, 1, 1, 1}},
		{facts: []CodeScore{{guess: []int{1, 2, 2, 2}, score: Score{misplaced: 1}}}, guess: []int{0, 1, 0, 0}},
	}

	for _, scenario := range possible {
		if !GuessIsPossible(scenario.facts, scenario.guess) {
			t.Error("Expected", scenario.guess, "to be valid for facts:", scenario.facts)
		}
	}
}

func TestFindKnuthPaperSolution(t *testing.T) {
	facts := []CodeScore{
		{guess: []int{0, 0, 1, 1}, score: Score{correct: 1}},
		{guess: []int{0, 2, 3, 3}, score: Score{misplaced: 1}},
		{guess: []int{2, 4, 1, 5}, score: Score{correct: 1, misplaced: 2}},
		{guess: []int{0, 3, 5, 1}, score: Score{correct: 1, misplaced: 1}},
	}

	rules := Rules{4, 6}
	combinations := NumberOfCombinations(rules)
	matches := [][]int{}
	for i := 0; i < combinations; i++ {
		guess := IndexToCode(rules, i)

		if GuessIsPossible(facts, guess) {
			matches = append(matches, guess)
		}
	}

	if !reflect.DeepEqual(matches, [][]int{{2, 5, 2, 1}}) {
		t.Error("Expected only match to be 2521 but got", matches)
	}
}

func TestCodeToIndex(t *testing.T) {
	rules := Rules{2, 2}
	for i := 0; i < 4; i++ {
		code := IndexToCode(rules, i)
		index := CodeToIndex(rules, code)

		if index != i {
			t.Error("Expected code", code, "to have index", i, "but got", index)
		}
	}
}

func TestFindPossibleCodesIndicies(t *testing.T) {
	facts := []CodeScore{
		{guess: []int{1, 1}, score: Score{correct: 1}},
	}

	count, possibleCodes := FindPossibleCodes(Rules{2, 2}, facts)

	expectation := []uint64{1, 2}
	assert.Equal(t, expectation, possibleCodes.ToNums(), "code indices should match")
	assert.Equal(t, 2, count, "possible codes count")
}

func TestFindMaxPossibleCountForGuess(t *testing.T) {
	facts := []CodeScore{
		{guess: []int{0, 0, 1, 1}, score: Score{correct: 1}},
		{guess: []int{0, 2, 3, 3}, score: Score{misplaced: 1}},
		{guess: []int{2, 4, 1, 5}, score: Score{correct: 1, misplaced: 2}},
	}
	guess := []int{0, 3, 5, 1}

	rules := Rules{4, 6}
	result := FindMaxPossibleCountForGuess(rules, facts, guess)

	assert.Equal(t, 1, result, "expected max count of one")
}

func TestFindBestGuess(t *testing.T) {
	facts := []CodeScore{
		{guess: []int{0, 0, 1, 1}, score: Score{correct: 1}},
		{guess: []int{0, 2, 3, 3}, score: Score{misplaced: 1}},
		{guess: []int{2, 4, 1, 5}, score: Score{correct: 1, misplaced: 2}},
	}
	rules := Rules{4, 6}

	result := FindBestGuess(rules, facts)
	assert.ElementsMatch(t, []int{0, 0, 5, 1}, result)
}

func TestPossibleScores(t *testing.T) {
	result := PossibleScores(Rules{3, 4})
	expectation := []Score{
		{},
		{misplaced: 1},
		{misplaced: 2},
		{misplaced: 3},
		{correct: 1},
		{correct: 1, misplaced: 1},
		{correct: 1, misplaced: 2},
		{correct: 2},
		{correct: 3},
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].correct < result[j].correct ||
			(result[i].correct == result[j].correct &&
				result[i].misplaced < result[j].misplaced)
	})

	assert.Equal(t, expectation, result, "did not get expected scores")
}
