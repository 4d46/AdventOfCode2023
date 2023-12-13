package main

import (
	"fmt"
	"os"
	"strings"
)

type pattern struct {
	pattern                 string
	rowPatterns             []uint64
	columnsPatterns         []uint64
	hReflectionPoint        int
	vReflectionPoint        int
	hSmudgedReflectionPoint int
	vSmudgedReflectionPoint int
}

var mirrorExample1 = `#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.

#...##..#
#....#..#
..##..###
#####.##.
#####.##.
..##..###
#....#..#`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 13)

	// Load input
	// inputStr := mirrorExample1
	inputStr := loadFileContents("patterns.txt")

	// Parse input
	patterns := parseInput(inputStr)

	// spew.Dump("patterns: %v\n", patterns)
	fmt.Println()
	summarySum1 := 0
	for _, pattern := range patterns {
		fmt.Println()
		printPattern(pattern, pattern.hReflectionPoint, pattern.vReflectionPoint)
		summary := calculateSummary(pattern)
		summarySum1 += summary
		fmt.Printf("%d\n", summary)
		// fmt.Printf("Summary: %d\n", summary)
	}
	fmt.Printf("Part 1 summary sum: %d\n", summarySum1)

	fmt.Println()
	summarySum2 := 0
	for _, pattern := range patterns {
		fmt.Println()
		printPattern(pattern, pattern.hSmudgedReflectionPoint, pattern.vSmudgedReflectionPoint)
		summary := calculateSmudgedSummary(pattern)
		summarySum2 += summary
		fmt.Printf("%d\n", summary)
		// fmt.Printf("Summary: %d\n", summary)
	}
	fmt.Printf("Part 2 summary sum: %d\n", summarySum2)

}

// Calculate summary int for a pattern
func calculateSummary(p pattern) int {
	// Calculate summary int for a pattern
	// This is the sum of the decimal values of the uint64s in the pattern
	if p.hReflectionPoint == -1 && p.vReflectionPoint == -1 {
		panic("No reflection point found")
	} else if p.hReflectionPoint >= 0 && p.vReflectionPoint >= 0 {
		panic("Both reflection points found")
	}
	summary := (p.vReflectionPoint+1)*100 + (p.hReflectionPoint + 1)
	return summary
}

// Calculate smudged summary int for a pattern
func calculateSmudgedSummary(p pattern) int {
	// Calculate summary int for a pattern
	// This is the sum of the decimal values of the uint64s in the pattern
	if p.hSmudgedReflectionPoint == -1 && p.vSmudgedReflectionPoint == -1 {
		panic("No reflection point found")
	}
	// } else if p.hSmudgedReflectionPoint >= 0 && p.vSmudgedReflectionPoint >= 0 {
	// 	panic("Both reflection points found")
	// }
	summary := (p.vSmudgedReflectionPoint+1)*100 + (p.hSmudgedReflectionPoint + 1)
	return summary
}

// Find vertical reflection point in poattern
func findVerticalReflectionPoint(pattern pattern) int {
	length := len(pattern.rowPatterns)
	// Loop over candidate reflection points, from the first to the last but one row
	for i := 0; i < length-1; i++ {
		// Loop from this point, backwards and forwards,checking for a match
		// If we find a point that doesn't match this reflection point isn't valid, move onto the next one
		match := true
		for j := 0; i-j >= 0 && i+j+1 < length; j++ {
			if pattern.rowPatterns[i-j] != pattern.rowPatterns[i+j+1] {
				match = false
				break
			}
		}
		// We've found the match, return the reflection point
		if match {
			return i
		}
		// Otherwise try next flreflection point
	}
	// If we are here no reflection point was found
	// Return -1 to indicate this
	return -1
}

