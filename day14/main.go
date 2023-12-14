package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

const (
	north = iota
	east
	south
	west
)

type mapState struct {
	originalMap []string
	// mapState is a struct to hold the state of the map
	roundRocks [][]bool

	// Array holding the gap details for each direction
	// 0 = north, 1 = east, 2 = south, 3 = west
	gapEnds [4][][]int

	dim int
}

const mirrorMapStr = `O....#....
O.OO#....#
.....##...
OO.#O....O
.O.....O#.
O.#..O.#.#
..O..#O..O
.......O..
#....###..
#OO..#....`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 14)

	// Load the map
	// inputMap := mirrorMapStr
	inputMap := loadFileContents("mirror_map.txt")

	// Parse the map
	mirrorMap := parseMap(inputMap)

	// Print the map
	printMap(mirrorMap)

	// rotatedMap := rotateMap(mirrorMap, -1)

	// fmt.Println()
	// fmt.Println("Rotated map:")
	// printMap(rotatedMap)

	// rotatedMap = rotateMap(rotatedMap, 1)

	// fmt.Println()
	// fmt.Println("Rotated map:")
	// printMap(rotatedMap)

	fmt.Printf("\nRolled North:\n")
	// Roll the rocks north
	rotatedMap := rollRocksNorth(mirrorMap)
	printMap(rotatedMap)

	// Calculate the load of the map
	load := calculateLoad(rotatedMap)
	fmt.Printf("Load: %d\n", load)

	fmt.Println("|-----------------------------------------------------------------------------------------|")

	// Parse the map state
	m := parseMapState(inputMap)
	// Print the map state
	printMapState(m)

	// // Roll the rocks north
	// rollRocksMapState(&m, north)
	// fmt.Println("Rolled North:")
	// printMapState(m)

	// // Roll the rocks west
	// rollRocksMapState(&m, west)
	// fmt.Println("Rolled West:")
	// printMapState(m)

	// // Roll the rocks south
	// rollRocksMapState(&m, south)
	// fmt.Println("Rolled South:")
	// printMapState(m)

	// // Roll the rocks east
	// rollRocksMapState(&m, east)
	// fmt.Println("Rolled East:")
	// printMapState(m)

	// numCycles := 3
	// NOTE: This took >5 hours to run for the example input 😬, need quicker method
	// numCycles := 1000000000
	// // Rotate map to the starting position
	// cycledMap := rotateMap(mirrorMap, -1)
	// fmt.Printf("\nRolled Cycles (%d):\n", numCycles)
	// for i := 0; i < numCycles; i++ {
	// 	if i%1000000 == 0 {
	// 		fmt.Printf(" - Cycle %10d (%d/%d)\n", i, i/1000000, numCycles/1000000)
	// 	}
	// 	// Roll the rockcycle
	// 	cycledMap = rollRockCycle(cycledMap)
	// 	// printMap(cycledMap)
	// }
	// // Rotate map back to original orientation
	// cycledMap = rotateMap(cycledMap, 1)

	// // Calculate the load of the map
	// cycledLoad := calculateLoad(cycledMap)
	// fmt.Printf("Load after %d cycles: %d\n", numCycles, cycledLoad)

	// Start profiling
	f, err := os.Create("myprogram.prof")
	if err != nil {

		fmt.Println(err)
		return

	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// numCycles := 3
	// numCycles := 1000000000
	numCycles := 500000
	fmt.Printf("\nRolled Cycles (%d):\n", numCycles)
	start := time.Now()
	for i := 0; i < numCycles; i++ {
		if i%100000 == 0 {
			fmt.Printf(" - Cycle %10d (%d/%d)", i, i/100000, numCycles/100000)
			if i > 0 {
				timeElapsed := time.Since(start)
				fmt.Printf(" took %s", timeElapsed.Round(time.Second))
			}
			fmt.Println()
			start = time.Now()
		}
		// Roll the rockcycle
		rollRocksMapStateCycle(&m)
		// fmt.Println(strings.Join(mapStateToString(m), "\n"))
		// fmt.Println()
		// printMapState(m)
	}
	fmt.Println(strings.Join(mapStateToString(m), "\n"))
	// printMapState(m)

	// Calculate the load of the map
	cycledLoad := calculateLoad(mapStateToString(m))
	fmt.Printf("Part 2 Load after %d cycles: %d\n", numCycles, cycledLoad)

}

// Roll the rocks north
func rollRocksNorth(mirrorMap []string) []string {
	// Rotate map so the first row is the first column
	rotatedMap := rotateMap(mirrorMap, -1)

	// Roll the rocks north
	for i := 0; i < len(rotatedMap); i++ {
		// Roll the rocks in the row
		rotatedMap[i] = rollRocks(rotatedMap[i])

	}
	// Rotate map back to original orientation
	rotatedMap = rotateMap(rotatedMap, 1)

	return rotatedMap
}

