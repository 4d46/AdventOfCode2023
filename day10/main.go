package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type step struct {
	letter rune
	x, y   int
}

const pipeMapExample1 = `.....
.S-7.
.|.|.
.L-J.
.....`

const pipeMapExample2 = `-L|F7
7S-7|
L|7||
-L-J|
L|-JF`

const pipeMapExample3 = `7-F7-
.FJ|7
SJLL7
|F--J
LJ.LJ`

const pipeMapExample4 = `...........
.S-------7.
.|F-----7|.
.||.....||.
.||.....||.
.|L-7.F-J|.
.|..|.|..|.
.L--J.L--J.
...........`

const pipeMapExample5 = `.F----7F7F7F7F-7....
.|F--7||||||||FJ....
.||.FJ||||||||L7....
FJL7L7LJLJ||LJ.L-7..
L--J.L7...LJS7F-7L7.
....F-J..F7FJ|L7L7L7
....L7.F7||L7|.L7L7|
.....|FJLJ|FJ|F7|.LJ
....FJL-7.||.||||...
....L---J.LJ.LJLJ...`

const pipeMapExample6 = `FF7FSF7F7F7F7F7F---7
L|LJ||||||||||||F--J
FL-7LJLJ||||||LJL-77
F--JF--7||LJLJ7F7FJ-
L---JF-JLJ.||-FJLJJ7
|F|F-JF---7F7-L7L|7|
|FFJF7L7F-JF7|JL---7
7-L-JL7||F7|L7F-7F7|
L.L7LFJ|||||FJL7||LJ
L7JLJL-JLJLJL--JLJ.L`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 10)

	// Load the map
	// pipeMapStr := pipeMapExample1
	// pipeMapStr := pipeMapExample2
	// pipeMapStr := pipeMapExample3
	// pipeMapStr := pipeMapExample4
	// pipeMapStr := pipeMapExample5
	// pipeMapStr := pipeMapExample6
	pipeMapStr := loadFileContents("pipe_map.txt")

	// Create the map
	pipeMap := parsePipeMap(pipeMapStr)
	spew.Dump(pipeMap)

	start := findStart(pipeMap)
	fmt.Printf("Start: %v\n", start)

	path := findPath(pipeMap, start)
	// fmt.Printf("Path: %v\n", path)
	// Serialize the path
	fmt.Println(SerializePath(path))
	printPath(len(pipeMap[0]), len(pipeMap), path)

	fmt.Println("Path length:", len(path))
	stepsToFurthestPoint := (len(path) - 1) / 2
	fmt.Println("Steps to furthest point:", stepsToFurthestPoint)

	enclosedAreas := findEnclosedAreas(pipeMap, path)
	fmt.Println("Enclosed Areas")
	fmt.Println(strings.Join(enclosedAreas, "\n"))
	fmt.Printf("Enclosed area count: %d\n", countEnclosedAreas(enclosedAreas))
}

// Count the number of enclosed areas
func countEnclosedAreas(enclosedMap []string) int {
	// Count the number of enclosed areas
	count := 0
	for _, line := range enclosedMap {
		for _, char := range line {
			if char == '▒' {
				count++
			}
		}
	}

	return count
}

// Find enclosed areas
func findEnclosedAreas(pipeMap []string, path []step) []string {
	// Find enclosed areas
	// Create a map
	var mapStr []string = make([]string, len(pipeMap))
	for y := 0; y < len(pipeMap); y++ {
		pipeCount := 0
		for x := 0; x < len(pipeMap[0]); x++ {
			pathIndex := pathCellIndex(path, x, y)
			if pathIndex >= 0 {
				pipeCount += calcHand(path, pathCellIndex(path, x, y))
				mapStr[y] += "+"
			} else {
				if pipeCount >= 2 || pipeCount <= -2 {
					mapStr[y] += "▒"
				} else {
					mapStr[y] += "."
				}
			}
			// if pipeCount >= 0 {
			// 	mapStr[y] += fmt.Sprintf("%d", pipeCount)
			// } else {
			// 	if pipeCount == -1 {
			// 		mapStr[y] += "A"
			// 	} else if pipeCount == -2 {
			// 		mapStr[y] += "B"
			// 	} else {
			// 		mapStr[y] += "Z"
			// 	}
			// }
		}
		// wasPathCell := false
		// for x := 0; x < len(pipeMap[0]); x++ {
		// 	if isPathCell(path, x, y) {
		// 		if wasPathCell {
		// 			pipeCount++
		// 		}
		// 		wasPathCell = true
		// 		mapStr[y] += "+"
		// 	} else {
		// 		if wasPathCell {
		// 			pipeCount++
		// 			wasPathCell = false
		// 		}
		// 		if pipeCount%2 == 1 {
		// 			mapStr[y] += "▒"
		// 		} else {
		// 			mapStr[y] += "."
		// 		}
		// 	}

		// }
	}

	return mapStr
}