// Find horizontal reflection point in poattern
func findHorizontalReflectionPoint(pattern pattern) int {
	// spew.Dump("pattern: %v\n", pattern)
	length := len(pattern.columnsPatterns)
	// fmt.Printf("LENGTH: %d\n", length)
	// Loop over candidate reflection points, from the first to the last but one row (which represents columns)
	for i := 0; i < length-1; i++ {
		// Loop from this point, backwards and forwards,checking for a match
		// If we find a point that doesn't match this reflection point isn't valid, move onto the next one
		match := true
		for j := 0; i-j >= 0 && i+j+1 < length; j++ {
			fmt.Printf("i: %d, j: %d %d->%d\n", i, j, pattern.columnsPatterns[i-j], pattern.columnsPatterns[i+j+1])
			if pattern.columnsPatterns[i-j] != pattern.columnsPatterns[i+j+1] {
				match = false
				break
			}
		}
		// We've found the match, return the reflection point
		if match {
			return i
		}
		// Otherwise try next reflection point
	}
	// If we are here no reflection point was found
	// Return -1 to indicate this
	return -1
}

// Find Smudged vertical reflection point in poattern
// Check for a reflection point that is smudged by one point
func findSmudgedVerticalReflectionPoint(pattern pattern) int {
	length := len(pattern.rowPatterns)
	// Loop over candidate reflection points, from the first to the last but one row
	for i := 0; i < length-1; i++ {
		// Loop from this point, backwards and forwards, checking for number of differnces
		// If we find a point that doesn't match this reflection point isn't valid, increase the count by 1
		diffCount := 0
		for j := 0; i-j >= 0 && i+j+1 < length; j++ {
			diffCount += countBitsDifferent(pattern.rowPatterns[i-j], pattern.rowPatterns[i+j+1])
		}
		// We've found the match, return the reflection point
		if diffCount == 1 {
			return i
		}
		// Otherwise try next reflection point
	}
	// If we are here no reflection point was found
	// Return -1 to indicate this
	return -1
}

// Find Smudged horizontal reflection point in poattern
func findSmudgedHorizontalReflectionPoint(pattern pattern) int {
	// spew.Dump("pattern: %v\n", pattern)
	length := len(pattern.columnsPatterns)
	// fmt.Printf("LENGTH: %d\n", length)
	// Loop over candidate reflection points, from the first to the last but one row (which represents columns)
	for i := 0; i < length-1; i++ {
		// Loop from this point, backwards and forwards, checking for number of differnces
		// If we find a point that doesn't match this reflection point isn't valid, increase the count by 1
		diffCount := 0
		for j := 0; i-j >= 0 && i+j+1 < length; j++ {
			// fmt.Printf("i: %d, j: %d %d->%d\n", i, j, pattern.columnsPatterns[i-j], pattern.columnsPatterns[i+j+1])
			diffCount += countBitsDifferent(pattern.columnsPatterns[i-j], pattern.columnsPatterns[i+j+1])
		}
		// We've found the match, return the reflection point
		if diffCount == 1 {
			return i
		}
		// Otherwise try next reflection point
	}
	// If we are here no reflection point was found
	// Return -1 to indicate this
	return -1
}

// Count number of bits different between two uint64s
func countBitsDifferent(uint64a, uint64b uint64) int {
	diffCount := 0
	for i := 0; i < 64; i++ {
		if uint64a&1 != uint64b&1 {
			diffCount++
		}
		uint64a >>= 1
		uint64b >>= 1
	}
	return diffCount
}

// Print a pattern
func printPattern(pattern pattern, hReflectionPoint, vReflectionPoint int) {
	// reflectionPointHorizontal := pattern.hReflectionPoint
	// reflectionPointVertical := pattern.vReflectionPoint
	reflectionPointHorizontal := hReflectionPoint
	reflectionPointVertical := vReflectionPoint
	if reflectionPointHorizontal >= 0 {
		reflectionPointHorizontal++
	}
	if reflectionPointVertical >= 0 {
		reflectionPointVertical++
	}

	// Split pattern into rows
	rows := strings.Split(pattern.pattern, "\n")

	if reflectionPointHorizontal != -1 && reflectionPointVertical != -1 {
		// Vertical reflection
		// Print rows before reflection point
		for i := 0; i < reflectionPointVertical; i++ {
			fmt.Println(rows[i][:reflectionPointHorizontal] + "|" + rows[i][reflectionPointHorizontal:])
		}
		// Print reflection point
		// Calculate row length
		rowLength := len(rows[0])
		if reflectionPointHorizontal >= 0 {
			rowLength++
		}
		fmt.Println(strings.Repeat("-", rowLength))
		// Print rows after reflection point
		for i := reflectionPointVertical; i < len(rows); i++ {
			fmt.Println(rows[i][:reflectionPointHorizontal] + "|" + rows[i][reflectionPointHorizontal:])
		}
	} else if reflectionPointVertical != -1 {
		// Vertical reflection
		// Print rows before reflection point
		for i := 0; i < reflectionPointVertical; i++ {
			fmt.Println(rows[i])
		}
		// Print reflection point
		fmt.Println(strings.Repeat("-", len(rows[0])))
		// Print rows after reflection point
		for i := reflectionPointVertical; i < len(rows); i++ {
			fmt.Println(rows[i])
		}
	} else if reflectionPointHorizontal != -1 {
		// Horizontal reflection
		// Print each row, inserting a character at the point of the  reflection point
		for _, row := range rows {
			fmt.Println(row[:reflectionPointHorizontal] + "|" + row[reflectionPointHorizontal:])
		}
	} else {
		// No reflection, just print the pattern
		for _, row := range rows {
			fmt.Println(row)
		}
	}
}