// Roll Rock Cycle
func rollRockCycle(mirrorMap []string) []string {
	// Assume the map is already rotated so that north is facing left
	// rotatedMap := rotateMap(mirrorMap, -1)
	rotatedMap := mirrorMap
	same := true
	// Perform a roll north, west, south, east
	for range []string{"north", "west", "south", "east"} {
		// fmt.Printf("\nRolled %s:\n", direction)
		// Roll the rocks
		for i := 0; i < len(rotatedMap); i++ {
			// Roll the rocks in the row
			rotatedMap[i] = rollRocks(rotatedMap[i])

		}
		// Rotate map clockwise for the next direction
		rotatedMap = rotateMap(rotatedMap, 1)
	}
	for i := range mirrorMap {
		if mirrorMap[i] != rotatedMap[i] {
			same = false
		}
	}
	if same {
		fmt.Println("Same")
	}
	// Rotate map back to original orientation
	// rotatedMap = rotateMap(rotatedMap, 1)

	return rotatedMap
}

// Calculate load of the map
func calculateLoad(mirrorMap []string) int {
	// Calculate the load of the map
	load := 0
	for i, line := range mirrorMap {
		load += strings.Count(line, "O") * (len(mirrorMap) - i)
	}

	return load
}

// Rotate the map
func rotateMap(mirrorMap []string, direction int) []string {
	var resultingMap []string

	// Rotate map anticlockwise
	if direction == -1 {
		// Create a new map, rows will be columns
		resultingMap = make([]string, len(mirrorMap[0]))
		// Rotate map anticlockwise
		// Loop through the columns, starting with the last column, scan from top to bottom creating a new row
		newRowCount := 0
		for i := len(mirrorMap[0]) - 1; i >= 0; i-- {
			for j := 0; j < len(mirrorMap); j++ {
				resultingMap[newRowCount] += string(mirrorMap[j][i])
			}
			newRowCount++
		}

	} else if direction == 1 {
		// Create a new map, rows will be columns
		resultingMap = make([]string, len(mirrorMap[0]))
		// Rotate map clockwise
		// Start with first column, scan from bottom to top creating a new row
		newRowCount := 0
		for i := 0; i < len(mirrorMap[0]); i++ {
			for j := len(mirrorMap) - 1; j >= 0; j-- {
				resultingMap[newRowCount] += string(mirrorMap[j][i])
			}
			newRowCount++
		}
	} else {
		panic("Invalid direction")
	}

	return resultingMap
}

// Roll the rocks in a row
func rollRocks(row string) string {
	var rolledRow string
	// Split lines between fixed rocks
	lines := strings.Split(row, "#")

	// Loop over all lines
	for i := range lines {
		// Count the number of rocks in the line
		rockCount := strings.Count(lines[i], "O")
		// Create new line of the same length as the original line
		// but with all the rocks at the front
		lines[i] = strings.Repeat("O", rockCount) + strings.Repeat(".", len(lines[i])-rockCount)
	}

	rolledRow = strings.Join(lines, "#")
	// fmt.Printf("Rolled row: %s\n", rolledRow)

	return rolledRow
}

// Parse the map
func parseMap(inputMap string) []string {
	// Split the input string into lines
	lines := splitLines(inputMap)

	return lines
}

