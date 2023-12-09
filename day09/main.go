package main

import (
	"fmt"
	"os"
	"strings"
)

// type extrapolation struct {
// 	points [2]int
// 	diff   int
// }

const OasisReportExample1 = `0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 9)

	// Load the input file
	// input := OasisReportExample1
	input := loadFileContents("oasis_report.txt")

	// Parse the input file
	report := parseReport(input)
	// spew.Dump(report)

	extrapolatedValueSum := 0
	extrapolatedBackwardValueSum := 0
	// Loop through the rows
	for i, _ := range report {
		// Find the extrapolation
		deltaStack := createDeltaStack(report[i])
		fmt.Printf("Report: %v\nDelta Stack: %v\n", report[i], deltaStack)
		calcExtrapolation(&report[i], &deltaStack)
		// fmt.Printf("Report: %v\nDelta Stack: %v\n", report[i], deltaStack)
		// fmt.Printf("Extrapolation: %v\n", report[i])
		calcExtrapolationBackwards(&report[i], &deltaStack)
		fmt.Printf("Report: %v\nDelta Stack: %v\n", report[i], deltaStack)
		// fmt.Printf("Extrapolation: %v\n", report[i])

		extrapolatedValueSum += report[i][len(report[i])-1]
		extrapolatedBackwardValueSum += report[i][0]
	}
	fmt.Printf("Extrapolated Value Sum:\t\t\t%d\n", extrapolatedValueSum)
	fmt.Printf("Extrapolated Backwards Value Sum:\t%d\n", extrapolatedBackwardValueSum)

}

func calcExtrapolationBackwards(point *[]int, deltaStack *[][]int) {
	// Find the depth of the delta stack
	depth := len(*deltaStack)

	// Loop through the deltas
	for i := depth - 2; i >= 0; i-- {
		// Find the first points of this and the next depth
		firstPointPos := 0
		firstPoint := (*deltaStack)[i][firstPointPos]
		// This should be one less than the last point
		deltaPos := 0
		delta := (*deltaStack)[i+1][deltaPos]

		// Find the extrapolation
		extrapolation := firstPoint - delta
		// Append extrapolation to current deltaStack row
		(*deltaStack)[i] = prepend((*deltaStack)[i], extrapolation)
	}
	// Find the last points of this and the next depth
	firstPointPos := 0
	lastPoint := (*point)[firstPointPos]
	deltaPos := 0
	delta := (*deltaStack)[0][deltaPos]
	// fmt.Printf("Last Point: %d, Delta: %d\n", lastPoint, delta)
	// Find the extrapolation
	extrapolation := lastPoint - delta
	// Append extrapolation to current deltaStack row
	(*point) = prepend((*point), extrapolation)
}

// Add points to the beginniong of an array
func prepend(slice []int, elements ...int) []int {
	// Create a slice to hold the result
	result := make([]int, 0, len(slice)+len(elements))

	// Append the elements
	result = append(result, elements...)

	// Append the slice
	result = append(result, slice...)

	return result
}

func calcExtrapolation(point *[]int, deltaStack *[][]int) {
	// Find the depth of the delta stack
	depth := len(*deltaStack)

	// Loop through the deltas
	for i := depth - 2; i >= 0; i-- {
		// Find the last points of this and the next depth
		lastPointPos := len((*deltaStack)[i]) - 1
		lastPoint := (*deltaStack)[i][lastPointPos]
		// This should be one less than the last point
		deltaPos := len((*deltaStack)[i+1]) - 1
		delta := (*deltaStack)[i+1][deltaPos]

		// Find the extrapolation
		extrapolation := lastPoint + delta
		// Append extrapolation to current deltaStack row
		(*deltaStack)[i] = append((*deltaStack)[i], extrapolation)
	}
	// Find the last points of this and the next depth
	lastPointPos := len(*point) - 1
	lastPoint := (*point)[lastPointPos]
	deltaPos := len((*deltaStack)[0]) - 1
	delta := (*deltaStack)[0][deltaPos]
	// fmt.Printf("Last Point: %d, Delta: %d\n", lastPoint, delta)
	// Find the extrapolation
	extrapolation := lastPoint + delta
	// Append extrapolation to current deltaStack row
	(*point) = append((*point), extrapolation)
}

func createDeltaStack(points []int) [][]int {
	// Create a slice to hold the delta stack
	deltaStack := make([][]int, 0, 5)

	// Calculate the deltas array, continue adding depths until all deltas are 0
	allZero := false
	var deltas []int
	for i := 0; !allZero; i++ {
		if i == 0 {
			allZero, deltas = calculateDifferences(points)
		} else {
			allZero, deltas = calculateDifferences(deltaStack[i-1])
		}

		deltaStack = append(deltaStack, deltas)
		// fmt.Printf("Delta Stack depth: (%d) %v\n", i, deltaStack)
	}

	return deltaStack
}

func calculateDifferences(points []int) (bool, []int) {
	allZero := true

	// Create a slice to hold the differences
	differences := make([]int, len(points)-1)

	// Loop through the points
	for i, _ := range points {
		// Skip the first point
		if i == 0 {
			continue
		}

		// Calculate the difference
		differences[i-1] = points[i] - points[i-1]
		if differences[i-1] != 0 {
			allZero = false
		}
	}

	return allZero, differences
}

// func findExtrapolations(points []int) [][]extrapolation {
//         // Create a slice to hold the extrapolations
//         extrapolations := make([][]extrapolation, 0,5)

//         // Loop through the points
//         for i, _ := range points {
//                 // Create a slice to hold the extrapolations
//                 extrapolations[i] = make([]extrapolation, len(points))

//                 // Loop through the points
//                 for j, _ := range points {
//                         // Skip the same point
//                         if i == j {
//                                 continue
//                         }

//                         // Find the extrapolation
//                         extrapolations[i][j] = findExtrapolation(points[i], points[j])
//                 }
//         }

//         return extrapolations
// }

// // Find the extrapolation between two points
// func findExtrapolation(points []int) extrapolation {
// 	result := extrapolation{}

// 	// Find the length of the array
// 	length := len(points)

// 	// Store the last two points
// 	result.points[0] = points[length-2]
// 	result.points[1] = points[length-1]

// 	// Find the difference between the last two points
// 	result.diff = points[length-1] - points[length-2]

// 	return result
// }

// Parse the input file
func parseReport(input string) [][]int {
	// Split the input into lines
	lines := splitLines(input)

	// Create a slice to hold the report
	report := make([][]int, len(lines))

	// Loop through the lines
	for i, line := range lines {
		// Skip empty lines
		if line == "" {
			continue
		}

		// Split the line into fields
		fields := splitFields(line)

		// Create a slice to hold the fields
		report[i] = make([]int, len(fields))

		// Loop through the fields
		for j, field := range fields {
			// Convert the field to an integer
			report[i][j] = parseInt(field)
		}
	}

	return report
}

// Split a string into lines
func splitLines(input string) []string {
	// Split the input into lines
	lines := strings.Split(input, "\n")

	return lines
}

// Split a string into fields
func splitFields(input string) []string {
	// Split the input into fields
	fields := strings.Split(input, " ")

	// Remove empty fields
	for i := 0; i < len(fields); i++ {
		if fields[i] == "" {
			fields = append(fields[:i], fields[i+1:]...)
		}
	}
	return fields
}

// Parse an integer from a string
func parseInt(input string) int {
	// Parse the integer
	var value int
	_, err := fmt.Sscanf(input, "%d", &value)
	if err != nil {
		panic(err)
	}

	return value
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
