package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	north = iota
	east
	south
	west
)

type coord struct {
	x, y int
}

type trench struct {
	direction int
	distance  int
	color     string
}

type border struct {
	width, height int
	start         coord
	points        [][]borderPoint
}

type borderMap struct {
	width, height int
	start         coord
	points        map[uint64]borderPoint
}

type borderPoint struct {
	pos, lastPos, nextPos coord
	direction             int
	colour                string
}

const ExamplePath1 = `R 6 (#70c710)
D 5 (#0dc571)
L 2 (#5713f0)
D 2 (#d2c081)
R 2 (#59c680)
D 2 (#411b91)
L 5 (#8ceee2)
U 2 (#caa173)
L 1 (#1b58a2)
U 2 (#caa171)
R 2 (#7807d2)
U 3 (#a77fa3)
L 2 (#015232)
U 2 (#7a21e3)`

// north = 3, east = 0, south = 1, west = 2

// Example 2 represents the following path:
// ###  ###
// # #  # #
// # #### #
// #      #
// ########

const ExamplePath2 = `R 2 (#000020)
D 2 (#000021)
R 3 (#000030)
U 2 (#000023)
R 2 (#000020)
D 4 (#000041)
L 7 (#000072)
U 4 (#000043)`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 18)

	// Load input
	// input := ExamplePath1
	// input := ExamplePath2
	input := loadFileContents("trenches.txt")

	// Parse input
	fmt.Print("Parsing Input... ")
	trenches := parsePath(input)
	fmt.Println("Done")

	// Print trenches
	// printTrenches(trenches)

	// Trace path border
	// border, volume := tracePathBorder(trenches)
	fmt.Print("Tracing Path Border... ")
	border, pointsOnRow := tracePathBorderMap(trenches)
	fmt.Println("Done")

	// Print border
	// printBorder(border)

	// Determine interior map
	// interiorMap := determineInteriorMap(&border)

	// Print interior map
	// printInteriorMap(interiorMap)

	// Print complete map
	// printCompleteMap(&border, interiorMap)

	// Count trench & interior points
	// volume := countTrenchAndInteriorPoints(&border, interiorMap)

	// volume += countInteriorMapPoints(&border, &pointsOnRow)
	fmt.Print("Counting Map Points... ")
	volume := countInteriorMapPointsAndTrench(&border, &pointsOnRow)
	fmt.Println("Done")

	// Print volume
	fmt.Println("Volume:", volume)
}

// Encode coordinate as uint64
func encodeCoordAsUint64(coord coord) uint64 {
	return uint64(coord.y)<<32 | uint64(coord.x)
}

// Decode uint64 as coordinate
func decodeUint64AsCoord(u uint64) coord {
	return coord{int(u), int(u >> 32)}
}

// Find maxSize of grid
func findMaxSize(trenches []trench) (coord, coord) {
	maxX, maxY := 0, 0
	minX, minY := 0, 0
	currPos := coord{0, 0}
	for _, trench := range trenches {
		switch trench.direction {
		case north:
			currPos.y -= trench.distance
		case east:
			currPos.x += trench.distance
		case south:
			currPos.y += trench.distance
		case west:
			currPos.x -= trench.distance
		}

		if currPos.x < minX {
			minX = currPos.x
		}
		if currPos.y < minY {
			minY = currPos.y
		}
		if currPos.x > maxX {
			maxX = currPos.x
		}
		if currPos.y > maxY {
			maxY = currPos.y
		}
	}

	return coord{minX, minY}, coord{maxX, maxY}
}

