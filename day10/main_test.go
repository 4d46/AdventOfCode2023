package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// | is a vertical pipe connecting north and south.
// - is a horizontal pipe connecting east and west.
// L is a 90-degree bend connecting north and east.
// J is a 90-degree bend connecting north and west.
// 7 is a 90-degree bend connecting south and west.
// F is a 90-degree bend connecting south and east.
// pipeMap1 := []string{
// 	".....",
// 	".....",
// 	".....",
// 	".....",
// 	".....",
// }

func TestFindNextStep(t *testing.T) {
	{
		pipeMap := []string{
			".....",
			".7...",
			".|...",
			".L...",
			".....",
		}

		// Test case 1: current step is '|', last step is south
		lastStep := step{letter: '7', x: 1, y: 1}
		currentStep := step{letter: '|', x: 1, y: 2}
		expected := step{letter: 'L', x: 1, y: 3}
		result := findNextStep(pipeMap, lastStep, currentStep)
		if result != expected {
			t.Errorf("Expected next %c step to be %v, but got %v", currentStep.letter, expected, result)
		}
		// Test case 2: current step is '|', last step is north
		lastStep = step{letter: 'L', x: 1, y: 3}
		currentStep = step{letter: '|', x: 1, y: 2}
		expected = step{letter: '7', x: 1, y: 1}
		result = findNextStep(pipeMap, lastStep, currentStep)
		if result != expected {
			t.Errorf("Expected next %c step to be %v, but got %v", currentStep.letter, expected, result)
		}
	}

	{
		pipeMap := []string{
			".....",
			".F-J.",
			".....",
		}

		// Test case 3: current step is '-', last step is east
		lastStep := step{letter: 'F', x: 1, y: 1}
		currentStep := step{letter: '-', x: 2, y: 1}
		expected := step{letter: 'J', x: 3, y: 1}
		result := findNextStep(pipeMap, lastStep, currentStep)
		if result != expected {
			t.Errorf("Expected next %c step to be %v, but got %v", currentStep.letter, expected, result)
		}
		// Test case 4: current step is '-', last step is west
		lastStep = step{letter: 'J', x: 3, y: 1}
		currentStep = step{letter: '-', x: 2, y: 1}
		expected = step{letter: 'F', x: 1, y: 1}
		result = findNextStep(pipeMap, lastStep, currentStep)
		if result != expected {
			t.Errorf("Expected next %c step to be %v, but got %v", currentStep.letter, expected, result)
		}
	}

	{
		pipeMap := []string{
			".....",
			".|...",
			".L-..",
			".....",
		}

		// Test case 5: current step is 'L', last step is west
		lastStep := step{letter: '-', x: 2, y: 2}
		currentStep := step{letter: 'L', x: 1, y: 2}
		expected := step{letter: '|', x: 1, y: 1}
		result := findNextStep(pipeMap, lastStep, currentStep)
		if result != expected {
			t.Errorf("Expected next %c step to be %v, but got %v", currentStep.letter, expected, result)
		}
		// Test case 6: current step is 'L', last step is south
		lastStep = step{letter: '|', x: 1, y: 1}
		currentStep = step{letter: 'L', x: 1, y: 2}
		expected = step{letter: '-', x: 2, y: 2}
		result = findNextStep(pipeMap, lastStep, currentStep)
		if result != expected {
			t.Errorf("Expected next %c step to be %v, but got %v", currentStep.letter, expected, result)
		}
	}

	{
		pipeMap := []string{
			".....",
			"..|..",
			".-J..",
			".....",
		}

		// Test case 7: current step is 'J', last step is east
		lastStep := step{letter: '-', x: 1, y: 2}
		currentStep := step{letter: 'J', x: 2, y: 2}
		expected := step{letter: '|', x: 2, y: 1}
		result := findNextStep(pipeMap, lastStep, currentStep)
		if result != expected {
			t.Errorf("Expected next %c step to be %v, but got %v", currentStep.letter, expected, result)
		}
		// Test case 8: current step is 'J', last step is south
		lastStep = step{letter: '|', x: 2, y: 1}
		currentStep = step{letter: 'J', x: 2, y: 2}
		expected = step{letter: '-', x: 1, y: 2}
		result = findNextStep(pipeMap, lastStep, currentStep)
		if result != expected {
			t.Errorf("Expected next %c step to be %v, but got %v", currentStep.letter, expected, result)
		}
	}

	{
		pipeMap := []string{
			".....",
			".-7..",
			"..|..",
			".....",
		}

		// Test case 9: current step is '7', last step is east
		lastStep := step{letter: '-', x: 1, y: 1}
		currentStep := step{letter: '7', x: 2, y: 1}
		expected := step{letter: '|', x: 2, y: 2}
		result := findNextStep(pipeMap, lastStep, currentStep)
		if result != expected {
			t.Errorf("Expected next %c step to be %v, but got %v", currentStep.letter, expected, result)
		}
		// Test case 10: current step is '7', last step is north
		lastStep = step{letter: '|', x: 2, y: 2}
		currentStep = step{letter: '7', x: 2, y: 1}
		expected = step{letter: '-', x: 1, y: 1}
		result = findNextStep(pipeMap, lastStep, currentStep)
		if result != expected {
			t.Errorf("Expected next %c step to be %v, but got %v", currentStep.letter, expected, result)
		}
	}

	{
		pipeMap := []string{
			".....",
			"..F-.",
			"..|..",
			".....",
		}

		// Test case 11: current step is 'F', last step is north
		lastStep := step{letter: '|', x: 2, y: 2}
		currentStep := step{letter: 'F', x: 2, y: 1}
		expected := step{letter: '-', x: 3, y: 1}
		result := findNextStep(pipeMap, lastStep, currentStep)
		if result != expected {
			t.Errorf("Expected next %c step to be %v, but got %v", currentStep.letter, expected, result)
		}
		// Test case 12: current step is 'F', last step is west
		lastStep = step{letter: '-', x: 3, y: 1}
		currentStep = step{letter: 'F', x: 2, y: 1}
		expected = step{letter: '|', x: 2, y: 2}
		result = findNextStep(pipeMap, lastStep, currentStep)
		if result != expected {
			t.Errorf("Expected next %c step to be %v, but got %v", currentStep.letter, expected, result)
		}
	}
}

