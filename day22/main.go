package main

import (
	"fmt"
	"os"
	"strings"
)

type blockTree struct {
	block brick
}

type brickMap map[int]brick
type Layer []map[int]int

type brick struct {
	// Give brick a unique id when read in so we can reference it later on
	id         int
	start, end Point
}

type Point struct {
	x, y, z int
}

const debug = false

const fallingSandBricksExample1 = `1,0,1~1,2,1
0,0,2~2,0,2
0,2,3~2,2,3
0,0,4~0,2,4
2,0,5~2,2,5
0,1,6~2,1,6
1,1,8~1,1,9`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 22)

	// Load input
	// input := fallingSandBricksExample1
	input := loadFileContents("falling_bricks.txt")

	// Parse input as starting position
	bricks, bottoms, tops := parseInput(input)

	fmt.Printf("Start: %v\n", bricks)

	// printBricks(bricks)

	// make blocks fall, until they all stop.
	// Also build and return a tree that represents the touching blocks
	makeBlocksFall(&bricks, &bottoms, &tops)

	// Print the tree
	// fmt.Printf("Drop : %v\n", bricks)
	// fmt.Println(bottoms)
	// fmt.Println(tops)

	// printBricks(bricks)

	printBrick(bricks[1], 10, 10)

	// Count removable bricks
	removableBricks := countRemovableBricks(&bricks, &bottoms, &tops)
	fmt.Printf("Removable bricks: %d\n", removableBricks)
}

// Make blocks fall, until they all stop.
// Also build and return a tree that represents the touching blocks
func makeBlocksFall(bricks *brickMap, bottoms *Layer, tops *Layer) *blockTree {

	brickFell := true

	// Keep dropping bricks one step until none fall
	for brickFell {
		brickFell = false

		// Start with the base layer and work upwards, dropping any block in the layer above that can drop
		for z := 0; z < len(*bottoms)-1; z++ {
			if debug {
				fmt.Printf("Checking layer %d\n", z+1)
			}
			// Loop over bricks in layer above
			for _, brick := range (*bottoms)[z+1] {
				// Check if brick can drop
				if canDrop(brick, z, bricks, bottoms, tops) {
					// Drop brick
					dropBrick(brick, bricks, bottoms, tops)
					// Mark that a brick fell
					brickFell = true
				}
			}
		}
		if debug {
			fmt.Println()
		}
	}

	return nil
}

// Drop a brick
func dropBrick(brick int, bricks *brickMap, bottoms *Layer, tops *Layer) {
	// Get brick
	b := (*bricks)[brick]

	if debug {
		fmt.Printf("\tDropping brick %d from [%v ➜ %v] to ", brick, b.start, b.end)
	}

	// Remove brick from current layer
	delete((*bottoms)[b.start.z], brick)
	delete((*tops)[b.end.z], brick)

	// Drop brick
	b.start.z--
	b.end.z--
	(*bricks)[brick] = b

	if debug {
		fmt.Printf("[%v ➜ %v]\n", b.start, b.end)
	}

	// Add brick to new layer
	(*bottoms)[b.start.z][brick] = brick
	(*tops)[b.end.z][brick] = brick
}

// Check is a brick can drop
func canDrop(brick int, lowerLayer int, bricks *brickMap, bottoms *Layer, tops *Layer) bool {
	// Get brick
	b := (*bricks)[brick]

	if debug {
		fmt.Printf("\tChecking if brick %d [%v ➜ %v] can drop from %d", brick, b.start, b.end, lowerLayer+1)
	}

	// Check if brick is already at bottom
	if b.end.z == 0 {
		if debug {
			fmt.Printf("  ❌ brick is already at the bottom\n")
		}
		return false
	}

	// Check if brick is sitting on any other bricks from the lower layer
	for _, lowerBrick := range (*bottoms)[lowerLayer] {
		if isSupported((*bricks)[lowerBrick], (*bricks)[brick]) {
			if debug {
				fmt.Printf("  ❌ brick is supported by brick %d\n", lowerBrick)
			}
			return false
		}
	}

	if debug {
		fmt.Println("\n\t ↳ brick is not supported by any bricks in the lower layer")
	}

	// Brick not being held up by anything, so it can drop
	return true
}

