package main

import (
	"testing"
)

type TestScenario struct {
	code        []int
	guess       []int
	expectation Score
}

func TestScoreCalculation(t *testing.T) {
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
