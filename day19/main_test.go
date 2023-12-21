package main

import (
	"testing"
)

func TestCalculateValidPartCombinations(t *testing.T) {
	// Test case 1: Valid limits
	limits1 := partLimits{
		valid: true,
		min:   map[string]int{"x": 1, "m": 2, "a": 3, "s": 4},
		max:   map[string]int{"x": 5, "m": 6, "a": 7, "s": 8},
	}
	expected1 := 5 * 5 * 5 * 5
	if result := calculateValidPartCombinations(limits1); result != expected1 {
		t.Errorf("Test case 1 failed: expected %d, got %d", expected1, result)
	}

	// Test case 2: Invalid limits
	limits2 := partLimits{
		valid: false,
		min:   map[string]int{"x": 5, "m": 6, "a": 7, "s": 8},
		max:   map[string]int{"x": 1, "m": 2, "a": 3, "s": 4},
	}
	expected2 := 0
	if result := calculateValidPartCombinations(limits2); result != expected2 {
		t.Errorf("Test case 2 failed: expected %d, got %d", expected2, result)
	}

	// Add more test cases as needed...
}

func TestWalkRuleSetCRN(t *testing.T) {
	var startPoint string
	var result int
	ruleSet := makeRuleSet()

	startPoint = "crn"
	result = walkRuleSet(ruleSet, startPoint, makeMaxPartLimits(), 0, false)
	expected := 85632000000000
	if result != expected {
		t.Errorf("TestWalkRuleSet '%s': expected %d, got %d", startPoint, expected, result)
	}
}

func TestWalkRuleSetQKQ(t *testing.T) {
	var startPoint string
	var result int
	ruleSet := makeRuleSet()

	startPoint = "qkq"
	result = walkRuleSet(ruleSet, startPoint, makeMaxPartLimits(), 0, true)
	expected := 176192000000000
	if result != expected {
		t.Errorf("TestWalkRuleSet '%s': expected %d, got %d", startPoint, expected, result)
	}
}

func makeRuleSet() ruleSet {
	rs := ruleSet{}
	rs.workflows = make(map[string]workflow)

	// crn{x>2662:A,R}
	rs.workflows["crn"] = workflow{
		rules: []rule{
			{
				end:      false,
				operand:  "x",
				operator: encodeOperator(">"),
				value:    2662,
				target:   "A",
			},
			{
				end:      true,
				operand:  "",
				operator: 0,
				value:    0,
				target:   "R",
			},
		},
	}

	// qkq{x<1416:A,crn}
	rs.workflows["qkq"] = workflow{
		rules: []rule{
			{
				end:      false,
				operand:  "x",
				operator: encodeOperator("<"),
				value:    1416,
				target:   "A",
			},
			{
				end:      true,
				operand:  "",
				operator: 0,
				value:    0,
				target:   "crn",
			},
		},
	}

	return rs
}
