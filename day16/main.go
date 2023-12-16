package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	north = 0x1
	east  = 0x2
	south = 0x4
	west  = 0x8
)

type Grid struct {
	layout     [][]rune
	light      [][]int
	maxX, maxY int
}

const layoutStr1 = `.|...\....
|.-.\.....
.....|-...
........|.
..........
.........\
..../.\\..
.-.-/..|..
.|....-|.\
..//.|....`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 16)

	// input := layoutStr1
	input := loadFileContents("layout.txt")

	// Parse layout
	grid := parseLayout(input)

	// Print layout
	printLayout(grid)

	// Track beam
	trackBeam(&grid, 0, 0, east)

	// Print energised grid
	printEnergisedGrid(grid)

	// Count energised cells for Part 1
	count := countEnergisedCells(grid)
	fmt.Printf("Energised cells Part 1: %d\n", count)

	// Part 2
	// Remember most energised start point
	mostEnergised := 0
	mostEnergisedX, mostEnergisedY, mostEnergisedDir := 0, 0, 0

	// Loop through all cells from the north, going south
	for x := range grid.light[0] {
		// Reset energised grid
		resetGrid(&grid)

		// Track beam
		trackBeam(&grid, x, 0, south)

		// Count energised cells
		count := countEnergisedCells(grid)

		// Check if this is the most energised
		if count > mostEnergised {
			mostEnergised = count
			mostEnergisedX, mostEnergisedY, mostEnergisedDir = x, 0, south
		}
	}

	// Loop through all cells from the east, going west
	for y := range grid.light {
		// Reset energised grid
		resetGrid(&grid)

		// Track beam
		trackBeam(&grid, len(grid.light[0])-1, y, west)

		// Count energised cells
		count := countEnergisedCells(grid)

		// Check if this is the most energised
		if count > mostEnergised {
			mostEnergised = count
			mostEnergisedX, mostEnergisedY, mostEnergisedDir = 0, y, west
		}
	}

	// Loop through all cells from the south, going north
	for x := range grid.light[0] {
		// Reset energised grid
		resetGrid(&grid)

		// Track beam
		trackBeam(&grid, x, len(grid.light)-1, north)

		// Count energised cells
		count := countEnergisedCells(grid)

		// Check if this is the most energised
		if count > mostEnergised {
			mostEnergised = count
			mostEnergisedX, mostEnergisedY, mostEnergisedDir = x, len(grid.light)-1, north
		}
	}

	// Loop through all cells from the west, going east
	for y := range grid.light {
		// Reset energised grid
		resetGrid(&grid)

		// Track beam
		trackBeam(&grid, 0, y, east)

		// Count energised cells
		count := countEnergisedCells(grid)

		// Check if this is the most energised
		if count > mostEnergised {
			mostEnergised = count
			mostEnergisedX, mostEnergisedY, mostEnergisedDir = len(grid.light[0])-1, y, east
		}
	}

	// Print most energised
	fmt.Printf("Part 2 Most energised: [%d, %d] %d Count: %d\n", mostEnergisedX, mostEnergisedY, mostEnergisedDir, mostEnergised)
}

// Track beam
func trackBeam(grid *Grid, posX, posY, dir int) {
	hitEnd := false
	for !hitEnd {
		// Move to next point
		posX, posY, dir, hitEnd = traceBeam(grid, posX, posY, dir)
	}
}