// Is this cell part of the path?
func isPathCell(path []step, x int, y int) bool {
	// Is this cell part of the path?
	for _, s := range path {
		if s.x == x && s.y == y {
			return true
		}
	}

	return false
}

func pathCellIndex(path []step, x int, y int) int {
	// Is this cell part of the path?
	for i, s := range path {
		if s.x == x && s.y == y {
			return i
		}
	}

	return -1
}

func calcHand(path []step, index int) int {
	// Calculate the handedness of the path at a given index
	// The handedness is the direction of the next step relative to the current step
	// The handedness is calculated by looking at the direction of the next step relative to the current step
	handedness := 0
	if path[index].y < path[index+1].y {
		// Next step is north, so increase count by 1
		handedness += 1
	} else if path[index].y > path[index+1].y {
		// Next step is south, so decrease count by 1
		handedness -= 1
	}
	if index == 0 {
		index = len(path) - 1
	}
	if path[index-1].y < path[index].y {
		// Previous step is north, so increase count by 1
		handedness += 1
	} else if path[index-1].y > path[index].y {
		// Previous step is south, so decrease count by 1
		handedness -= 1
	}
	// Next step is in the same row, so don't change count
	return handedness

}

// Find the starting point
func findStart(pipeMap []string) step {
	for y, line := range pipeMap {
		for x, char := range line {
			if char == 'S' {
				return step{char, x, y}
			}
		}
	}

	return step{}
}

// Walk the path
func findPath(pipeMap []string, start step) []step {
	// Walk the path
	path := make([]step, 0, 10)
	path = append(path, start)

	// Find the first step
	firstStep := findFirstStep(pipeMap, start)
	path = append(path, firstStep)

	// fmt.Printf("After first step, current path:\n")
	// printPath(len(pipeMap[0]), len(pipeMap), path)

	lastStep := start
	nextStep := firstStep
	for nextStep.letter != 'S' {
		// Remember the last step
		currentStep := nextStep
		// Find the next step
		nextStep = findNextStep(pipeMap, lastStep, nextStep)
		// fmt.Printf("Next step: {%c x:%3d y:%3d}\n", nextStep.letter, nextStep.x, nextStep.y)
		path = append(path, nextStep)
		lastStep = currentStep

		// printPath(len(pipeMap[0]), len(pipeMap), path)
		if nextStep.letter == 0 {
			panic("No next step found")
		}
	}

	return path
}

// Print path on a map
func printPath(width int, height int, path []step) {
	// Print the path on a map
	// Create a map
	var mapStr []string
	for i := 0; i < height; i++ {
		mapStr = append(mapStr, strings.Repeat(".", width))
	}

	// Add the path to the map
	for _, s := range path {
		mapStr[s.y] = mapStr[s.y][:s.x] + string(s.letter) + mapStr[s.y][s.x+1:]
	}

	fmt.Printf("%s\n", strings.Join(mapStr, "\n"))
}

