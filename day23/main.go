package main

import (
	"fmt"
	"os"
	"strings"
)

const exampleTrail1 = `#.#####################
#.......#########...###
#######.#########.#.###
###.....#.>.>.###.#.###
###v#####.#v#.###.#.###
###.>...#.#.#.....#...#
###v###.#.#.#########.#
###...#.#.#.......#...#
#####.#.#.#######.#.###
#.....#.#.#.......#...#
#.#####.#.#.#########v#
#.#...#...#...###...>.#
#.#.#v#######v###.###v#
#...#.>.#...>.>.#.###.#
#####v#.#.###v#.#.###.#
#.....#...#...#.#.#...#
#.#########.###.#.#.###
#...###...#...#...#.###
###.###.#.###v#####v###
#...#...#.#.>.>.#.>.###
#.###.###.#.###.#.#v###
#.....###...###...#...#
#####################.#`

const (
	// Directions
	north = iota
	east
	south
	west
)

type void struct{}

var EMPTY void

type Point struct {
	x int
	y int
}

type TrailMap struct {
	forest map[Point]bool
	slopes map[Point]rune
	width  int
	height int
	start  Point
	end    Point
}

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 23)

	// Load input plot string
	// trailMapStr := exampleTrail1
	trailMapStr := loadFileContents("hiking_trails.txt")

	// Parse trail map
	trailMap := parseTrailMap(trailMapStr)

	// Print start point & Size
	fmt.Printf("Start: %v ", trailMap.start)
	fmt.Printf("Size: %v x %v\n", trailMap.width, trailMap.height)

	// printTrailMap(trailMap, nil)

	longestPath := walkPath(&trailMap, trailMap.start)

	fmt.Printf("Longest path: %v\n", longestPath)

	// Identify length of path in last square
	// fmt.Printf("Length of path: %v\n", trailMap.plotSteps[trailMap.end])

	// printTrailMap(trailMap, nil)

	// Initialize list of plots that can be reached in the specified number of steps
	// Initialise with start point
	// 	var stepPlots []Point
	// 	stepPlots = append(stepPlots, start)

	// 	const maxSteps = 64
	// 	// Iterate over steps
	// 	for steps := 1; steps <= maxSteps; steps++ {
	// 		// Identify plots that can be reached in the specified number of steps
	// 		stepPlots = identifyStepPlots(garden, steps, stepPlots)
	// 		// spew.Dump(stepPlots)

	// 		if steps&3 == 0 {
	// 			fmt.Printf("Step %v, step plots = %d\n", steps, len(stepPlots))
	// 		}
	// 		// Print garden
	// 		// fmt.Printf("\nStep %v\n", steps)
	// 		// printGarden(garden, steps)
	// 	}

	// // Count plots that can be reached in the specified number of steps
	// fmt.Println("Counting plots...")
	// numPlots := countPlotsReachedBySteps(garden, maxSteps)
	// fmt.Printf("Number of plots that can be reached in %v steps: %v\n", maxSteps, numPlots)
}

func walkPath(tm *TrailMap, start Point) int {
	var longestPath int

	// Start at the start point
	currentPlot := start

	// Take first step
	longestPath = takeNextStep(tm, make(map[Point]void), currentPlot, south, 0)
	fmt.Println()
	return longestPath
}

func takeNextStep(tm *TrailMap, route map[Point]void, currentPlot Point, prevDirection int, count int) int {

	// Check if we have hit the end
	if currentPlot == tm.end {
		// Record number of steps to get to this point, if it is greater than a previous number of steps
		// if prevCount, ok := tm.plotSteps[currentPlot]; !ok || count > prevCount {
		// Print the number of steps to get to this point
		fmt.Printf("Reached end in %v steps\n", count)
		// fmt.Print(".")
		return count
	}

	// If the number of steps to this square is less than than previously, stop walking this path
	// if prevCount, ok := tm.plotSteps[currentPlot]; ok && count <= prevCount {
	// 	return
	// }
	// Check if we have already been here and if so don't continue this path
	if _, ok := route[currentPlot]; ok {
		fmt.Print("x")
		return 0
	}

	// Record number of steps to get to this point
	// route[currentPlot] = count
	// Record that is this plot has been visited
	route[currentPlot] = EMPTY

	longestPath := 0

	// Try to find new points in each direction
	for direction := range []int{north, east, south, west} {
		// Don't go back the way we came
		if direction == oppositeDirection(prevDirection) {
			continue
		}

		// Find coordinates of new point
		nextPlot := calculateNextPlot(currentPlot, direction)

		// Check if new point is outside garden
		if !tm.isInside(nextPlot) {
			// It isn't, so skip this direction
			continue
		}

		// Check if new point is forest, if so ignore
		if tm.forest[nextPlot] {
			continue
		}

		// Ignore for part 2
		// // Check if nextPlot is a slope and if so, whether we can travel down it
		// if slope, ok := tm.slopes[nextPlot]; ok {
		// 	if !canTravelDownSlope(slope, direction) {
		// 		continue
		// 	}
		// }

		// Duplicate the route
		newRoute := make(map[Point]int)
		for k, v := range route {
			newRoute[k] = v
		}
		// Looks valid so take next step
		pathLen := takeNextStep(tm, newRoute, nextPlot, direction, count+1)
		if pathLen > longestPath {
			longestPath = pathLen
		}
	}
	return longestPath
}

