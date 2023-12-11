package main

import (
	"fmt"
	"os"
	"strings"
)

type Coordinate struct {
	x int
	y int
}

const GalaxyMapExample1 = `...#......
.......#..
#.........
..........
......#...
.#........
.........#
..........
.......#..
#...#.....`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 11)

	// Load the input data
	galaxyMapStr := GalaxyMapExample1
	//galaxyMapStr:=loadFileContents("input.txt")

	// Parse the input data
	galaxyMap := parseGalaxyMap(galaxyMapStr)
	// Print the galaxy map
	fmt.Println(galaxyMap)

	// Expand the Galaxy Map
	blankRows := findEmptyRows(galaxyMap)
	blankColumns := findEmptyColumns(galaxyMap)
	// Expand the galaxy map
	expandedGalaxyMap := expandGalaxyMap(galaxyMap, blankRows, blankColumns)
	// Print the expanded galaxy map
	// fmt.Println(expandedGalaxyMap)
	// Print the expanded galaxy map
	printGalaxyMap(expandedGalaxyMap)
}

// Expand the galaxy map
func expandGalaxyMap(galaxyMap []Coordinate, blankRows []int, blankColumns []int) []Coordinate {
	// Create an empty array of coordinates
	expandedGalaxyMap := []Coordinate{}

	// Loop over the galaxy map
	for _, coordinate := range galaxyMap {
		x := coordinate.x
		y := coordinate.y

		// Loop over the blank rows
		for _, blankRow := range blankRows {
			// Increment the y value if the blank row is above the coordinate
			if blankRow < y {
				y++
			}
		}

		// Loop over the blank columns
		for _, blankColumn := range blankColumns {
			// Increment the x value if the blank column is to the left of the coordinate
			if blankColumn < x {
				x++
			}
		}

		// Add the coordinate to the expanded galaxy map
		expandedGalaxyMap = append(expandedGalaxyMap, Coordinate{x, y})
	}

	return expandedGalaxyMap
}

// Find empty rows
func findEmptyRows(galaxyMap []Coordinate) []int {
	// Create an empty array of rows
	emptyRowMap := map[int]int{}

	for y := 0; y < len(galaxyMap); y++ {
		emptyRowMap[y] = y
	}

	// Loop over the galaxy map, remove from the array if the row is not empty
	for _, coordinate := range galaxyMap {
		// Remove the row index from the map
		delete(emptyRowMap, coordinate.y)
	}

	// Create an empty array of rows
	emptyRows := []int{}

	// Copy row from map into array
	for _, row := range emptyRowMap {
		emptyRows = append(emptyRows, row)
	}

	return emptyRows
}

// Find empty columns
func findEmptyColumns(galaxyMap []Coordinate) []int {
	// Create an empty array of columns
	emptyColumnMap := map[int]int{}

	for x := 0; x < len(galaxyMap); x++ {
		emptyColumnMap[x] = x
	}

	// Loop over the galaxy map, remove from the array if the column is not empty
	for _, coordinate := range galaxyMap {
		// Remove the column index from the map
		delete(emptyColumnMap, coordinate.x)
	}

	// Create an empty array of columns
	emptyColumns := []int{}

	// Copy column from map into array
	for _, column := range emptyColumnMap {
		emptyColumns = append(emptyColumns, column)
	}

	return emptyColumns
}

// Print the galaxy map
func printGalaxyMap(galaxyMap []Coordinate) {
	// Find the maximum x and y values
	maxX := 0
	maxY := 0
	for _, coordinate := range galaxyMap {
		if coordinate.x > maxX {
			maxX = coordinate.x
		}
		if coordinate.y > maxY {
			maxY = coordinate.y
		}
	}

	// Create an empty array of rows
	galaxyMapRows := []string{}

	// Loop over the y values
	for y := 0; y <= maxY; y++ {
		// Create an empty row
		row := ""

		// Loop over the x values
		for x := 0; x <= maxX; x++ {
			// Set the character to a space
			char := '.'

			// Loop over the galaxy map
			for _, coordinate := range galaxyMap {
				// If the coordinate matches the x and y values then set the character to a '#'
				if coordinate.x == x && coordinate.y == y {
					char = '#'
				}
			}

			// Add the character to the row
			row += string(char)
		}

		// Add the row to the array of rows
		galaxyMapRows = append(galaxyMapRows, row)
	}

	// Join the rows into a string
	galaxyMapStr := strings.Join(galaxyMapRows, "\n")

	// Print the galaxy map
	fmt.Println(galaxyMapStr)
}

// Parse the galaxy map string into an array of coordinates
func parseGalaxyMap(galaxyMapStr string) []Coordinate {
	// Create an empty array of coordinates
	galaxyMap := []Coordinate{}

	// Split lines into an array
	galaxyMapLines := strings.Split(galaxyMapStr, "\n")

	// Loop through the galaxy map string
	for y, line := range galaxyMapLines {
		for x, char := range line {
			// If the character is a '#' then add the coordinate to the array
			if char == '#' {
				galaxyMap = append(galaxyMap, Coordinate{x, y})
			}
		}
	}

	return galaxyMap
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
