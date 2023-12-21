package main

import (
	"fmt"
	"os"
	"strings"
)

const examplePlot1 = `...........
.....###.#.
.###.##..#.
..#.#...#..
....#.#....
.##..S####.
.##..#...#.
.......##..
.##.#.####.
.##..##.##.
...........`

const (
	// Directions
	north = iota
	east
	south
	west
)

type Point struct {
	x int
	y int
}

type Garden struct {
	rocks     map[Point]bool
	plotSteps map[Point]map[int]bool
	width     int
	height    int
}

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 21)

	// Load input plot string
	// plotStr := examplePlot1
	plotStr := loadFileContents("garden.txt")

	// Parse garden
	garden, start := parseGarden(plotStr)

	// Print start point & Size
	fmt.Printf("Start: %v ", start)
	fmt.Printf("Size: %v x %v\n", garden.width, garden.height)

	// Initialize list of plots that can be reached in the specified number of steps
	// Initialise with start point
	var stepPlots []Point
	stepPlots = append(stepPlots, start)

	const maxSteps = 64
	// Iterate over steps
	for steps := 1; steps <= maxSteps; steps++ {
		// Identify plots that can be reached in the specified number of steps
		stepPlots = identifyStepPlots(garden, steps, stepPlots)
		// spew.Dump(stepPlots)

		if steps&3 == 0 {
			fmt.Printf("Step %v, step plots = %d\n", steps, len(stepPlots))
		}
		// Print garden
		// fmt.Printf("\nStep %v\n", steps)
		// printGarden(garden, steps)
	}

	// Count plots that can be reached in the specified number of steps
	fmt.Println("Counting plots...")
	numPlots := countPlotsReachedBySteps(garden, maxSteps)
	fmt.Printf("Number of plots that can be reached in %v steps: %v\n", maxSteps, numPlots)
}

// Count plots that can be reached in the specified number of steps
func countPlotsReachedBySteps(garden Garden, steps int) int {
	// Initialize count
	count := 0
	for _, plot := range garden.plotSteps {
		if plot[steps] {
			count++
		}
	}
	return count
}

// Identify the plots that can be reached for the specified number of steps
// This assumes previous steps have already been calculated
func identifyStepPlots(garden Garden, steps int, previousStepPlots []Point) []Point {
	// Initialize list of plots that can be reached in the specified number of steps
	stepPlotMap := make(map[Point]bool)

	// Iterate over previous plots
	for _, plot := range previousStepPlots {
		// Try to find new points in each direction
		for _, direction := range []int{north, east, south, west} {
			// Find coordinates of new point
			nextPlot := calculateNextPlot(plot, direction)

			// Check if new point is outside garden
			if !garden.isInside(nextPlot) {
				// It isn't, so skip this direction
				continue
			}

			// Check if new point is a rock
			if garden.rocks[nextPlot] {
				// It is, so skip this direction
				continue
			}

			// Valid step, so add the step number to the plot steps and add the plot to the list of step plots
			garden.plotSteps[nextPlot][steps] = true
			if _, ok := stepPlotMap[nextPlot]; !ok {
				stepPlotMap[nextPlot] = true
			}
		}
	}
	stepPlots := make([]Point, len(stepPlotMap))
	i := 0
	for k := range stepPlotMap {
		stepPlots[i] = k
		i++
	}

	return stepPlots
}

// Print the garden, for a given number of steps
func printGarden(garden Garden, steps int) {
	// Iterate over rows
	for y := 0; y < garden.height; y++ {
		// Iterate over columns
		for x := 0; x < garden.width; x++ {
			// Check if plot is a rock
			if garden.rocks[Point{x: x, y: y}] {
				// It is, so print it
				fmt.Print("#")
			} else {
				// It isn't, so check if plot has been reached in the specified number of steps
				if garden.plotSteps[Point{x: x, y: y}][steps] {
					// It has, so print it
					fmt.Print("O")
				} else {
					// It hasn't, so print empty plot
					fmt.Print(".")
				}
			}
		}

		// Print new line
		fmt.Println()
	}
}

// Calculate the coordinates of the next plot in the specified direction
func calculateNextPlot(plot Point, direction int) Point {
	// Initialize next plot
	nextPlot := Point{
		x: plot.x,
		y: plot.y,
	}

	// Calculate coordinates of next plot
	switch direction {
	case north:
		nextPlot.y--
	case east:
		nextPlot.x++
	case south:
		nextPlot.y++
	case west:
		nextPlot.x--
	}

	return nextPlot
}

// Function to identify if Point is inside the garden
func (garden Garden) isInside(point Point) bool {
	// Check if point is inside garden
	if point.x < 0 || point.x >= garden.width || point.y < 0 || point.y >= garden.height {
		return false
	}

	return true
}

// Parse garden from string
func parseGarden(plotStr string) (Garden, Point) {
	// Initialize garden
	garden := Garden{
		rocks:     make(map[Point]bool),
		plotSteps: make(map[Point]map[int]bool),
	}

	// Initialize start point
	start := Point{
		x: 0,
		y: 0,
	}

	// Split the plot string into lines
	plotLines := strings.Split(plotStr, "\n")
	garden.height = len(plotLines)
	garden.width = len(plotLines[0])

	// Iterate over lines
	for y, line := range plotLines {
		// Iterate over characters
		for x, plot := range line {
			switch plot {
			case '.':
				// Add empty plot to garden
				garden.plotSteps[Point{x: x, y: y}] = make(map[int]bool)
			case 'S':
				// Set start point
				start.x = x
				start.y = y
				// Add empty plot to garden
				garden.plotSteps[Point{x: x, y: y}] = make(map[int]bool)
			case '#':
				// Add rock to garden
				garden.rocks[Point{x: x, y: y}] = true
			}
		}
	}

	return garden, start
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