// Trace beam
func traceBeam(grid *Grid, posX, posY, dir int) (int, int, int, bool) {
	// Look at current cell contents to help determine next action
	cell := grid.layout[posY][posX]

	// Record the incoming beam direction
	// record the light beam
	grid.light[posY][posX] |= dir

	switch cell {
	case '|':
		// Switch based on light direction
		switch dir {
		case north:
			// Light continues north

			// Calculate next position and direction
			nextX, nextY, nextDir := posX, posY-1, north

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&north > 0
			return nextX, nextY, nextDir, previousBeam
		case east, west:
			// Beam split into two directions, north and south

			// Calculate next position and direction
			nextXN, nextYN, nextDirN := posX, posY-1, north
			nextXS, nextYS, nextDirS := posX, posY+1, south

			// Check if this is retracing a previous beam
			previousBeamN := outGrid(grid, nextXN, nextYN) || grid.light[nextYN][nextXN]&north > 0
			previousBeamS := outGrid(grid, nextXS, nextYS) || grid.light[nextYS][nextXS]&south > 0

			if previousBeamN && !previousBeamS {
				// North beam already traced, go south
				return nextXS, nextYS, nextDirS, previousBeamS
			} else if !previousBeamN && previousBeamS {
				// South beam already traced, go north
				return nextXN, nextYN, nextDirN, previousBeamN
			} else if previousBeamN && previousBeamS {
				// Both beams already traced, return that previous beams found
				return nextDirN, nextDirN, north, true
			} else {
				// Neither beam traced, call function to start tracing one direction and then
				// return the other to continue that direction
				hitEnd := false
				for !hitEnd {
					// Move to next point
					nextXS, nextYS, nextDirS, hitEnd = traceBeam(grid, nextXS, nextYS, nextDirS)
				}
				// Now continue the other direction
				return nextXN, nextYN, nextDirN, previousBeamN
			}
		case south:
			// Light continues south
			// record the light beam
			grid.light[posY][posX] |= south

			// Calculate next position and direction
			nextX, nextY, nextDir := posX, posY+1, south

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&south > 0

			return nextX, nextY, nextDir, previousBeam
		default:
			panic("Invalid direction")
		}
	case '-':
		// Switch based on light direction
		switch dir {
		case north, south:
			// Beam split into two directions, east and west

			// Calculate next position and direction
			nextXE, nextYE, nextDirE := posX+1, posY, east
			nextXW, nextYW, nextDirW := posX-1, posY, west

			// Check if this is retracing a previous beam
			previousBeamE := outGrid(grid, nextXE, nextYE) || grid.light[nextYE][nextXE]&east > 0
			previousBeamW := outGrid(grid, nextXW, nextYW) || grid.light[nextYW][nextXW]&west > 0

			if previousBeamE && !previousBeamW {
				// East beam already traced, go west
				return nextXW, nextYW, nextDirW, previousBeamW
			} else if !previousBeamE && previousBeamW {
				// West beam already traced, go east
				return nextXE, nextYE, nextDirE, previousBeamE
			} else if previousBeamE && previousBeamW {
				// Both beams already traced, return that previous beams found
				return nextDirE, nextDirE, east, true
			} else {
				// Neither beam traced, call function to start tracing one direction and then
				// return the other to continue that direction
				hitEnd := false
				for !hitEnd {
					// Move to next point
					nextXE, nextYE, nextDirE, hitEnd = traceBeam(grid, nextXE, nextYE, nextDirE)
				}
				// Now continue the other direction
				return nextXW, nextYW, nextDirW, previousBeamW
			}
		case east:
			// Light continues east

			// Calculate next position and direction
			nextX, nextY, nextDir := posX+1, posY, east

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&east > 0

			return nextX, nextY, nextDir, previousBeam
		case west:
			// Light continues west

			// Calculate next position and direction
			nextX, nextY, nextDir := posX-1, posY, west

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&west > 0

			return nextX, nextY, nextDir, previousBeam
		default:
			panic("Invalid direction")
		}
	case '/':
		// Switch based on light direction
		switch dir {
		case north:
			// Light continues east

			// Calculate next position and direction
			nextX, nextY, nextDir := posX+1, posY, east

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&east > 0

			return nextX, nextY, nextDir, previousBeam
		case east:
			// Light continues north

			// Calculate next position and direction
			nextX, nextY, nextDir := posX, posY-1, north

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&north > 0

			return nextX, nextY, nextDir, previousBeam
		case south:
			// Light continues west

			// Calculate next position and direction
			nextX, nextY, nextDir := posX-1, posY, west

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&west > 0

			return nextX, nextY, nextDir, previousBeam
		case west:
			// Light continues south

			// Calculate next position and direction
			nextX, nextY, nextDir := posX, posY+1, south

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&south > 0

			return nextX, nextY, nextDir, previousBeam
		default:
			panic("Invalid direction")
		}
	case '\\':
		// Switch based on light direction
		switch dir {
		case north:
			// Light continues west

			// Calculate next position and direction
			nextX, nextY, nextDir := posX-1, posY, west

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&west > 0

			return nextX, nextY, nextDir, previousBeam
		case east:
			// Light continues south

			// Calculate next position and direction
			nextX, nextY, nextDir := posX, posY+1, south

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&south > 0

			return nextX, nextY, nextDir, previousBeam
		case south:
			// Light continues east

			// Calculate next position and direction
			nextX, nextY, nextDir := posX+1, posY, east

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&east > 0

			return nextX, nextY, nextDir, previousBeam
		case west:
			// Light continues north

			// Calculate next position and direction
			nextX, nextY, nextDir := posX, posY-1, north

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&north > 0

			return nextX, nextY, nextDir, previousBeam
		default:
			panic("Invalid direction")
		}
	case '.':
		// Switch based on light direction
		switch dir {
		case north:
			// Light continues north

			// Calculate next position and direction
			nextX, nextY, nextDir := posX, posY-1, north

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&north > 0

			return nextX, nextY, nextDir, previousBeam
		case east:
			// Light continues east

			// Calculate next position and direction
			nextX, nextY, nextDir := posX+1, posY, east

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&east > 0

			return nextX, nextY, nextDir, previousBeam
		case south:
			// Light continues south

			// Calculate next position and direction
			nextX, nextY, nextDir := posX, posY+1, south

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&south > 0

			return nextX, nextY, nextDir, previousBeam
		case west:
			// Light continues west

			// Calculate next position and direction
			nextX, nextY, nextDir := posX-1, posY, west

			// Check if this is retracing a previous beam
			previousBeam := outGrid(grid, nextX, nextY) || grid.light[nextY][nextX]&west > 0

			return nextX, nextY, nextDir, previousBeam
		default:
			panic("Invalid direction")
		}
	default:
		panic("Invalid cell")
	}
}

