package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

const (
	north = iota
	east
	south
	west
)

// Struct to remember the lowest total heat loss upto this block
// Also remember the last block that was visited
type block struct {
	heatLossSum  int
	lastX, lastY int
	lastDir      int
	lastStreak   int
}

// 2D map of blocks, with a "depth" representing the number squares taken in 1 direction
// to reach the current square in the map
type coolingMap struct {
	heatloss   [][]int
	computed   [][][3]block
	maxX, maxY int
}

const exampleMapStr = `2413432311323
3215453535623
3255245654254
3446585845452
4546657867536
1438598798454
4457876987766
3637877979653
4654967986887
4564679986453
1224686865563
2546548887735
4322674655533`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 17)

	// Load example map
	input := exampleMapStr

	// Parse input map
	heatLossMap := parseInputMap(input)

	// Print input map
	fmt.Println("Input map:")
	printMap(heatLossMap)

	// Compute heat loss sums
	computeHeatSums(&heatLossMap)

	// Print computed map
	printComputedMap(heatLossMap)

	// Find lowest heat loss
	lowestHeatLoss := findLowestHeatLoss(heatLossMap)
	fmt.Printf("Lowest heat loss: %d\n", lowestHeatLoss)

	// Debug
	spew.Dump(heatLossMap.computed[heatLossMap.maxX-1][heatLossMap.maxY-2])
	spew.Dump(heatLossMap.computed[heatLossMap.maxX-2][heatLossMap.maxY-1])

}

func computeHeatSums(heatLossMap *coolingMap) {
	// // Define start point at the end of the map and work backwards
	// startX, startY := heatLossMap.maxX-1, heatLossMap.maxY-1
	// Define start point at the start of the map and work forwards
	startX, startY := 0, 0

	// From start point start working backwards
	calculateNextStep(heatLossMap, startX, startY, -1, 0, 0, true)
}

func calculateNextStep(heatLossMap *coolingMap, posX, posY, dir, sumSoFar int, streak int, firstStep bool) {

	// Record sum to this point if
	sumSoFar += heatLossMap.heatloss[posY][posX]
	// If sum isn't lower than the previous lowest sum to this block, stop processing
	if sumSoFar >= heatLossMap.computed[posY][posX][streak].heatLossSum {
		return
	}
	// This must be the new lowest sum to this block, record it and then continue
	heatLossMap.computed[posY][posX][streak].heatLossSum = sumSoFar
	if firstStep {
		heatLossMap.computed[posY][posX][streak].lastX = posX
		heatLossMap.computed[posY][posX][streak].lastY = posY
	} else {
		lastX, lastY := calculateNextPosition(posX, posY, oppositeDirection(dir))
		heatLossMap.computed[posY][posX][streak].lastX = lastX
		heatLossMap.computed[posY][posX][streak].lastY = lastY
	}

	// Loop through all directions
	for nextDir := range [4]int{north, east, south, west} {
		// Don't loop back on yourself, skip if next direction is opposite of current direction
		// unless this is the first step where the direction is invalid
		if !firstStep && nextDir == oppositeDirection(dir) {
			continue
		}
		// If we have been in this direction for 3 blocks, we can't go any further, skip
		nextStreak := 0
		if nextDir == dir {
			nextStreak = streak + 1
		}
		if nextStreak >= 3 {
			continue
		}
		// Calculate next position
		nextX, nextY := calculateNextPosition(posX, posY, nextDir)

		// If next position is outside of map, skip
		if outsideMap(heatLossMap, nextX, nextY) {
			continue
		}
		// Calculate next step
		calculateNextStep(heatLossMap, nextX, nextY, nextDir, sumSoFar, nextStreak, false)
	}
}

func calculateNextPosition(posX, posY, dir int) (int, int) {
	switch dir {
	case north:
		return posX, posY - 1
	case east:
		return posX + 1, posY
	case south:
		return posX, posY + 1
	case west:
		return posX - 1, posY
	default:
		panic("Invalid direction")
	}
}

func oppositeDirection(dir int) int {
	switch dir {
	case north:
		return south
	case east:
		return west
	case south:
		return north
	case west:
		return east
	default:
		panic("Invalid direction")
	}
}

func outsideMap(heatLossMap *coolingMap, posX, posY int) bool {
	return posX < 0 || posX >= heatLossMap.maxX || posY < 0 || posY >= heatLossMap.maxY
}

// Parse input map
func parseInputMap(input string) coolingMap {
	// Split input into lines
	lines := append([]string{}, strings.Split(input, "\n")...)

	// Create map
	heatLossMap := coolingMap{}
	heatLossMap.heatloss = make([][]int, len(lines))
	heatLossMap.computed = make([][][3]block, len(lines))

	// Parse lines
	for y, line := range lines {
		heatLossMap.heatloss[y] = make([]int, len(line))
		for x, char := range line {
			heatLossMap.heatloss[y][x] = int(char - '0')
		}
		heatLossMap.computed[y] = make([][3]block, len(line))
		for x := range heatLossMap.computed[y] {
			for streak := 0; streak < 3; streak++ {
				heatLossMap.computed[y][x][streak] = block{heatLossSum: math.MaxInt, lastX: -1, lastY: -1, lastDir: -1, lastStreak: -1}
			}
		}
	}
	heatLossMap.maxY = len(lines)
	heatLossMap.maxX = len(lines[0])

	return heatLossMap
}

// Print map
func printMap(heatLossMap coolingMap) {
	for _, line := range heatLossMap.heatloss {
		for _, block := range line {
			fmt.Printf("%d", block)
		}
		fmt.Println()
	}
}

// Print computed map
func printComputedMap(heatLossMap coolingMap) {
	// Rune map
	runeMap := make([][]rune, heatLossMap.maxY)
	for y := range runeMap {
		runeMap[y] = make([]rune, heatLossMap.maxX)
		for x := range runeMap[y] {
			runeMap[y][x] = '.'
		}
	}

	// nextX, nextY := 0, 0
	nextX, nextY := heatLossMap.maxX-1, heatLossMap.maxY-1
	// // Loop until we reach the end of the map
	// for nextX != heatLossMap.maxX-1 || nextY != heatLossMap.maxY-1 {
	// Loop until we reach the start of the map
	for nextX != 0 || nextY != 0 {
		runeMap[nextY][nextX] = '#'
		// Find next block
		// Loop over all 3 streaks and find the lowest heat loss sum
		nextLowestBlock := heatLossMap.computed[nextY][nextX][0]
		for _, streakBlock := range heatLossMap.computed[nextY][nextX] {
			if streakBlock.heatLossSum < nextLowestBlock.heatLossSum {
				nextLowestBlock = streakBlock
			}
		}
		nextX, nextY = nextLowestBlock.lastX, nextLowestBlock.lastY
	}
	runeMap[0][0] = '#'
	for _, line := range runeMap {
		fmt.Println(string(line))
	}
}

func findLowestHeatLoss(heatLossMap coolingMap) int {
	// Find lowest heat loss sum for the value in the bottom right corner
	lowestHeatLoss := math.MaxInt
	// for _, streak := range heatLossMap.computed[0][0] {
	for _, streak := range heatLossMap.computed[heatLossMap.maxY-1][heatLossMap.maxX-1] {
		if streak.heatLossSum < lowestHeatLoss {
			lowestHeatLoss = streak.heatLossSum
		}
	}
	return lowestHeatLoss
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