// Trace path border
func tracePathBorder(trenches []trench) (border, int) {
	count := 0
	borderPath := border{}

	// Find boundaries
	minPos, maxPos := findMaxSize(trenches)
	borderPath.width = maxPos.x - minPos.x + 1
	borderPath.height = maxPos.y - minPos.y + 1
	borderPath.start = coord{0 - minPos.x, 0 - minPos.y}
	fmt.Println("start:", borderPath.start)
	fmt.Println("width:", borderPath.width, "height:", borderPath.height)

	// Create grid of the right size
	borderPath.points = make([][]borderPoint, borderPath.height)
	for j := range borderPath.points {
		borderPath.points[j] = make([]borderPoint, borderPath.width)
		for i := range borderPath.points[j] {
			borderPath.points[j][i] = borderPoint{
				pos:       coord{i, j},
				direction: -1,
				colour:    "",
			}
		}
	}

	pos := borderPath.start
	lastPos := pos
	firstNewPoint := coord{-1, -1}

	// // Add first point
	// borderPath.points[pos.y][pos.x] = borderPoint{
	// 	pos:       pos,
	// 	lastPos:   coord{-1, -1},
	// 	direction: trenches[0].direction,
	// 	colour:    trenches[0].color,
	// }
	// Loop over trenches
	for _, trench := range trenches {
		// Calculate new position
		newPos := pos
		switch trench.direction {
		case north:
			newPos.y -= trench.distance
			for j := pos.y - 1; j >= newPos.y; j-- {
				count++
				thisPos := coord{pos.x, j}
				borderPath.points[j][pos.x] = borderPoint{
					pos:       thisPos,
					lastPos:   lastPos,
					direction: trench.direction,
					colour:    trench.color,
				}
				// Set nextPos of last point to this point
				borderPath.points[lastPos.y][lastPos.x].nextPos = thisPos

				// spew.Dump(borderPath.points[j][pos.x])
				if firstNewPoint.x == -1 {
					firstNewPoint = thisPos
				}

				lastPos = thisPos

			}
		case east:
			newPos.x += trench.distance
			for i := pos.x + 1; i <= newPos.x; i++ {
				count++
				thisPos := coord{i, pos.y}
				borderPath.points[pos.y][i] = borderPoint{
					pos:       thisPos,
					lastPos:   lastPos,
					direction: trench.direction,
					colour:    trench.color,
				}
				// Set nextPos of last point to this point
				borderPath.points[lastPos.y][lastPos.x].nextPos = thisPos

				// spew.Dump(borderPath.points[pos.y][i])
				if firstNewPoint.x == -1 {
					firstNewPoint = thisPos
				}

				lastPos = thisPos

			}
		case south:
			newPos.y += trench.distance
			for j := pos.y + 1; j <= newPos.y; j++ {
				count++
				thisPos := coord{pos.x, j}
				borderPath.points[j][pos.x] = borderPoint{
					pos:       thisPos,
					lastPos:   lastPos,
					direction: trench.direction,
					colour:    trench.color,
				}
				// Set nextPos of last point to this point
				borderPath.points[lastPos.y][lastPos.x].nextPos = thisPos

				// spew.Dump(borderPath.points[j][pos.x])
				if firstNewPoint.x == -1 {
					firstNewPoint = thisPos
				}

				lastPos = thisPos

			}
		case west:
			newPos.x -= trench.distance
			for i := pos.x - 1; i >= newPos.x; i-- {
				count++
				thisPos := coord{i, pos.y}
				borderPath.points[pos.y][i] = borderPoint{
					pos:       thisPos,
					lastPos:   lastPos,
					direction: trench.direction,
					colour:    trench.color,
				}
				// Set nextPos of last point to this point
				borderPath.points[lastPos.y][lastPos.x].nextPos = thisPos

				// spew.Dump(borderPath.points[pos.y][i])
				if firstNewPoint.x == -1 {
					firstNewPoint = thisPos
				}

				lastPos = thisPos

			}
		}

		pos = newPos
	}
	// Connect last point to first point
	borderPath.points[lastPos.y][lastPos.x].nextPos = firstNewPoint
	borderPath.points[firstNewPoint.y][firstNewPoint.x].lastPos = lastPos

	return borderPath, count
}

// Determine interior map
func determineInteriorMap(border *border) [][]bool {
	// Create grid of the right size
	interiorMap := make([][]bool, border.height)
	for j := range interiorMap {
		interiorMap[j] = make([]bool, border.width)
	}

	// Print detail of point 5,0
	// fmt.Println()
	// fmt.Printf("Point: %v\n", border.points[157][6])
	// fmt.Printf("Point: %v\n", border.points[157][10])
	// fmt.Printf("Point: %v\n", border.points[6][0])
	// fmt.Printf("Point: %v\n", border.points[5][0])

	// Loop over border points
	for j := range border.points {
		insideCount := 0
		for i := range border.points[j] {
			// Check if point has a direction
			yDiff := border.points[j][i].nextPos.y - border.points[j][i].lastPos.y
			insideCount += yDiff
			// Debug
			// if j == 157 {
			// 	fmt.Println("i:", i, ", yDiff:", yDiff, ", insideCount:", insideCount)
			// }

			if yDiff == 0 && abs(insideCount) >= 2 {
				// We are inside the shape and not on a trench, so mark as interior
				interiorMap[j][i] = true
			}
		}
	}

	return interiorMap
}

