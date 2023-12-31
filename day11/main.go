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
	// galaxyMapStr := GalaxyMapExample1
	galaxyMapStr := loadFileContents("galaxy_map.txt")

	// Parse the input data
	galaxyMap := parseGalaxyMap(galaxyMapStr)
	// Print the galaxy map
	fmt.Println(galaxyMap)
	// printGalaxyMap(galaxyMap)

	// Expand the Galaxy Map
	blankRows := findEmptyRows(galaxyMap)
	fmt.Println("Blank Rows:")
	fmt.Println(blankRows)
	blankColumns := findEmptyColumns(galaxyMap)
	fmt.Println("Blank Columns:")
	fmt.Println(blankColumns)
	// Expand the galaxy map
	expandedGalaxyMap := expandGalaxyMap(galaxyMap, blankRows, blankColumns)
	printGalaxyMapDiff(galaxyMap, expandedGalaxyMap)
	// Print the expanded galaxy map
	// fmt.Println(expandedGalaxyMap)
	// Print the expanded galaxy map
	// printGalaxyMap(expandedGalaxyMap)

	// Calculate the distances between all galaxies
	distances := calculateDistances(expandedGalaxyMap)
	// fmt.Println(distances)
	totalDistance := calculateSumOfDistances(distances)
	fmt.Printf("Sum of all distances: %d\n", totalDistance)
}

// Function to print the difference between two galaxy maps
func printGalaxyMapDiff(galaxyMap1 []Coordinate, galaxyMap2 []Coordinate) {
	// Loop over the galaxy map
	for i, _ := range galaxyMap1 {
		// Print the galaxy map
		fmt.Printf("%d: %d,%d ➞ %d,%d\n", i, galaxyMap1[i].x, galaxyMap1[i].y, galaxyMap2[i].x, galaxyMap2[i].y)
	}
}

// Calculate the sum of all distance between all galaxies
func calculateSumOfDistances(distances map[string]int) int {
	// Create an empty array of distances
	sumOfDistances := 0

	// Loop over the distances
	for _, distance := range distances {
		// Add the distance to the sum of distances
		sumOfDistances += distance
	}

	// Return the sum of distances
	return sumOfDistances
}

// Calculate the distances between all galaxies
func calculateDistances(galaxyMap []Coordinate) map[string]int {
	// Create an empty map of distances
	distances := map[string]int{}

	// Loop over the galaxy map
	for i, galaxy1 := range galaxyMap {
		// Loop over the galaxy map
		for j, galaxy2 := range galaxyMap {
			if i == j {
				continue
			}

			// Calculate the distance between the two galaxies
			distance := calculateDistance(galaxy1, galaxy2)

			var name string
			if i < j {
				name = fmt.Sprintf("%d_%d", i, j)
			} else {
				name = fmt.Sprintf("%d_%d", j, i)
			}
			// Add the distance to the map of distances
			distances[name] = distance
		}
	}

	// Return the map of distances
	return distances
}

// func calculateDistances(galaxyMap []Coordinate) map[int]map[int]int {
// 	// Create an empty map of distances
// 	distances := map[int]map[int]int{}

// 	// Loop over the galaxy map
// 	for i, galaxy1 := range galaxyMap {
// 		// Create an empty map of distances
// 		distances[i] = map[int]int{}

// 		// Loop over the galaxy map
// 		for j, galaxy2 := range galaxyMap {
// 			// Calculate the distance between the two galaxies
// 			distance := calculateDistance(galaxy1, galaxy2)

// 			// Add the distance to the map of distances
// 			distances[i][j] = distance
// 		}
// 	}

// 	// Return the map of distances
// 	return distances
// }

// Calculate the distance between two galaxies in steps
func calculateDistance(galaxy1 Coordinate, galaxy2 Coordinate) int {
	// Calculate the distance between the two galaxies
	distance := abs(galaxy1.x-galaxy2.x) + abs(galaxy1.y-galaxy2.y)

	// Return the distance
	return distance
}

// Integer Abs function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
			// Increment the y value if the blank row is above the original coordinate
			if blankRow < coordinate.y {
				// y++
				// Replacing the 1 row with 1000000 rows, so add another 999999 to the y value
				y += 999999
			}
		}

		// Loop over the blank columns
		for _, blankColumn := range blankColumns {
			// Increment the x value if the blank column is to the left of the coordinate
			if blankColumn < coordinate.x {
				// x++
				// Replacing the 1 column with 1000000 columns, so add another 999999 to the x value
				x += 999999
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
