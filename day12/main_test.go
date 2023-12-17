package main

import (
	"reflect"
	"testing"
)

// Test for classifySpringMap
func TestClassifySpringMap(t *testing.T) {
	// Define the test cases
	var testCases = []struct {
		springMap  string
		springList []int
	}{
		{"#.#.###", []int{1, 1, 3}},
		{".#...#....###.", []int{1, 1, 3}},
		{".#.###.#.######", []int{1, 3, 1, 6}},
		{"####.#...#...", []int{4, 1, 1}},
		{"#....######..#####.", []int{1, 6, 5}},
		{".###.##....#", []int{3, 2, 1}},
	}

	// Loop over the test cases
	for _, tc := range testCases {
		// Call classifySpringMap
		springList := classifySpringMap(tc.springMap)
		// Check the returned spring list
		if !reflect.DeepEqual(springList, tc.springList) {
			t.Errorf("classifySpringMap(%s) = %v, want %v", tc.springMap, springList, tc.springList)
		}
	}
}

// Test for twoToPowerOf
func TestTwoToPowerOf(t *testing.T) {
	// Define the test cases
	var testCases = []struct {
		n        int
		expected int
	}{
		{0, 1},
		{1, 2},
		{2, 4},
		{3, 8},
		{4, 16},
		{5, 32},
		{6, 64},
		{7, 128},
		{8, 256},
		{9, 512},
		{10, 1024},
		{11, 2048},
		{12, 4096},
		{13, 8192},
		{14, 16384},
	}

	// Loop over the test cases
	for _, tc := range testCases {
		// Call twoToPowerOf
		result := twoToPowerOf(tc.n)
		// Check the returned result
		if result != tc.expected {
			t.Errorf("twoToPowerOf(%d) = %d, want %d", tc.n, result, tc.expected)
		}
	}
}

// Test for computePossibleSpringMapCombinationsRecursive
func TestComputePossibleSpringMapCombinationsRecursive(t *testing.T) {
	ttable := buildSpringMapTruthTable()

	// Define the test cases
	var testCases = []struct {
		springMap *springMapEntry
		expected  int
	}{
		// Test case 1
		{
			springMap: &springMapEntry{
				// Define the springMapEntry fields here
				springMap:  "???.###",
				springList: []int{1, 1, 3},
			},

			expected: 1, // Define the expected result here
		},
		// Test case 2
		{
			springMap: &springMapEntry{
				// Define the springMapEntry fields here
				springMap:  ".??..??...?##.",
				springList: []int{1, 1, 3},
			},

			expected: 4, // Define the expected result here
		},
		// Test case 3
		{
			springMap: &springMapEntry{
				// Define the springMapEntry fields here
				springMap:  "?#?#?#?#?#?#?#?",
				springList: []int{1, 3, 1, 6},
			},

			expected: 1, // Define the expected result here
		},
		// Test case 4
		{
			springMap: &springMapEntry{
				// Define the springMapEntry fields here
				springMap:  "????.#...#...",
				springList: []int{4, 1, 1},
			},

			expected: 1, // Define the expected result here
		},
		// Test case 5
		{
			springMap: &springMapEntry{
				// Define the springMapEntry fields here
				springMap:  "????.######..#####.",
				springList: []int{1, 6, 5},
			},

			expected: 4, // Define the expected result here
		},
		// Test case 6
		{
			springMap: &springMapEntry{
				// Define the springMapEntry fields here
				springMap:  "?###????????",
				springList: []int{3, 2, 1},
			},

			expected: 10, // Define the expected result here
		},

		// Add more test cases as needed
	}

	// Loop over the test cases
	for _, tc := range testCases {
		// Call computePossibleSpringMapCombinationsRecursive
		result := computePossibleSpringMapCombinationsRecursive(tc.springMap, &ttable, true)
		// Check the returned result
		if result != tc.expected {
			t.Errorf("computePossibleSpringMapCombinationsRecursive(%v) = %d, want %d", tc.springMap, result, tc.expected)
		}
	}
}