// Trace path border
func tracePathBorderMap(trenches []trench) (borderMap, [][]uint64) {
	borderPath := borderMap{}

	// Find boundaries
	minPos, maxPos := findMaxSize(trenches)
	borderPath.width = maxPos.x - minPos.x + 1
	borderPath.height = maxPos.y - minPos.y + 1
	fmt.Println("minPos:", minPos, "maxPos:", maxPos)

	borderPath.start = coord{0 - minPos.x, 0 - minPos.y}
	fmt.Println("start coord:", borderPath.start, " size:", borderPath.width, "x", borderPath.height)

	// Create an array of rows, each containing a slice of encoded points on that row
	pointsOnRow := make([][]uint64, borderPath.height)

	// Create grid of the right size
	borderPath.points = make(map[uint64]borderPoint)

	pos := borderPath.start
	lastPos := pos
	firstNewPoint := coord{-1, -1}

	// // Add first point
	// borderPath.points[pos.y][pos.x] = borderPoint{
	// 	pos:       pos,
	// 	lastPos:   coord{-1, -1},
	// 	direction: trenches[0].direction,
	// 	colour:    trenches[0].color,
	// }
	// borderPath.points[encodeCoordAsUint64(borderPath.start)] = borderPoint{
	// 	pos:       borderPath.start,
	// 	lastPos:   lastPos,
	// 	direction: trenches[0].direction,
	// 	colour:    trenches[0].color,
	// }
	// pointsOnRow[borderPath.start.y] = append(pointsOnRow[borderPath.start.y], encodeCoordAsUint64(borderPath.start))
	// Loop over trenches
	for _, trench := range trenches {
		// Calculate new position
		// fmt.Println("trench:", trench)
		newPos := pos
		switch trench.direction {
		case north:
			newPos.y -= trench.distance
			// fmt.Println("N = pos.y:", pos.y, ", newPos.y:", newPos.y, ", trench.distance:", trench.distance)
			for j := pos.y - 1; j >= newPos.y; j-- {
				thisPos := coord{pos.x, j}
				borderPath.points[encodeCoordAsUint64(thisPos)] = borderPoint{
					pos:       thisPos,
					lastPos:   lastPos,
					direction: trench.direction,
					colour:    trench.color,
				}
				// Remember this point is on this row
				pointsOnRow[j] = append(pointsOnRow[j], encodeCoordAsUint64(thisPos))

				// Set nextPos of last point to this point
				updatePreviousPointsNextPos(&borderPath, lastPos, thisPos)

				// spew.Dump(thisPos)
				// fmt.Println(thisPos)
				if firstNewPoint.x == -1 {
					firstNewPoint = thisPos
				}

				lastPos = thisPos

			}
		case east:
			newPos.x += trench.distance
			// fmt.Println("pos.x:", pos.x, ", newPos.x:", newPos.x)
			for i := pos.x + 1; i <= newPos.x; i++ {
				// fmt.Println("i:", i)
				thisPos := coord{i, pos.y}
				borderPath.points[encodeCoordAsUint64(thisPos)] = borderPoint{
					pos:       thisPos,
					lastPos:   lastPos,
					direction: trench.direction,
					colour:    trench.color,
				}
				// Remember this point is on this row
				pointsOnRow[pos.y] = append(pointsOnRow[pos.y], encodeCoordAsUint64(thisPos))

				// Set nextPos of last point to this point
				updatePreviousPointsNextPos(&borderPath, lastPos, thisPos)

				// spew.Dump(thisPos)
				// fmt.Println(thisPos)
				if firstNewPoint.x == -1 {
					firstNewPoint = thisPos
				}

				lastPos = thisPos

			}
		case south:
			newPos.y += trench.distance
			for j := pos.y + 1; j <= newPos.y; j++ {
				thisPos := coord{pos.x, j}
				borderPath.points[encodeCoordAsUint64(thisPos)] = borderPoint{
					pos:       thisPos,
					lastPos:   lastPos,
					direction: trench.direction,
					colour:    trench.color,
				}
				// Remember this point is on this row
				pointsOnRow[j] = append(pointsOnRow[j], encodeCoordAsUint64(thisPos))

				// Set nextPos of last point to this point
				updatePreviousPointsNextPos(&borderPath, lastPos, thisPos)

				// spew.Dump(borderPath.points[j][pos.x])
				// fmt.Println(thisPos)
				if firstNewPoint.x == -1 {
					firstNewPoint = thisPos
				}

				lastPos = thisPos

			}
		case west:
			newPos.x -= trench.distance
			for i := pos.x - 1; i >= newPos.x; i-- {
				thisPos := coord{i, pos.y}
				borderPath.points[encodeCoordAsUint64(thisPos)] = borderPoint{
					pos:       thisPos,
					lastPos:   lastPos,
					direction: trench.direction,
					colour:    trench.color,
				}
				// Remember this point is on this row
				pointsOnRow[pos.y] = append(pointsOnRow[pos.y], encodeCoordAsUint64(thisPos))

				// Set nextPos of last point to this point
				updatePreviousPointsNextPos(&borderPath, lastPos, thisPos)

				// spew.Dump(borderPath.points[pos.y][i])
				// fmt.Println(thisPos)
				if firstNewPoint.x == -1 {
					firstNewPoint = thisPos
				}

				lastPos = thisPos

			}
		}
		// if t == 1 {
		// 	spew.Dump(borderPath.points)
		// }
		pos = newPos
	}
	// Connect last point to first point
	lastEncodedPoint := encodeCoordAsUint64(lastPos)
	lastPoint := borderPath.points[lastEncodedPoint]
	lastPoint.nextPos = firstNewPoint
	borderPath.points[lastEncodedPoint] = lastPoint
	firstEncodedPoint := encodeCoordAsUint64(firstNewPoint)
	firstPoint := borderPath.points[firstEncodedPoint]
	firstPoint.lastPos = lastPos
	borderPath.points[firstEncodedPoint] = firstPoint

	return borderPath, pointsOnRow
}