// Function that checks if a point is outside the grid
func outGrid(grid *Grid, x, y int) bool {
	return x < 0 || y < 0 || x >= grid.maxX || y >= grid.maxY
}

// Count energised cells
func countEnergisedCells(grid Grid) int {
	count := 0
	for _, row := range grid.light {
		for _, col := range row {
			if col > 0 {
				count++
			}
		}
	}

	return count
}

// Parse layout string into a grid
func parseLayout(layoutStr string) Grid {
	// Split string into lines
	lines := splitLines(layoutStr)

	// Create grid
	grid := Grid{
		layout: make([][]rune, len(lines)),
		light:  make([][]int, len(lines)),
		maxX:   len(lines[0]),
		maxY:   len(lines),
	}

	// Parse lines
	for y, line := range lines {
		// Create row
		grid.layout[y] = append(grid.layout[y], []rune(line)...)
		grid.light[y] = make([]int, len(line))
	}

	return grid
}

// Reset grid
func resetGrid(grid *Grid) {
	// Reset grid
	for y := range grid.light {
		for x := range grid.light[y] {
			grid.light[y][x] = 0
		}
	}
}

// Split string into lines
func splitLines(str string) []string {
	// Split string into lines
	lines := []string{}
	for _, line := range strings.Split(str, "\n") {
		lines = append(lines, line)
	}

	return lines
}

// Print layout
func printLayout(grid Grid) {
	// Print layout
	for _, row := range grid.layout {
		for _, col := range row {
			fmt.Printf("%c", col)
		}
		fmt.Println()
	}
}

// Print energised grid
func printEnergisedGrid(grid Grid) {
	// Print layout
	for _, row := range grid.light {
		for _, col := range row {
			if col > 0 {
				fmt.Printf("%c", '#')
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
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