// Parse input string into a slice of patterns
func parseInput(inputStr string) []pattern {
	// Split input into lines
	lines := strings.Split(inputStr, "\n")

	// Create slice of patterns
	var patterns []pattern

	inPattern := false
	var nextPattern pattern
	// Parse each line into a pattern
	for _, line := range lines {
		if len(line) != 0 {
			if inPattern {
				// Add line to pattern
				nextPattern.pattern += "\n" + line
			} else {
				// Start new pattern
				nextPattern = pattern{pattern: line}
				inPattern = true
			}
		} else {
			// End of pattern
			inPattern = false
			// Add pattern to slice
			patterns = append(patterns, nextPattern)
			nextPattern = pattern{}
		}
	}
	if nextPattern.pattern != "" {
		// Add pattern to slice
		patterns = append(patterns, nextPattern)
	}

	for i, _ := range patterns {
		patterns[i].rowPatterns = parseRowPatterns(patterns[i].pattern)
		patterns[i].columnsPatterns = parseColumnPatterns(patterns[i].pattern)
		patterns[i].hReflectionPoint = findHorizontalReflectionPoint(patterns[i])
		patterns[i].vReflectionPoint = findVerticalReflectionPoint(patterns[i])
		patterns[i].hSmudgedReflectionPoint = findSmudgedHorizontalReflectionPoint(patterns[i])
		patterns[i].vSmudgedReflectionPoint = findSmudgedVerticalReflectionPoint(patterns[i])
	}

	return patterns
}

// Parse pattern rows into an array of uint64
// Each row should be used as a bit pattern for a uint64
func parseRowPatterns(pattern string) []uint64 {
	// Split pattern into rows
	rows := strings.Split(pattern, "\n")

	// Create array of uint64
	var uint64s []uint64

	// Parse each row into a uint64
	for _, row := range rows {
		// Convert row into uint64
		uint64s = append(uint64s, parseRow(row))
	}

	return uint64s
}

// Parse pattern columns into an array of uint64
// Each column should be used as a bit pattern for a uint64
func parseColumnPatterns(pattern string) []uint64 {
	// Split pattern into rows
	rows := strings.Split(pattern, "\n")

	// Create array of uint64
	var uint64s []uint64

	// Parse each column into a uint64
	for i := 0; i < len(rows[0]); i++ {
		var uint64 uint64
		for _, row := range rows {
			uint64 <<= 1
			if row[i] == '#' {
				uint64 |= 1
			}
		}
		uint64s = append(uint64s, uint64)
	}

	return uint64s
}

// Parse pattern row into a uint64
// Each row should be used as a bit pattern for a uint64
func parseRow(row string) uint64 {
	// Convert row into uint64
	var uint64 uint64
	for _, char := range row {
		uint64 <<= 1
		if char == '#' {
			uint64 |= 1
		}
	}

	return uint64
}

// Load file contents into a string and return it
func loadFileContents(filename string) string {
	// Read contents of file into a string
	fileBytes, err := os.ReadFile(filename) // just pass the file name
	if err != nil {
		panic(err)
	}

	return string(fileBytes) // convert content to a 'string'
}