// Update previous points nextPos
func updatePreviousPointsNextPos(border *borderMap, lastPos, thisPos coord) {
	lastEncodedPoint := encodeCoordAsUint64(lastPos)
	lastPoint := border.points[lastEncodedPoint]
	lastPoint.nextPos = thisPos
	border.points[lastEncodedPoint] = lastPoint
}

func countInteriorMapPointsAndTrench(border *borderMap, pointsOnRow *[][]uint64) int {
	count := 0
	// fmt.Println()
	// Loop over each row
	for rowIndex := range *pointsOnRow {
		// Sort points on row
		rowPoints := (*pointsOnRow)[rowIndex]
		sort.Slice(rowPoints, func(i, j int) bool { return rowPoints[i] < rowPoints[j] })
		inStartX := -1
		insideCount := 0

		// Debug
		// if rowIndex == 0 {
		// fmt.Printf("%2d:", rowIndex)
		// for i := range rowPoints {
		// 	fmt.Printf(" %2d", rowPoints[i])
		// }
		// fmt.Println()
		// fmt.Printf("%2d:", rowIndex)
		// for i := range rowPoints {
		// 	fmt.Printf(" %2d", (*border).points[rowPoints[i]].pos.x)
		// }
		// fmt.Println()
		// }

		for i := 0; i < len(rowPoints); i++ {
			point := border.points[rowPoints[i]]
			yDiff := point.nextPos.y - point.lastPos.y
			insideCount += yDiff

			// The trench and inside counts.  So start counting if we are inside the trench
			// if abs(insideCount) >= 2 {
			if abs(insideCount) >= 1 {
				// We are inside the trench, so remember the start point if we don't have a start point set yet
				if inStartX == -1 {
					inStartX = point.pos.x
				}
			} else if inStartX != -1 {
				// We have left the trench or interior, so count the points between the start and end point
				// fmt.Println("rowIndex:", rowIndex, ", inStartX:", inStartX, ", point.pos.x:", point.pos.x, ", count:", point.pos.x-inStartX+1)
				count += point.pos.x - inStartX + 1
				inStartX = -1
			}
			// fmt.Printf("rowIndex: %2d, i: %2d, x: %2d nextx: %2d, lastx: %2d, y: %2d nexty: %2d, lasty: %2d, yDiff: %2d, insideCount: %2d, inStartX: %2d, count: %2d\n", rowIndex, i, point.pos.x, point.nextPos.x, point.lastPos.x, point.pos.y, point.nextPos.y, point.lastPos.y, yDiff, insideCount, inStartX, count)
		}
		// fmt.Printf("After row %2d, count: %2d\n", rowIndex, count)
	}

	// for j := range border.points {
	// 	// insideCount := 0
	// 	// for i := range border.points[j] {
	// 	// 	// Check if point has a direction
	// 	// 	yDiff := border.points[j][i].nextPos.y - border.points[j][i].lastPos.y
	// 	// 	insideCount += yDiff

	// 	// 	if yDiff == 0 && abs(insideCount) >= 2 {
	// 	// 		count++
	// 	// 	}
	// 	// }
	// }

	return count
}