func canTravelDownSlope(slope rune, direction int) bool {
	switch slope {
	case '>':
		return direction == east
	case 'v':
		return direction == south
	case '<':
		return direction == west
	case '^':
		return direction == north
	}
	panic("Invalid slope")
}

func oppositeDirection(direction int) int {
	switch direction {
	case north:
		return south
	case east:
		return west
	case south:
		return north
	case west:
		return east
	}
	panic("Invalid direction")
}

// Count plots that can be reached in the specified number of steps
// func countPlotsReachedBySteps(tm TrailMap, steps int) int {
// 	// Initialize count
// 	count := 0
// 	for _, plot := range tm.plotSteps {
// 		if plot[steps] {
// 			count++
// 		}
// 	}
// 	return count
// }

// Identify the plots that can be reached for the specified number of steps
// This assumes previous steps have already been calculated
// func identifyStepPlots(tm TrailMap, steps int, previousStepPlots []Point) []Point {
// 	// Initialize list of plots that can be reached in the specified number of steps
// 	stepPlotMap := make(map[Point]bool)

// 	// Iterate over previous plots
// 	for _, plot := range previousStepPlots {
// 		// Try to find new points in each direction
// 		for _, direction := range []int{north, east, south, west} {
// 			// Find coordinates of new point
// 			nextPlot := calculateNextPlot(plot, direction)

// 			// Check if new point is outside garden
// 			if !tm.isInside(nextPlot) {
// 				// It isn't, so skip this direction
// 				continue
// 			}

// 			// Check if new point is forest
// 			if tm.forest[nextPlot] {
// 				// It is, so skip this direction
// 				continue
// 			}

// 			// Check if nextPlot is a slope and if so, whether we can travel down it
// 			if slope, ok := tm.slopes[nextPlot]; ok {
// 				// It is, so check if we can travel down it
// 				switch slope {
// 				case '>':

// 			// Valid step, so add the step number to the plot steps and add the plot to the list of step plots
// 			garden.plotSteps[nextPlot][steps] = true
// 			if _, ok := stepPlotMap[nextPlot]; !ok {
// 				stepPlotMap[nextPlot] = true
// 			}
// 		}
// 	}
// 	stepPlots := make([]Point, len(stepPlotMap))
// 	i := 0
// 	for k := range stepPlotMap {
// 		stepPlots[i] = k
// 		i++
// 	}

// 	return stepPlots
// }

// Print the trail map
func printTrailMap(tm TrailMap) {
	// Iterate over rows
	for y := 0; y < tm.height; y++ {
		// Iterate over columns
		for x := 0; x < tm.width; x++ {
			// Check if plot is forest
			if tm.forest[Point{x: x, y: y}] {
				// It is, so print it
				fmt.Print("#")
			} else {
				// It hasn't, so determine whether it is a slope
				if slope, ok := tm.slopes[Point{x: x, y: y}]; ok {
					// It is, so print it
					fmt.Printf("%c", slope)
				} else {
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
func (tm TrailMap) isInside(point Point) bool {
	// Check if point is inside garden
	if point.x < 0 || point.x >= tm.width || point.y < 0 || point.y >= tm.height {
		return false
	}

	return true
}

// Parse trail map from string
func parseTrailMap(plotStr string) TrailMap {
	// Initialize garden
	tm := TrailMap{
		forest: make(map[Point]bool),
		slopes: make(map[Point]rune),
	}

	// Initialize start point
	// tm.start = Point{
	// 	x: 1,
	// 	y: 0,
	// }

	// Split the plot string into lines
	plotLines := strings.Split(plotStr, "\n")
	tm.height = len(plotLines)
	tm.width = len(plotLines[0])

	// tm.end = Point{
	// 	x: tm.width - 1,
	// 	y: tm.height - 1,
	// }

	// Iterate over lines
	for y, line := range plotLines {
		// Iterate over characters
		for x, plot := range line {
			switch plot {
			case '.':
				// don't do anything, if it is missing we will add a value later on
				// tm.plotSteps[Point{x: x, y: y}] = math.MaxInt
				// case 'S':
				// 	// Set start point
				// 	start.x = x
				// 	start.y = y
				// 	// Add empty plot to garden
				// 	garden.plotSteps[Point{x: x, y: y}] = make(map[int]bool)
				// If this is on the first line it will be the start point and
				// if it is on the last line it will be the end point
				if y == 0 {
					// Set start point
					tm.start.x = x
					tm.start.y = y
				} else if y == tm.height-1 {
					// Set end point
					tm.end.x = x
					tm.end.y = y
				}
			case '#':
				// This is a forest, add it to Trail Map
				tm.forest[Point{x: x, y: y}] = true
			case '>', 'v', '<', '^':
				// This is a slope, add it to Trail Map
				tm.slopes[Point{x: x, y: y}] = plot
			}
		}
	}

	return tm
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