func parseMapState(inputMap string) mapState {
	// Split the input string into lines
	lines := splitLines(inputMap)
	// Create a new mapState
	var m mapState
	// Set the original map
	m.originalMap = lines

	// Check the map is square
	if len(lines) != len(lines[0]) {
		panic("Map is not square")
	} else {
		m.dim = len(lines)
	}

	// Loop over each row in the map setting the cube Rock locations
	//ToDo: Finish this

	// Calculate the gaps in each direction
	// 0 = north, 1 = east, 2 = south, 3 = west
	// Create a new array and parse square rock locations for each direction
	m.gapEnds[north] = make([][]int, len(lines))
	// Capture square rock location in columns from the north direction
	// Add the start boundary
	for i := range m.gapEnds[north] {
		m.gapEnds[north][i] = append(m.gapEnds[north][i], -1)
	}
	for j := range lines {
		for i := range lines[j] {
			// Check if the current location is a square rock
			if lines[j][i] == '#' {
				m.gapEnds[north][i] = append(m.gapEnds[north][i], j)
			}
		}
	}
	for i := range m.gapEnds[north] {
		m.gapEnds[north][i] = append(m.gapEnds[north][i], len(lines[i]))
	}

	// Capture square rock location in columns from the east direction
	m.gapEnds[east] = make([][]int, len(lines[0]))
	// Add the start boundary
	for i := range m.gapEnds[east] {
		m.gapEnds[east][i] = append(m.gapEnds[east][i], -1)
	}
	for j := range lines {
		for i := len(lines[j]) - 1; i >= 0; i-- {
			// Check if the current location is a square rock
			if lines[j][i] == '#' {
				m.gapEnds[east][j] = append(m.gapEnds[east][j], len(lines[j])-1-i)
			}
		}
	}
	for i := range m.gapEnds[north] {
		m.gapEnds[east][i] = append(m.gapEnds[east][i], len(lines[i]))
	}

	// Capture square rock location in columns from the south direction
	m.gapEnds[south] = make([][]int, len(lines))
	// Add the start boundary
	for i := range m.gapEnds[south] {
		m.gapEnds[south][i] = append(m.gapEnds[south][i], -1)
	}
	for j := len(lines) - 1; j >= 0; j-- {
		for i := len(lines[j]) - 1; i >= 0; i-- {
			// Check if the current location is a square rock
			if lines[j][i] == '#' {
				m.gapEnds[south][len(lines[j])-1-i] = append(m.gapEnds[south][len(lines[j])-1-i], len(lines)-1-j)
			}
		}
	}
	for i := range m.gapEnds[south] {
		m.gapEnds[south][i] = append(m.gapEnds[south][i], len(lines[i]))
	}

	// Capture square rock location in columns from the west direction
	m.gapEnds[west] = make([][]int, len(lines[0]))
	// Add the start boundary
	for i := range m.gapEnds[west] {
		m.gapEnds[west][i] = append(m.gapEnds[west][i], -1)
	}
	for j := len(lines) - 1; j >= 0; j-- {
		for i := range lines[j] {
			// Check if the current location is a square rock
			if lines[j][i] == '#' {
				m.gapEnds[west][len(lines)-1-j] = append(m.gapEnds[west][len(lines)-1-j], i)
			}
		}
	}
	for i := range m.gapEnds[west] {
		m.gapEnds[west][i] = append(m.gapEnds[west][i], len(lines[i]))
	}

	// Set the roundRocks
	m.roundRocks = make([][]bool, len(lines))
	// Loop over each row in the map setting the round Rock locations
	for i := range lines {
		m.roundRocks[i] = make([]bool, len(lines[i]))
		for j := range lines[i] {
			if string(lines[i][j]) == "O" {
				m.roundRocks[i][j] = true
			}
		}
	}

	return m
}

func rollRocksMapStateCycle(m *mapState) {
	// Roll the rocks north
	rollRocksMapState(m, north)
	// Roll the rocks west
	rollRocksMapState(m, west)
	// Roll the rocks south
	rollRocksMapState(m, south)
	// Roll the rocks east
	rollRocksMapState(m, east)
}