func TestFindPath1(t *testing.T) {
	const pipeMapExample1 = `.....
.S-7.
.|.|.
.L-J.
.....`

	expectedPipeMapExample1 := []step{{83, 1, 1}, {45, 2, 1}, {55, 3, 1}, {124, 3, 2}, {74, 3, 3}, {45, 2, 3}, {76, 1, 3}, {124, 1, 2}, {83, 1, 1}}

	start := step{letter: 'S', x: 1, y: 1}
	exampleMap := parsePipeMap(pipeMapExample1)
	foundPath := findPath(exampleMap, start)
	assert.Equal(t, expectedPipeMapExample1, foundPath, "paths should match")
	// t.Errorf("Expected path to be \n%s\n, but got \n%s\n", expectedPipeMapExample1, foundPath)
}

func TestFindPath1Noisy(t *testing.T) {
	const pipeMapExample = `-L|F7
7S-7|
L|7||
-L-J|
L|-JF`

	expectedPipeMapExample := []step{{83, 1, 1}, {45, 2, 1}, {55, 3, 1}, {124, 3, 2}, {74, 3, 3}, {45, 2, 3}, {76, 1, 3}, {124, 1, 2}, {83, 1, 1}}

	start := step{letter: 'S', x: 1, y: 1}
	exampleMap := parsePipeMap(pipeMapExample)
	foundPath := findPath(exampleMap, start)
	assert.Equal(t, expectedPipeMapExample, foundPath, "paths should match")
	// t.Errorf("Expected path to be \n%s\n, but got \n%s\n", expectedPipeMapExample1, foundPath)
}

func TestFindPath2Noisy(t *testing.T) {
	const pipeMapExample = `7-F7-
.FJ|7
SJLL7
|F--J
LJ.LJ`

	expectedPipeMapExample := []step{{83, 0, 2}, {74, 1, 2}, {70, 1, 1}, {74, 2, 1}, {70, 2, 0}, {55, 3, 0}, {124, 3, 1}, {76, 3, 2}, {55, 4, 2}, {74, 4, 3}, {45, 3, 3}, {45, 2, 3}, {70, 1, 3}, {74, 1, 4}, {76, 0, 4}, {124, 0, 3}, {83, 0, 2}}

	start := step{letter: 'S', x: 0, y: 2}
	exampleMap := parsePipeMap(pipeMapExample)
	foundPath := findPath(exampleMap, start)
	assert.Equal(t, expectedPipeMapExample, foundPath, "paths should match")
}
