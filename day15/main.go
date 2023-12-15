package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	gom "github.com/wk8/go-ordered-map"
)

type instruction struct {
	value     string
	hash      int
	label     string
	operator  rune
	operand   int
	labelHash int
}

type box struct {
	contents gom.OrderedMap
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
	printInstructions(&steps)

	sumPart1 := hashSum(&steps)
	fmt.Printf("Part 1 Hash Sum: %d\n", sumPart1)

	// Part 2

	// Define light boxes
	var lightBoxes [256]box
	// Initialise light boxes
	for i := range lightBoxes {
		lightBoxes[i].contents = *gom.New()
	}

	// Calculate the label hash for each step
	calculateLabelHashes(&steps)

	// Process Instructions
	processInstructions(&lightBoxes, &steps)

	printLightBoxes(&lightBoxes)

	// Calculate Focussing Power
	focussingPower := calculateFocussingPower(&lightBoxes)
	fmt.Printf("Focussing Power: %d\n", focussingPower)

}

// Print Lightboxes
func printLightBoxes(lightBoxes *[256]box) {
	for i := range lightBoxes {
		fmt.Printf("%3d: %s\n", i, lightBoxString(lightBoxes[i].contents))
	}
}

// Lightbox String
func lightBoxString(contents gom.OrderedMap) string {
	var s string
	s += fmt.Sprintf("[")
	for pair := contents.Oldest(); pair != nil; pair = pair.Next() {
		s += fmt.Sprintf("%s[%d] ", pair.Key, pair.Value)
	}
	s += fmt.Sprintf("]")
	return s
}

// Parse the input data
func parseInput(input string) []instruction {
	// Split the input into a slice of strings
	inputSlice := strings.Split(input, ",")
	// Create a slice of steps
	steps := make([]instruction, len(inputSlice))
	for i, s := range inputSlice {
		steps[i].value = s
	}
	return steps
}

// Process Instructions
func processInstructions(lightBoxes *[256]box, steps *[]instruction) {
	// Loop over instructions
	for i := range *steps {
		box := &lightBoxes[(*steps)[i].labelHash]

		switch (*steps)[i].operator {
		case '=':
			// Add lens to box
			box.contents.Set((*steps)[i].label, (*steps)[i].operand)
		case '-':
			// Remove lens from box
			box.contents.Delete((*steps)[i].label)
		default:
			panic("Unknown operator")
		}
	}
}

// Calculate Focussing Power
func calculateFocussingPower(lightBoxes *[256]box) int {
	total := 0
	// Loop over lightboxes
	for i := range lightBoxes {
		// Loop over lenses in lightbox
		lensPoistion := 1
		for pair := lightBoxes[i].contents.Oldest(); pair != nil; pair = pair.Next() {
			total += (i + 1) * lensPoistion * pair.Value.(int)
			lensPoistion++
		}
	}

	return total
}

// Calculate the hash for each step
func calculatePartHashs(steps *[]instruction) {
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

// Calculate the label hash for each step
func calculateLabelHashes(steps *[]instruction) {
	for i := range *steps {
		decodeInstruction(&(*steps)[i])
		fmt.Printf("Label: %s %c %d\n", (*steps)[i].label, (*steps)[i].operator, (*steps)[i].operand)
		var hash int
		for _, c := range (*steps)[i].label {
			hash += int(c)
			hash *= 17
			hash &= 0xFF
		}
		(*steps)[i].labelHash = hash
	}
}

// Decode the instruction and set the label, operator and operand
func decodeInstruction(instruction *instruction) {
	// Check if instruction is remove instruction
	if strings.HasSuffix(instruction.value, "-") {
		instruction.operator = '-'
		instruction.label = strings.TrimSuffix(instruction.value, "-")
	} else {
		instruction.operator = '='
		instructionParts := strings.Split(instruction.value, "=")

		instruction.label = instructionParts[0]
		instructionOperand, err := strconv.Atoi(instructionParts[1])
		if err != nil {
			panic(err)
		}
		instruction.operand = instructionOperand
	}
}

// Calculate the sum of the hash values
func hashSum(steps *[]instruction) int {
	var sum int
	for i := range *steps {
		sum += (*steps)[i].hash
	}
	return sum
}

// Print the steps
func printInstructions(steps *[]instruction) {
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
