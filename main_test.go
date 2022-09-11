package main

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"sort"
	"testing"
)

func TestCalculateScore(t *testing.T) {
	type TestScenario struct {
		code        []int
		guess       []int
		expectation Score
	}

	scenarios := []TestScenario{
		{code: []int{1, 2, 3, 4}, guess: []int{1, 1, 1, 1}, expectation: Score{rightValueAndPosition: 1}},
		{code: []int{1, 2, 2, 2}, guess: []int{3, 1, 3, 3}, expectation: Score{rightValueWrongPosition: 1}},
		{code: []int{1, 2, 3, 4}, guess: []int{1, 1, 2, 5}, expectation: Score{rightValueAndPosition: 1, rightValueWrongPosition: 1}},
		{code: []int{1, 1, 2, 2}, guess: []int{1, 2, 1, 2}, expectation: Score{rightValueAndPosition: 2, rightValueWrongPosition: 2}},
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
		numberOfColours int
		index           int
		expectation     []int
	}

	scenarios := []TestScenario{
		{numberOfColours: 4, index: 0, expectation: []int{0, 0, 0, 0}},
		{numberOfColours: 4, index: 1, expectation: []int{0, 0, 0, 1}},
		{numberOfColours: 4, index: 255, expectation: []int{3, 3, 3, 3}},
		{numberOfColours: 10, index: 4, expectation: []int{0, 0, 0, 4}},
	}

	for _, scenario := range scenarios {
		code := IndexToCode(len(scenario.expectation), scenario.numberOfColours, scenario.index)

		if !reflect.DeepEqual(code, scenario.expectation) {
			t.Error("Expected code", scenario.expectation, "but got", code, "for index",
				scenario.index, "and", scenario.numberOfColours, "colours")
		}
	}
}

func TestGuessIsPossible(t *testing.T) {
	type TestScenario struct {
		facts []CodeScore
		guess []int
	}

	possible := []TestScenario{
		{facts: []CodeScore{{guess: []int{0, 1, 2, 3}, score: Score{rightValueAndPosition: 1}}}, guess: []int{1, 1, 1, 1}},
		{facts: []CodeScore{{guess: []int{1, 2, 2, 2}, score: Score{rightValueWrongPosition: 1}}}, guess: []int{0, 1, 0, 0}},
	}

	for _, scenario := range possible {
		if !GuessIsPossible(scenario.facts, scenario.guess) {
			t.Error("Expected", scenario.guess, "to be valid for facts:", scenario.facts)
		}
	}
}

func TestFindKnuthPaperSolution(t *testing.T) {
	facts := []CodeScore{
		{guess: []int{0, 0, 1, 1}, score: Score{rightValueAndPosition: 1}},
		{guess: []int{0, 2, 3, 3}, score: Score{rightValueWrongPosition: 1}},
		{guess: []int{2, 4, 1, 5}, score: Score{rightValueAndPosition: 1, rightValueWrongPosition: 2}},
		{guess: []int{0, 3, 5, 1}, score: Score{rightValueAndPosition: 1, rightValueWrongPosition: 1}},
	}

	combinations := NumberOfCombinations(4, 6)
	matches := [][]int{}
	for i := 0; i < combinations; i++ {
		guess := IndexToCode(4, 6, i)

		if GuessIsPossible(facts, guess) {
			matches = append(matches, guess)
		}
	}

	if !reflect.DeepEqual(matches, [][]int{{2, 5, 2, 1}}) {
		t.Error("Expected only match to be 2521 but got", matches)
	}
}

func TestCodeToIndex(t *testing.T) {
	for i := 0; i < 4; i++ {
		code := IndexToCode(2, 2, i)
		index := CodeToIndex(2, 2, code)

		if index != i {
			t.Error("Expected code", code, "to have index", i, "but got", index)
		}
	}
}

func TestFindPossibleCodesIndicies(t *testing.T) {
	facts := []CodeScore{
		{guess: []int{1, 1}, score: Score{rightValueAndPosition: 1}},
	}

	possibleCodes := FindPossibleCodesIndicies(2, 2, facts)

	expectation := []uint64{1, 2}
	assert.Equal(t, expectation, possibleCodes.ToNums(), "code indices should match")
}

func TestPossibleNonFinalScores(t *testing.T) {
	result := PossibleScores(3)
	expectation := []Score{
		{},
		{rightValueWrongPosition: 1},
		{rightValueWrongPosition: 2},
		{rightValueWrongPosition: 3},
		{rightValueAndPosition: 1},
		{rightValueAndPosition: 1, rightValueWrongPosition: 1},
		{rightValueAndPosition: 1, rightValueWrongPosition: 2},
		{rightValueAndPosition: 2},
		{rightValueAndPosition: 3},
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].rightValueAndPosition < result[j].rightValueAndPosition ||
			(result[i].rightValueAndPosition == result[j].rightValueAndPosition &&
				result[i].rightValueWrongPosition < result[j].rightValueWrongPosition)
	})

	assert.Equal(t, expectation, result, "did not get expected scores")
}
