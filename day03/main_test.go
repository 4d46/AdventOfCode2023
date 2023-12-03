package main

import "testing"

func TestIsPartAdjacentToLabel(t *testing.T) {
	// Define a test case
	part := Part{name: '@', position: Coordinate{5, 5}}
	label := Label{name: "123", start: Coordinate{6, 5}, end: Coordinate{9, 5}}

	// Call the function with the test case
	result := isPartAdjactentToLabel(part, label)

	// Check the result
	if result != true {
		t.Errorf("Expected true, but got %v", result)
	}
}

func TestNotYIsPartAdjacentToLabel(t *testing.T) {
	// Define a test case
	part := Part{name: '@', position: Coordinate{5, 5}}
	label := Label{name: "123", start: Coordinate{6, 7}, end: Coordinate{9, 7}}

	// Call the function with the test case
	result := isPartAdjactentToLabel(part, label)

	// Check the result
	if result != false {
		t.Errorf("Expected true, but got %v", result)
	}
}

func TestNotXIsPartAdjacentToLabel(t *testing.T) {
	// Define a test case
	part := Part{name: '@', position: Coordinate{5, 5}}
	label := Label{name: "123", start: Coordinate{7, 6}, end: Coordinate{10, 6}}

	// Call the function with the test case
	result := isPartAdjactentToLabel(part, label)

	// Check the result
	if result != false {
		t.Errorf("Expected false, but got %v", result)
	}
}

func TestNot515IsPartAdjacentToLabel(t *testing.T) {
	// Define a test case
	part := Part{name: '*', position: Coordinate{62, 1}}
	label := Label{name: "515", start: Coordinate{58, 1}, end: Coordinate{58 + len("515") - 1, 1}}

	// Call the function with the test case
	result := isPartAdjactentToLabel(part, label)

	// Check the result
	if result != false {
		t.Errorf("Expected false, but got %v", result)
	}
}
