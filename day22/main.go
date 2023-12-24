package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"
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

	// fmt.Printf("P Brick 1: {%s}\n", bricks[48].String())
	// printBrick(bricks[48], 10, 10)
	// printBrick(bricks[804], 10, 10)
	// printBrick(bricks[1135], 10, 10)

	// for i := 0; i < 6; i++ {
	// 	fmt.Printf("Layer %d %v:\n", i, tops[i])
	// }

	fmt.Printf("Layers ")

	// fmt.Printf("Start: %v\n", bricks)

	// printBricks(bricks)

	// make blocks fall, until they all stop.
	// Also build and return a tree that represents the touching blocks
	makeBlocksFall(&bricks, &bottoms, &tops)

	// Print the brick 48
	// printBrick(bricks[48], 10, 10)
	// // Print the bricks whose tops are in layer 3
	// for _, brick := range tops[3] {
	// 	switch brick {
	// 	case 475, 192, 269, 539, 883, 646, 693, 142, 1157, 242, 767, 280, 77, 205, 467, 776:
	// 		continue
	// 	}
	// 	printBrick(bricks[brick], 10, 10)
	// }

	// Print the tree
	// fmt.Printf("Drop : %v\n", bricks)
	// fmt.Println(bottoms)
	// fmt.Println(tops)

	// printBricks(bricks)

	// printBrick(bricks[1], 10, 10)
	// fmt.Printf("D Brick 1: {%s}\n", bricks[48].String())
	// // for i := 0; i < len(bottoms); i++ {
	// for i := 0; i < 10; i++ {
	// 	printLayerBricks(bricks, i)
	// }

	// Count removable bricks
	removableBricksCount := countRemovableBricks(&bricks, &bottoms, &tops)
	fmt.Printf("Removable bricks: %d\n", removableBricksCount)

	// Part 2
	// Print the bricks in the first layer where bottoms occur
	printedLayers := 0
	for j := len(bottoms) - 1; j > 0 && printedLayers < 1; j-- {
		fmt.Printf("len bottoms: %d; j: %d\n", len(bottoms), j)
		if len(bottoms[j]) > 0 {
			printLayerBricks(bricks, j)
			printedLayers++
		}
	}

	// Find sum of bricks combinations that would disintegrate, for every brick
	var sum int
	// droppedBrickSumCache := make(map[int]int)
	// // Loop over layers from top to bottom
	// for z := len(bottoms) - 1; z >= 0; z-- {
	// 	// Loop over bricks in layer, so we calculate each brick once but from top to bottom
	// 	for _, brick := range tops[z] {
	// 		droppedBricks := droppedBricksIfRemoved(brick, &bricks, &bottoms, &tops)
	// 		// Calculate sum of dropped bricks, if this brick was removed
	// 		total := 0
	// 		for _, droppedBrick := range droppedBricks {
	// 			droppedBrickSum, ok := droppedBrickSumCache[droppedBrick]
	// 			if !ok {
	// 				panic(fmt.Sprintf("Dropped brick %d not in cache", droppedBrick))
	// 			}
	// 			total += droppedBrickSum + 1
	// 		}
	// 		// Add sum of dropped bricks to cache for this brick
	// 		droppedBrickSumCache[brick] = total
	// 		sum += total
	// 	}
	// }
	sum = countBlocksFall(&bricks, &bottoms, &tops)
	fmt.Printf("Part 2 Sum of brick removal disintegrations: %d\n", sum)
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

func countBlocksFall(bricks *brickMap, bottoms *Layer, tops *Layer) int {
	var count int

	// Keep dropping bricks one step until none fall
	for i, removedBrick := range *bricks {
		if i&32 == 0 {
			fmt.Printf("Checking brick %d\n", i)
		}
		// Create a new deep copy of bricks, bottoms and tops
		lbricks := make(brickMap, len(*bricks))
		for k, v := range *bricks {
			lbricks[k] = v
		}
		lbottoms := make(Layer, len(*bottoms))
		for k, v := range *bottoms {
			lbottoms[k] = make(map[int]int, len(v))
			for k2, v2 := range v {
				lbottoms[k][k2] = v2
			}
		}
		ltops := make(Layer, len(*tops))
		for k, v := range *tops {
			ltops[k] = make(map[int]int, len(v))
			for k2, v2 := range v {
				ltops[k][k2] = v2
			}
		}

		// Remove brick from bricks
		delete(lbricks, removedBrick.id)
		// Remove brick from bottoms
		delete(lbottoms[(*bricks)[removedBrick.id].start.z], removedBrick.id)
		// Remove brick from tops
		delete(ltops[(*bricks)[removedBrick.id].end.z], removedBrick.id)

		// Start with the layer above the removed brick work upwards, dropping any block in the layer above that can drop
		for z := (*bricks)[removedBrick.id].start.z; z < len(lbottoms)-1; z++ {
			// for z := 0; z < len(lbottoms)-1; z++ {
			if debug {
				fmt.Printf("Checking layer %d\n", z+1)
			}
			// Loop over bricks in layer above
			for _, brick := range (lbottoms)[z+1] {
				// Check if brick can drop
				if canDrop(brick, z, &lbricks, &lbottoms, &ltops) {
					// Drop brick
					dropBrick(brick, &lbricks, &lbottoms, &ltops)
					count++
				}
			}
		}
		if debug {
			fmt.Println()
		}
	}

	return count
}