// Print interior map
func printInteriorMap(interiorMap [][]bool) {
	fmt.Println("Interior map:")
	for _, row := range interiorMap {
		for _, point := range row {
			if point {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

// Print complete map
func printCompleteMap(border *border, interiorMap [][]bool) {
	fmt.Println("Complete map:")
	for j := range border.points {
		fmt.Printf("%03d ", j)
		for i := range border.points[j] {
			// Check if point has a direction
			if border.points[j][i].direction != -1 {
				fmt.Print(convertDirectionToString(border.points[j][i].direction))
			} else if interiorMap[j][i] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

// Count trench & interior points
func countTrenchAndInteriorPoints(border *border, interiorMap [][]bool) int {
	count := 0

	// Count all border points
	for j := range border.points {
		for i := range border.points[j] {
			// Check if point has a direction
			if border.points[j][i].direction != -1 {
				count++
			}
		}
	}

	// Count all elements on the interior map
	for j := range interiorMap {
		for i := range interiorMap[j] {
			// Check if point is on the interior map
			if interiorMap[j][i] {
				// Point is on the interior map, so count it
				count++
			}
		}
	}

	return count
}

// Decode hex as instruction
func decodeHexInstruction(hex string) (int, int) {
	// fmt.Printf("hex dir: %d %c\n", mod4(int(hex[len(hex)-1]-47)), rune(hex[len(hex)-1]))
	dir := mod4(int(hex[len(hex)-1] - 47))

	distHex := hex[:len(hex)-1]
	dist, err := strconv.ParseInt(strings.TrimLeft(distHex, "#"), 16, 64)
	if err != nil {
		panic(err)
	}

	return dir, int(dist)
}

// Function that calculate the mod 4 of an integer
func mod4(x int) int {
	return x & 3
}

// Abs function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Parse path
func parsePath(path string) []trench {
	trenches := make([]trench, 0)

	// Split path into lines and loop over them
	for _, line := range strings.Split(path, "\n") {
		// Split line into direction, distance and colour
		parts := strings.Split(line, " ")

		// Part 1
		// dir := convertDirection(parts[0])
		// dist, _ := strconv.Atoi(parts[1])

		// Part 2
		dir, dist := decodeHexInstruction(strings.Trim(parts[2], "()"))

		trenches = append(trenches, trench{
			direction: dir,
			distance:  dist,
			color:     strings.Trim(parts[2], "()"),
		})
	}

	return trenches
}

// Convert direction
func convertDirection(dir string) int {
	switch dir {
	case "R":
		return east
	case "D":
		return south
	case "L":
		return west
	case "U":
		return north
	}

	return -1
}

// Print trenches
func printTrenches(trenches []trench) {
	fmt.Println("Trenches:")
	for _, trench := range trenches {
		fmt.Printf("  %5s %02d %s\n", convertDirectionToString(trench.direction), trench.distance, trench.color)
	}
}

func printBorder(border border) {
	fmt.Println("Border:")
	for _, row := range border.points {
		for _, point := range row {
			if point.direction == -1 {
				fmt.Print(".")
			} else {
				fmt.Print(convertDirectionToString(point.direction))
			}
		}
		fmt.Println()
	}
}

// Convert direction to string
func convertDirectionToString(dir int) string {
	switch dir {
	case north:
		return "n"
	case east:
		return "e"
	case south:
		return "s"
	case west:
		return "w"
	}

	return ""
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
