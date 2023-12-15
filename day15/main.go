package main

import (
	"fmt"
	"os"
	"strings"
)

type step struct {
	value string
	hash  int
}

const exampleInitialisationSequence1 = `HASH`
const exampleInitialisationSequence2 = `rn=1,cm-,qp=3,cm=2,qp-,pc=4,ot=9,ab=5,pc-,pc=6,ot=7`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 15)

	// Load the input data
	// input := exampleInitialisationSequence1
	// input := exampleInitialisationSequence2
	input := loadFileContents("initialization_sequence.txt")

	// Parse the input data
	steps := parseInput(input)

	// Part 1
	calculatePartHashs(&steps)
	printSteps(&steps)

	sumPart1 := hashSum(&steps)
	fmt.Printf("Part 1 Hash Sum: %d\n", sumPart1)
}

// Parse the input data
func parseInput(input string) []step {
	// Split the input into a slice of strings
	inputSlice := strings.Split(input, ",")
	// Create a slice of steps
	steps := make([]step, len(inputSlice))
	for i, s := range inputSlice {
		steps[i].value = s
	}
	return steps
}

// Calculate the hash for each step
func calculatePartHashs(steps *[]step) {
	for i := range *steps {
		var hash int
		for j := range (*steps)[i].value {
			hash += int((*steps)[i].value[j])
			hash *= 17
			hash &= 0xFF
		}
		(*steps)[i].hash = hash
	}
}

// Calculate the sum of the hash values
func hashSum(steps *[]step) int {
	var sum int
	for i := range *steps {
		sum += (*steps)[i].hash
	}
	return sum
}

// Print the steps
func printSteps(steps *[]step) {
	for i := range *steps {
		fmt.Printf("%6d: (%d) %s \n", i, (*steps)[i].hash, (*steps)[i].value)
	}
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