// Drop a brick
func dropBrick(brick int, bricks *brickMap, bottoms *Layer, tops *Layer) {
	// Get brick
	b := (*bricks)[brick]
	ldebug := debug
	if b.id == -1 {
		ldebug = true
	}

	if ldebug {
		fmt.Printf("\tDropping brick %d from [%v ➜ %v] to ", brick, b.start, b.end)
	}

	// Remove brick from current layer
	delete((*bottoms)[b.start.z], brick)
	delete((*tops)[b.end.z], brick)

	// Drop brick
	b.start.z--
	b.end.z--
	(*bricks)[brick] = b

	if ldebug {
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
	ldebug := debug
	if b.id == -1 {
		ldebug = true
	}

	if ldebug {
		fmt.Printf("\tChecking if brick %d [%v ➜ %v] can drop from %d\n", brick, b.start, b.end, lowerLayer+1)
	}

	// Check if brick is already at bottom
	if b.start.z == 1 {
		if ldebug {
			fmt.Printf("  ❌ brick is already at the bottom\n")
		}
		return false
	}

	// Check if brick is sitting on any other bricks from the lower layer
	for _, lowerBrick := range (*tops)[lowerLayer] {
		if isSupported((*bricks)[lowerBrick], (*bricks)[brick]) {
			if ldebug {
				fmt.Printf("  ❌ brick is supported by brick %d{%s}\n", lowerBrick, (*bricks)[lowerBrick].String())
			}
			return false
		}
	}

	if ldebug {
		fmt.Println("\n\t ↳ brick is not supported by any bricks in the lower layer\n")
	}

	// Brick not being held up by anything, so it can drop
	return true
}

// Check if a brick is supported, this is just comparing whether the 2 bricks overlap in the x and y plane
func isSupported(brick, supportingBrick brick) bool {
	// if brick.id == 48 && supportingBrick.id == 804 {
	// 	fmt.Printf("***Checking if brick %d is supported by brick %d\t{%s}{%s}", brick.id, supportingBrick.id, brick.String(), supportingBrick.String())
	// }
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
	// fmt.Printf("Brick %d is removable\n", brick)
	// fmt.Printf("%04d {%d,%d,%d-%d,%d,%d}\n", (*bricks)[brick].id, (*bricks)[brick].start.x, (*bricks)[brick].start.y, (*bricks)[brick].start.z, (*bricks)[brick].end.x, (*bricks)[brick].end.y, (*bricks)[brick].end.z)
	// Brick removal doesn't break the structure
	return true
}

// Function that returns the brick ids that would fall if the specified brick was removed
func droppedBricksIfRemoved(brick int, bricks *brickMap, bottoms *Layer, tops *Layer) []int {
	var droppedBricks []int

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
		// With candidate brick removed, no other bricks in this layer can support this brick, add it to the list of dropped bricks
		if !supported {
			droppedBricks = append(droppedBricks, higherBrick)
		}
	}
	// fmt.Printf("Brick %c is removable\n", "ABCDEFGHIJKLMNO"[brick])
	// fmt.Printf("Brick %d is removable\n", brick)
	// fmt.Printf("%04d {%d,%d,%d-%d,%d,%d}\n", (*bricks)[brick].id, (*bricks)[brick].start.x, (*bricks)[brick].start.y, (*bricks)[brick].start.z, (*bricks)[brick].end.x, (*bricks)[brick].end.y, (*bricks)[brick].end.z)
	// Brick removal doesn't break the structure
	return droppedBricks
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

// Print bricks in a layer
func printLayerBricks(bricks brickMap, layer int) {
	var layerBricks []brick
	for _, brick := range bricks {
		if brick.start.z == layer {
			layerBricks = append(layerBricks, brick)
		}
	}
	if len(layerBricks) == 0 {
		// Nothing in the layer
		return
	}
	fmt.Printf("Layer %d:\n", layer)

	// Sort bricks by id
	slices.SortFunc(layerBricks, func(a, b brick) int { return cmp.Compare(a.id, b.id) })
	// Print bricks
	for _, brick := range layerBricks {
		fmt.Printf("L%3d B%04d {%d,%d,%d-%d,%d,%d}\n", layer, brick.id, brick.start.x, brick.start.y, brick.start.z, brick.end.x, brick.end.y, brick.end.z)
	}
	fmt.Println()
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
	fmt.Printf("Printing brick [%4d] %v:\n", b.id, b)
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

// Serialise brick details
func (b brick) String() string {
	return fmt.Sprintf("%d,%d,%d~%d,%d,%d", b.start.x, b.start.y, b.start.z, b.end.x, b.end.y, b.end.z)
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