func rollRocksMapState(m *mapState, dir int) {
	// Depending on the direction, roll the rocks in the map in that direction
	// 0 = north, 1 = east, 2 = south, 3 = west
	switch dir {
	case north:
		// Roll the rocks north, loop over each column
		for i := 0; i < len(m.gapEnds[north]); i++ {
			// fmt.Printf("Column %d\n", i)
			// Loop over each gap between rocks in the column, counting rolling rocks
			for e := 0; e < len(m.gapEnds[north][i])-1; e++ {
				rollingRocks := 0
				// Count the number of rocks in the gap and remove rolling rocks
				// fmt.Printf("Gap %d: %d to %d\n", e, m.gapEnds[north][i][e], m.gapEnds[north][i][e+1])
				for j := m.gapEnds[north][i][e] + 1; j < m.gapEnds[north][i][e+1]; j++ {
					if m.roundRocks[j][i] {
						rollingRocks++
						m.roundRocks[j][i] = false
					}
				}
				// fmt.Printf("  Rolling rocks: %d\n", rollingRocks)
				// Add any rolling rocks to the top of the gap
				for j := m.gapEnds[north][i][e] + 1; j < m.gapEnds[north][i][e]+rollingRocks+1; j++ {
					m.roundRocks[j][i] = true
				}
			}

		}
	case east:
		// Roll the rocks east
		for i := 0; i < len(m.gapEnds[east]); i++ {
			// fmt.Printf("Column %d\n", i)
			// Loop over each gap between rocks in the column, counting rolling rocks
			for e := 0; e < len(m.gapEnds[east][i])-1; e++ {
				rollingRocks := 0
				// Count the number of rocks in the gap and remove rolling rocks
				// fmt.Printf("Gap %d: %d to %d\n", e, m.gapEnds[east][i][e], m.gapEnds[east][i][e+1])
				for j := m.gapEnds[east][i][e] + 1; j < m.gapEnds[east][i][e+1]; j++ {
					if m.roundRocks[i][m.dim-1-j] {
						rollingRocks++
						m.roundRocks[i][m.dim-1-j] = false
					}
				}
				// Add any rolling rocks to the top of the gap
				for j := m.gapEnds[east][i][e] + 1; j < m.gapEnds[east][i][e]+rollingRocks+1; j++ {
					m.roundRocks[i][m.dim-1-j] = true
				}
			}
		}
	case south:
		// Roll the rocks south
		for i := 0; i < len(m.gapEnds[south]); i++ {
			// fmt.Printf("Column %d\n", i)
			// Loop over each gap between rocks in the column, counting rolling rocks
			for e := 0; e < len(m.gapEnds[south][i])-1; e++ {
				rollingRocks := 0
				// Count the number of rocks in the gap and remove rolling rocks
				// fmt.Printf("Gap %d: %d to %d\n", e, m.gapEnds[south][i][e], m.gapEnds[south][i][e+1])
				for j := m.gapEnds[south][i][e] + 1; j < m.gapEnds[south][i][e+1]; j++ {
					if m.roundRocks[m.dim-1-j][m.dim-1-i] {
						rollingRocks++
						m.roundRocks[m.dim-1-j][m.dim-1-i] = false
					}
				}
				// Add any rolling rocks to the top of the gap
				for j := m.gapEnds[south][i][e] + 1; j < m.gapEnds[south][i][e]+rollingRocks+1; j++ {
					m.roundRocks[m.dim-1-j][m.dim-1-i] = true
				}
			}
		}

	case west:
		// Roll the rocks west
		for i := 0; i < len(m.gapEnds[west]); i++ {
			// fmt.Printf("Column %d\n", i)
			// Loop over each gap between rocks in the column, counting rolling rocks
			for e := 0; e < len(m.gapEnds[west][i])-1; e++ {
				rollingRocks := 0
				// Count the number of rocks in the gap and remove rolling rocks
				// fmt.Printf("Gap %d: %d to %d\n", e, m.gapEnds[west][i][e], m.gapEnds[west][i][e+1])
				for j := m.gapEnds[west][i][e] + 1; j < m.gapEnds[west][i][e+1]; j++ {
					if m.roundRocks[m.dim-1-i][j] {
						rollingRocks++
						m.roundRocks[m.dim-1-i][j] = false
					}
				}
				// Add any rolling rocks to the top of the gap
				for j := m.gapEnds[west][i][e] + 1; j < m.gapEnds[west][i][e]+rollingRocks+1; j++ {
					m.roundRocks[m.dim-1-i][j] = true
				}
			}
		}
	default:
		panic("Invalid direction")
	}
}

// Print the map state
func printMapState(m mapState) {
	// Print the map
	for _, line := range m.originalMap {
		fmt.Printf("%s\n", line)
	}

	// Print the gaps
	fmt.Println("Gaps:")
	for i := range m.gapEnds {
		fmt.Printf("  %d: %v\n", i, m.gapEnds[i])
	}

	// Print the roundRock locations
	fmt.Println("Round Rocks:")
	// for i := range m.roundRocks {
	// 	fmt.Printf("  %d: %v\n", i, m.roundRocks[i])
	// }
	rockMap := mapStateToString(m)
	for _, line := range rockMap {
		fmt.Printf("%s\n", line)
	}
}

// Print the map state
func mapStateToString(m mapState) []string {
	var resultingMap []string
	// Loop over all rows
	for i := range m.originalMap {
		// Create a stringbuilder to hold the row
		sb := strings.Builder{}
		sb.Grow(len(m.originalMap[i]))
		// Loop over all columns
		for j := range m.originalMap[i] {
			// Check if the current location is a square rock
			if m.roundRocks[i][j] {
				sb.WriteString("O")
			} else if m.originalMap[i][j] == '#' {
				sb.WriteString("#")
			} else {
				sb.WriteString(".")
			}
		}
		resultingMap = append(resultingMap, sb.String())
	}
	return resultingMap
}

// Split the lines in a string array
func splitLines(input string) []string {
	// Split the input string into lines
	lines := []string{}
	for _, line := range strings.Split(input, "\n") {
		if len(line) > 0 {
			lines = append(lines, string(line))
		}
	}

	return lines
}

// Print the map
func printMap(mirrorMap []string) {
	// Print the map
	for _, line := range mirrorMap {
		fmt.Printf("%s\n", line)
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
