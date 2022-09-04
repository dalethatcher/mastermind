package main

import (
	"reflect"
	"testing"
)

func TestScoreCalculation(t *testing.T) {
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
