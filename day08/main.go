package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// Define a multipath structure
// type multipath struct {
// 	paths [][]string
// }

type multipath struct {
	first   []string
	current []string
	next    []string
	steps   int
}

type looppath struct {
	first   string
	current string
	next    string
	length  int
}

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 8)

	// Get input string
	// inputExample1Str := `RL

	// AAA = (BBB, CCC)
	// BBB = (DDD, EEE)
	// CCC = (ZZZ, GGG)
	// DDD = (DDD, DDD)
	// EEE = (EEE, EEE)
	// GGG = (GGG, GGG)
	// ZZZ = (ZZZ, ZZZ)`

	// inputExample1aStr := `LLR

	// AAA = (BBB, BBB)
	// BBB = (AAA, ZZZ)
	// ZZZ = (ZZZ, ZZZ)`

	// inputExample2Str := `LR

	// 11A = (11B, XXX)
	// 11B = (XXX, 11Z)
	// 11Z = (11B, XXX)
	// 22A = (22B, XXX)
	// 22B = (22C, 22C)
	// 22C = (22Z, 22Z)
	// 22Z = (22B, 22B)
	// XXX = (XXX, XXX)`

	inputStr := loadFileContents("directions.txt")
	// inputStr := inputExample1Str
	// inputStr := inputExample1aStr
	// inputStr := inputExample2Str

	// Load directions
	// directions := parseDirections(inputExample1Str)
	directions := parseDirections(inputStr)
	// spew.Dump(directions)

	// Load rules
	rules := parseRules(inputStr)
	// spew.Dump(rules)

	// Walk the path, part 1
	// path := walkPath(directions, rules)
	// NOTE: number of steps is one less than the number of elements in the path
	// fmt.Printf("Part 1 Path Start: %s  End: %s  Number of Steps: %d\n", path[0], path[len(path)-1], len(path)-1)

	fmt.Println("--------------------")

	// Walk the path, part 2
	// NOTE: Brute force walk not going to work. Trying alternative approach
	// paths2 := walkPath2f(directions, rules)
	// spew.Dump(paths2)
	// fmt.Printf("Part 2 Number of: Paths = %d, Steps = %d\n", len(paths2.first), paths2.steps)

	// Need to find size of loops for each loop
	// Find all the start points that end with A
	startPoints := findStartPoints(rules[0], "A")

	//  Create an array of looppaths based on the found startpoints
	var looppaths []looppath = make([]looppath, len(startPoints))
	for i, _ := range startPoints {
		looppaths[i].first = startPoints[i]
		looppaths[i].current = startPoints[i]
	}

	for i, _ := range looppaths {
		looppaths[i].walkPath2l(directions, rules)
	}
	spew.Dump(looppaths)

	lcm := looppaths[0].length
	if len(looppaths) == 2 {
		lcm = findLCM(looppaths[0].length, looppaths[1].length)
	} else if len(looppaths) > 2 {
		extraValues := make([]int, len(looppaths)-2)
		// Additional values for LCM calculation
		for i := 2; i < len(looppaths); i++ {
			extraValues[i-2] = looppaths[i].length
		}
		lcm = findLCM(looppaths[0].length, looppaths[1].length, extraValues...)
	}
	fmt.Printf("Part 2 Number of: Paths = %d, Steps = %d\n", len(looppaths), lcm)

}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func findLCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = findLCM(result, integers[i])
	}

	return result
}

// Parse directions
func parseDirections(input string) string {
	// Read and return first line of string
	lines := strings.Split(input, "\n")
	return strings.Trim(lines[0], " ")
}

// Parse rules
func parseRules(input string) []map[string]string {
	// Split into lines, then loop over lines, ignoring the first line
	lines := strings.Split(input, "\n")
	rules := make([]map[string]string, 2)
	rules[0] = make(map[string]string)
	rules[1] = make(map[string]string)

	for i := 1; i < len(lines); i++ {
		// If line is empty, skip it
		if len(lines[i]) == 0 {
			continue
		}
		// Split line into source, ldest and rdest
		parts := strings.Split(lines[i], "=")
		source := strings.Trim(parts[0], " \t")
		destinations := strings.Trim(parts[1], " ()")
		ldest := strings.Split(destinations, ",")[0]
		rdest := strings.Split(destinations, ",")[1]
		// Add rules to maps
		rules[0][source] = strings.Trim(ldest, " ")
		rules[1][source] = strings.Trim(rdest, " ")
	}
	return rules
}

