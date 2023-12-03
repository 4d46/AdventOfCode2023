package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// Coordinate Structure
type Coordinate struct {
	x int
	y int
}

// Part Structure
type Part struct {
	name     rune
	position Coordinate
}

// Label Structure
type Label struct {
	name  string
	start Coordinate
	end   Coordinate
}

// Schematic Structure
type Schematic struct {
	labels []Label
	parts  []Part
}

// Example Schematic
var exampleSchematicPart1 = `467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..`

var exampleSchematicPart2 = `467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 03)

	// Load engine schematic file
	schematicString := loadFileContents("engineschematic.txt")
	// schematicString := exampleSchematicPart1
	// schematicString := exampleSchematicPart2

	// Parse schematic string into a Schematic structure
	schematic := parseSchematic(schematicString)

	// spew.Dump(schematic)
	// Print number of parts and number of labels
	fmt.Printf("Number of parts: %d\n", len(schematic.parts))
	fmt.Printf("Number of labels: %d\n", len(schematic.labels))

	// Sum the labels that are part numbers in a schematic
	sum := sumPartNumbers(schematic)
	fmt.Printf("Sum of part numbers: %d\n", sum)

	// printSchematic(schematic)
	gearRatioSum := findGearRatioSum(schematic)
	fmt.Printf("Sum of gear ratios: %d\n", gearRatioSum)
}

// Function that sums the labels that are part numbers in a schematic
func sumPartNumbers(schematic Schematic) int {
	// Create a sum variable
	sum := 0

	// Loop through labels
	for _, label := range schematic.labels {
		// Check if label is a part number, by finding if it adjactent to a part
		// part, err := findAdjactentPart(schematic, label)
		_, err := findAdjactentPart(schematic, label)
		if err == nil {
			// Add part number to sum
			labelValue, err := strconv.Atoi(label.name)
			if err != nil {
				panic(err)
			}
			sum += labelValue
			// Print details of matching label
			// fmt.Printf("Label %s [%d,%d] is a part number for %c\n", label.name, label.start.x, label.start.y, part.name)
		}
		// } else {
		// 	fmt.Printf("Label %s [%d,%d] is not a part number\n", label.name, label.start.x, label.start.y)
		// }
	}

	// Return sum
	return sum
}

// Function to find the sum of the gear ratios
func findGearRatioSum(schematic Schematic) int {
	// Create a sum variable
	sum := 0

	// Loop over parts
	for _, part := range schematic.parts {
		// Check if part is a gear symbol
		if part.name == '*' {
			// Find labels adjactent to part
			labels := findAdjactentLabels(schematic, part)
			// Confirm there are only 2 labels
			if len(labels) == 2 {
				// Find the gear ratio
				firstGear, err := strconv.Atoi(labels[0].name)
				if err != nil {
					panic(err)
				}
				secondGear, err := strconv.Atoi(labels[1].name)
				if err != nil {
					panic(err)
				}
				gearRatio := firstGear * secondGear

				// Add gear ratio to sum
				sum += gearRatio
			}
		}
	}
	return sum
}

// Function that finds labels adjactent to a part
func findAdjactentLabels(schematic Schematic, part Part) []Label {
	// Create a slice of labels
	labels := []Label{}

	// Loop through labels
	for _, label := range schematic.labels {
		// Check if label is adjactent to part
		if isPartAdjactentToLabel(part, label) {
			// Add label to slice
			labels = append(labels, label)
		}
	}

	// Return labels
	return labels
}

// Function that finds a part adjactent to a label
func findAdjactentPart(schematic Schematic, label Label) (Part, error) {
	// Loop through parts
	for _, part := range schematic.parts {
		// Check if part is adjactent to label
		if isPartAdjactentToLabel(part, label) {
			// Return part number
			return part, nil
		}
	}

	// Return error
	return Part{}, fmt.Errorf("no part adjactent to label")
}

// Function that checks if a part is adjactent to a label
func isPartAdjactentToLabel(part Part, label Label) bool {
	adjactentX := false
	adjactentY := false

	// if label.name == "515" {
	// 	// Print out details of part and label
	// 	fmt.Printf("Part %c [%d,%d] is adjactent to label %s [%d,%d] [%d,%d]\n", part.name, part.position.x, part.position.y, label.name, label.start.x, label.start.y, label.end.x, label.end.y)
	// 	fmt.Printf("label.end.x >= part.position.x-1 %t\n", label.end.x >= part.position.x-1)
	// 	fmt.Printf("label.start.x <= part.position.x+1 %t\n", label.start.x <= part.position.x+1)
	// 	fmt.Printf("label.end.y >= part.position.y-1 %t\n", label.end.y >= part.position.y-1)
	// 	fmt.Printf("label.start.y <= part.position.y+1 %t\n", label.start.y <= part.position.y+1)
	// }

	// Check if part is adjactent to label
	if label.end.x >= part.position.x-1 && label.start.x <= part.position.x+1 {
		adjactentX = true
	}
	if label.end.y >= part.position.y-1 && label.start.y <= part.position.y+1 {
		adjactentY = true
	}
	distanceStart := math.Sqrt(math.Pow(float64(label.start.x-part.position.x), 2) + math.Pow(float64(label.start.y-part.position.y), 2))
	distanceEnd := math.Sqrt(math.Pow(float64(label.end.x-part.position.x), 2) + math.Pow(float64(label.end.y-part.position.y), 2))
	if distanceStart > 1.5 && distanceEnd > 1.5 && adjactentX && adjactentY {
		fmt.Printf("Part %c [%d,%d] is not adjactent to label %s [%d,%d] [%d,%d], distances %f %f\n", part.name, part.position.x, part.position.y, label.name, label.start.x, label.start.y, label.end.x, label.end.y, distanceStart, distanceEnd)
	}

	// Return false
	return adjactentX && adjactentY
}

// Function that prints a schematic
func printSchematic(schematic Schematic) {
	// Create an array of strings each containing 140 '.' characters
	schematicArray := make([]string, 140)
	for i := range schematicArray {
		schematicArray[i] = strings.Repeat(".", 140)
	}
	spew.Dump(schematicArray)

	// Loop through parts
	for _, part := range schematic.parts {
		// Add part to map string
		schematicArray[part.position.y] = schematicArray[part.position.y][:part.position.x] + string(part.name) + schematicArray[part.position.y][part.position.x+1:]
	}

	// Loop through labels
	for _, label := range schematic.labels {
		// Add label to map string
		schematicArray[label.start.y] = schematicArray[label.start.y][:label.start.x] + label.name + schematicArray[label.start.y][label.end.x:]
	}

	// Print schematic
	for _, line := range schematicArray {
		fmt.Println(line)
	}
}

// Parse schematic string into a Schematic structure
func parseSchematic(schematicString string) Schematic {
	// Create a new Schematic structure
	schematic := Schematic{}

	// Split schematic string into lines
	lines := strings.Split(schematicString, "\n")

	// Loop through lines
	for y, line := range lines {
		// Parse line for parts
		parts := parseLineForParts(line, y)
		// Add parts to schematic
		schematic.parts = append(schematic.parts, parts...)
		// Parse line for labels
		labels := parseLineForLabels(line, y)
		// Add labels to schematic
		schematic.labels = append(schematic.labels, labels...)
	}

	// Return schematic
	return schematic
}

// Parse a line for parts
func parseLineForParts(line string, y int) []Part {
	// Create a new slice of parts
	parts := []Part{}

	partSymbols := []rune{'-', '@', '*', '/', '&', '#', '%', '+', '=', '$'}

	// Loop through string looking for symbols, remebering x position
	for x, symbol := range line {
		// Check if symbol is a part
		if runeInList(symbol, partSymbols) {
			// Create a new part
			part := Part{name: symbol, position: Coordinate{x, y}}

			// Add part to slice
			parts = append(parts, part)
		}
	}

	// Return parts
	return parts
}

// Parse a line for labels
func parseLineForLabels(line string, y int) []Label {
	// Create a new slice of labels
	labels := []Label{}

	re := regexp.MustCompile(`\.?([0-9]{1,3})\.?`)
	matches := re.FindAllStringIndex(line, -1)

	// Convert matches to labels
	for _, match := range matches {
		// Get label name
		name := line[match[0]:match[1]]
		startX := match[0]
		endX := match[1]

		// If name starts with a . increase start coord by 1
		if name[0] == '.' {
			startX++
		}

		// If name ends with a . decrease end coord by 1
		if name[len(name)-1] == '.' {
			endX--
		}

		name = strings.Trim(name, ".")

		// fmt.Printf("name: %s, startX: %d, endX: %d\n", name, startX, endX)
		// Validate name
		if name != line[startX:endX] {
			err := fmt.Errorf("name does not match '%s' != '%s' [%d,%d]\n%s", name, line[startX:endX], startX, endX, line)
			panic(err)
		}

		// Create a new label
		label := Label{name: name,
			start: Coordinate{startX, y},
			end:   Coordinate{endX - 1, y},
		}

		// Add label to slice
		labels = append(labels, label)

	}

	// Return labels
	return labels
}

// Check if rune is in a list of runes
func runeInList(r rune, list []rune) bool {
	// Loop through list
	for _, l := range list {
		// Check if rune is in list
		if r == l {
			return true
		}
	}

	return false
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