// Find the next step
func findNextStep(pipeMap []string, lastStep step, currentStep step) step {
	// Find the next step
	nextStep := step{}

	// Check the next square in each direction, but take into account the edges of the map
	// The pipes are arranged in a two-dimensional grid of tiles:
	// | is a vertical pipe connecting north and south.
	// - is a horizontal pipe connecting east and west.
	// L is a 90-degree bend connecting north and east.
	// J is a 90-degree bend connecting north and west.
	// 7 is a 90-degree bend connecting south and west.
	// F is a 90-degree bend connecting south and east.
	switch currentStep.letter {
	case '|':
		if currentStep.y-lastStep.y < 0 {
			// last step was north, so next step will be north
			nextStep = step{rune(pipeMap[currentStep.y-1][currentStep.x]), currentStep.x, currentStep.y - 1}
		} else {
			// last step was south, so next step will be south
			nextStep = step{rune(pipeMap[currentStep.y+1][currentStep.x]), currentStep.x, currentStep.y + 1}
		}
	case '-':
		if currentStep.x-lastStep.x < 0 {
			// last step was west, so next step will be west
			nextStep = step{rune(pipeMap[currentStep.y][currentStep.x-1]), currentStep.x - 1, currentStep.y}
		} else {
			// last step was east, so next step will be east
			nextStep = step{rune(pipeMap[currentStep.y][currentStep.x+1]), currentStep.x + 1, currentStep.y}
		}
	case 'L':
		if currentStep.y-lastStep.y == 0 {
			// last step was west, so next step will be north
			nextStep = step{rune(pipeMap[currentStep.y-1][currentStep.x]), currentStep.x, currentStep.y - 1}
		} else {
			// last step was south, so next step will be east
			nextStep = step{rune(pipeMap[currentStep.y][currentStep.x+1]), currentStep.x + 1, currentStep.y}
		}
	case 'J':
		if currentStep.y-lastStep.y == 0 {
			// last step was east, so next step will be north
			nextStep = step{rune(pipeMap[currentStep.y-1][currentStep.x]), currentStep.x, currentStep.y - 1}
		} else {
			// last step was south, so next step will be west
			nextStep = step{rune(pipeMap[currentStep.y][currentStep.x-1]), currentStep.x - 1, currentStep.y}
		}
	case '7':
		if currentStep.y-lastStep.y == 0 {
			// last step was east, so next step will be south
			nextStep = step{rune(pipeMap[currentStep.y+1][currentStep.x]), currentStep.x, currentStep.y + 1}
		} else {
			// last step was north, so next step will be west
			nextStep = step{rune(pipeMap[currentStep.y][currentStep.x-1]), currentStep.x - 1, currentStep.y}
		}
	case 'F':
		if currentStep.y-lastStep.y == 0 {
			// last step was west, so next step will be south
			nextStep = step{rune(pipeMap[currentStep.y+1][currentStep.x]), currentStep.x, currentStep.y + 1}
		} else {
			// last step was north, so next step will be east
			nextStep = step{rune(pipeMap[currentStep.y][currentStep.x+1]), currentStep.x + 1, currentStep.y}
		}
	}

	return nextStep
}

// Find the next step
func findFirstStep(pipeMap []string, currentStep step) step {
	// Find the next step
	nextStep := step{}

	// Check the next square in each direction, but take into account the edges of the map

	// Check the next step in each direction
	// for _, dir := range []string{"↖", "⬆", "↗", "➡", "↘", "⬇", "↙", "⬅"} {
	for _, dir := range []string{"⬆", "➡", "⬇", "⬅"} {
		// Check the next step in this direction
		nextStep = findNextStepInDir(pipeMap, currentStep, dir)
		if nextStep.letter != 0 {
			break
		}
	}

	return nextStep
}

// Find the next step in a direction
func findNextStepInDir(pipeMap []string, currentStep step, dir string) step {
	// Find the next step in a direction
	nextStep := step{}

	// Check the next step in this direction
	switch dir {
	case "⬆":
		if currentStep.y > 0 {
			switch pipeMap[currentStep.y-1][currentStep.x] {
			case '|', '7', 'F':
				nextStep = step{rune(pipeMap[currentStep.y-1][currentStep.x]), currentStep.x, currentStep.y - 1}
			}
		}
	case "➡":
		if currentStep.x < len(pipeMap[0])-1 {
			switch pipeMap[currentStep.y][currentStep.x+1] {
			case '-', '7', 'J':
				nextStep = step{rune(pipeMap[currentStep.y][currentStep.x+1]), currentStep.x + 1, currentStep.y}
			}
		}
	case "⬇":
		if currentStep.y < len(pipeMap)-1 {
			switch pipeMap[currentStep.y+1][currentStep.x] {
			case '|', 'L', 'J':
				nextStep = step{rune(pipeMap[currentStep.y+1][currentStep.x]), currentStep.x, currentStep.y + 1}
			}
		}
	case "⬅":
		if currentStep.x > 0 {
			switch pipeMap[currentStep.y][currentStep.x-1] {
			case '-', 'L', 'F':
				nextStep = step{rune(pipeMap[currentStep.y][currentStep.x-1]), currentStep.x - 1, currentStep.y}
			}
		}
	}
	if nextStep.letter != '.' {
		return nextStep
	}

	return step{}
}

// Parse the pipe map string into a map
func parsePipeMap(pipeMapStr string) []string {
	// Split the string into lines
	lines := splitLines(pipeMapStr)

	return lines
}

// Split a string into lines
func splitLines(str string) []string {
	// Split the string into lines
	lines := make([]string, 0)
	for _, line := range strings.Split(str, "\n") {
		lines = append(lines, line)
	}

	return lines
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

func SerializePath(path []step) string {
	// Serialize the path
	pathStr := "{"
	for i, s := range path {
		if i > 0 {
			pathStr += ","
		}
		pathStr += fmt.Sprintf("{%d,%d,%d}", s.letter, s.x, s.y)
	}
	pathStr += "}"

	return pathStr
}