// Check if a brick is supported, this is just comparing whether the 2 bricks overlap in the x and y plane
func isSupported(brick, supportingBrick brick) bool {
	// Check if brick is supported by supporting brick
	if brick.start.x > supportingBrick.end.x || brick.end.x < supportingBrick.start.x ||
		brick.start.y > supportingBrick.end.y || brick.end.y < supportingBrick.start.y {
		return false
	}

	// Brick is supported
	return true
}

// Count removable bricks
func countRemovableBricks(bricks *brickMap, bottoms *Layer, tops *Layer) int {
	// Count removable bricks
	var removableBricks int

	// Loop over bricks
	for brick := range *bricks {
		// Check if brick is removable
		if structureSoundIfRemoved(brick, bricks, bottoms, tops) {
			// Brick is removable
			removableBricks++
		}
	}

	return removableBricks
}

// Function that checks if the structure is still sound if a brick is removed
func structureSoundIfRemoved(brick int, bricks *brickMap, bottoms *Layer, tops *Layer) bool {
	// Get brick in question
	b := (*bricks)[brick]

	// Loop over bricks in layer above
	for _, higherBrick := range (*bottoms)[b.end.z+1] {
		supported := false
		// Check if brick can drop by checking all bricks in this bricks layer, except the brick in question
		// Loop over bricks in layer except the brick in question
		for _, lowerBrick := range (*tops)[b.end.z] {
			if lowerBrick != brick {
				// Check if brick will hold it up
				if isSupported((*bricks)[higherBrick], (*bricks)[lowerBrick]) {
					// Brick is supported, so structure is sound. Continue to check next higher brick
					supported = true
					break
				}
			}
		}
		// No other bricks in this layer can support all the brick, so structure is not sound
		if !supported {
			return false
		}
	}
	// fmt.Printf("Brick %c is removable\n", "ABCDEFGHIJKLMNO"[brick])
	fmt.Printf("Brick %d is removable\n", brick)
	// Brick removal doesn't break the structure
	return true
}

// Parse input as starting position
func parseInput(input string) (brickMap, Layer, Layer) {
	// Split input into lines
	lines := strings.Split(input, "\n")

	// Create a slice of bricks
	bricks := make(brickMap, len(lines))

	var maxZ int
	// Parse each line into a brick
	for i, line := range lines {
		bricks[i] = parseLine(line, i)
		maxZ = biggestInt(maxZ, biggestInt(bricks[i].start.z, bricks[i].end.z))
	}

	// Create a slice of layer top & bottoms
	bottoms := make(Layer, maxZ+1)
	// Initialise map for each layer
	for i := range bottoms {
		bottoms[i] = make(map[int]int)
	}
	tops := make(Layer, maxZ+1)
	// Initialise map for each layer
	for i := range tops {
		tops[i] = make(map[int]int)
	}

	// Loop over loaded bricks and find the top and bottom of each layer
	for _, brick := range bricks {
		// Add brick to bottom layer
		bottoms[brick.start.z][brick.id] = brick.id
		tops[brick.end.z][brick.id] = brick.id
	}
	return bricks, bottoms, tops
}

func biggestInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Parse a line into a brick
func parseLine(line string, id int) brick {
	// Split line into start and end points
	points := strings.Split(line, "~")

	// Parse start and end points
	start := parsePoint(points[0])
	end := parsePoint(points[1])

	return brick{id, start, end}
}

// Parse a point
func parsePoint(number string) Point {
	var result [3]int

	// Split point into x, y, z
	_, err := fmt.Sscanf(number, "%d,%d,%d", &result[0], &result[1], &result[2])
	if err != nil {
		panic(err)
	}

	return Point{result[0], result[1], result[2]}
}

// Print bricks
func printBricks(bricks brickMap) {
	// Find max x and y
	var maxX, maxY int
	for _, brick := range bricks {
		maxX = biggestInt(maxX, biggestInt(brick.start.x, brick.end.x))
		maxY = biggestInt(maxY, biggestInt(brick.start.y, brick.end.y))
	}

	fmt.Println("Bricks:")
	for _, brick := range bricks {
		printBrick(brick, maxX, maxY)
	}
	fmt.Println()
}

// Print a brick
func printBrick(b brick, maxX int, maxY int) {
	fmt.Printf("Printing brick %v:\n", b)
	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
			if x >= b.start.x && x <= b.end.x && y >= b.start.y && y <= b.end.y {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
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