// Walk the path
func walkPath(directions string, rules []map[string]string) []string {
	// Start at the root
	current := "AAA"
	var path []string
	path = append(path, current)
	var pathLength int
	directionsLength := len(directions)
	// Loop over directions
	for dirPos := 0; current != "ZZZ" && pathLength < 1000000; dirPos++ {
		// fmt.Printf("Pos %d %d %c\n", dirPos, dirPos%directionsLength, directions[dirPos%directionsLength])
		// fmt.Printf("Current (%d) %s\n", pathLength, current)
		// If direction is left, follow left rule
		if directions[dirPos%directionsLength] == 'L' {
			current = rules[0][current]
		}
		// If direction is right, follow right rule
		if directions[dirPos%directionsLength] == 'R' {
			current = rules[1][current]
		}
		// Check if current is valid
		if current == "" {
			panic("Invalid path")
		}
		// Add current to path
		path = append(path, current)
		pathLength++
	}
	return path
}

func (lp *looppath) walkPath2l(directions string, rules []map[string]string) {
	var finished bool
	var directionsLength = len(directions)
	for dirPos := 0; !finished; dirPos++ {
		if dirPos%10000000 == 0 {
			fmt.Printf(" > Pos %d %d %c\n", dirPos, dirPos%directionsLength, directions[dirPos%directionsLength])
		}
		var activeRules *map[string]string
		switch directions[dirPos%directionsLength] {
		case 'L':
			activeRules = &rules[0]
		case 'R':
			activeRules = &rules[1]
		default:
			panic("Invalid direction")
		}
		// Increase loop length count
		lp.length++
		fmt.Printf("Loop start %s, Count %d\n", lp.first, lp.length)

		// Follow rule
		lp.next = (*activeRules)[lp.current]
		// Check if next is valid
		if lp.next == "" {
			panic("Invalid path")
		}

		// Check if path is finished
		finished = strings.HasSuffix(lp.next, "Z")

		// Copy next to current
		lp.current = lp.next
	}
}

// Next step
func (lp *looppath) nextStepl(rules *map[string]string) bool {
	// Increase loop length count
	lp.length++
	fmt.Printf("Loop start %s, Count %d\n", lp.first, lp.length)

	// Follow rule
	lp.next = (*rules)[lp.current]
	// Check if next is valid
	if lp.next == "" {
		panic("Invalid path")
	}

	// Check if path is finished
	finished := strings.HasSuffix(lp.next, "Z")

	// Copy next to current
	lp.current = lp.next

	return finished
}

// func walkPath2f(directions string, rules []map[string]string) multipath {
// 	const maxPathLength = 10000000
// 	paths := multipath{}

// 	// Find all the start points that end with A
// 	startPoints := findStartPoints(rules[0], "A")
// 	// Remember the start nodes and initialise the current nodes
// 	paths.first = startPoints
// 	paths.current = append(paths.current, startPoints...)
// 	paths.next = make([]string, len(startPoints))

// 	var finished bool
// 	var pathLength int
// 	var directionsLength = len(directions)
// 	for dirPos := 0; !finished && pathLength < maxPathLength; dirPos++ {
// 		if dirPos%10000000 == 0 {
// 			fmt.Printf(" > Pos %d %d %c\n", dirPos, dirPos%directionsLength, directions[dirPos%directionsLength])
// 		}
// 		var activeRules *map[string]string
// 		switch directions[dirPos%directionsLength] {
// 		case 'L':
// 			activeRules = &rules[0]
// 		case 'R':
// 			activeRules = &rules[1]
// 		default:
// 			panic("Invalid direction")
// 		}
// 		finished = nextStepf(&paths, activeRules)
// 	}

// 	if pathLength >= maxPathLength {
// 		panic("Path too long")
// 	}

// 	return paths
// }

// // Next step
// func nextStepf(paths *multipath, rules *map[string]string) bool {
// 	var finished bool
// 	numPaths := len(paths.current)

// 	// Increase step count
// 	paths.steps++
// 	// fmt.Printf("Step %d\n", paths.steps)

// 	// Loop over paths
// 	for pathPos := 0; pathPos < numPaths; pathPos++ {
// 		// Get current node
// 		current := paths.current[pathPos]

// 		// Follow rule
// 		next := (*rules)[current]
// 		// Check if next is valid
// 		if next == "" {
// 			panic("Invalid path")
// 		}
// 		// Add next to path
// 		// fmt.Printf("Path %d: %s -> %s\n", pathPos, current, next)
// 		paths.next[pathPos] = next
// 	}
// 	// Check if all paths are finished
// 	finished = true
// 	for pathPos := 0; pathPos < numPaths; pathPos++ {
// 		if !strings.HasSuffix(paths.next[pathPos], "Z") {
// 			finished = false
// 		}
// 	}
// 	// Copy next to current
// 	for pathPos := 0; pathPos < numPaths; pathPos++ {
// 		paths.current[pathPos] = paths.next[pathPos]
// 	}

// 	return finished
// }

// Walk the path for part 2
// func walkPath2(directions string, rules []map[string]string) multipath {
// 	const maxPathLength = 1000000
// 	paths := multipath{}

// 	// Find all the start points that end with A
// 	startPoints := findStartPoints(rules[0], "A")
// 	paths.paths = make([][]string, len(startPoints))
// 	for i := 0; i < len(startPoints); i++ {
// 		paths.paths[i] = make([]string, 0, 5000)
// 		paths.paths[i] = append(paths.paths[i], startPoints[i])
// 	}

// 	var finished bool
// 	var pathLength int
// 	var directionsLength = len(directions)
// 	for dirPos := 0; !finished && pathLength < maxPathLength; dirPos++ {
// 		if dirPos%1000000 == 0 {
// 			fmt.Printf(" . Pos %d %d %c\n", dirPos, dirPos%directionsLength, directions[dirPos%directionsLength])
// 		}
// 		var activeRules *map[string]string
// 		switch directions[dirPos%directionsLength] {
// 		case 'L':
// 			activeRules = &rules[0]
// 		case 'R':
// 			activeRules = &rules[1]
// 		default:
// 			panic("Invalid direction")
// 		}
// 		finished = nextStep(&paths, activeRules)
// 	}

// 	if pathLength >= maxPathLength {
// 		panic("Path too long")
// 	}

// 	return paths
// }

// // Next step
// func nextStep(paths *multipath, rules *map[string]string) bool {
// 	var finished bool
// 	numPaths := len(paths.paths)
// 	lastIndex := len(paths.paths[0]) - 1

// 	// Loop over paths
// 	for pathPos := 0; pathPos < numPaths; pathPos++ {
// 		// Get current node
// 		current := paths.paths[pathPos][lastIndex]

// 		// Follow rule
// 		next := (*rules)[current]
// 		// Check if next is valid
// 		if next == "" {
// 			panic("Invalid path")
// 		}
// 		// Add next to path
// 		paths.paths[pathPos] = append(paths.paths[pathPos], next)
// 	}
// 	// Check if all paths are finished
// 	finished = true
// 	newLastIndex := len(paths.paths[0]) - 1
// 	for pathPos := 0; pathPos < numPaths; pathPos++ {
// 		if !strings.HasSuffix(paths.paths[pathPos][newLastIndex], "Z") {
// 			finished = false
// 		}
// 	}
// 	return finished
// }

// Find all the start points
func findStartPoints(rules map[string]string, suffix string) []string {
	var startPoints []string
	for key, _ := range rules {
		if strings.HasSuffix(key, suffix) {
			startPoints = append(startPoints, key)
		}
	}
	return startPoints
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
