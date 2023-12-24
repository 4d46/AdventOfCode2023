package main

import "testing"

func TestCanDrop(t *testing.T) {
	bricks := brickMap{
		1: brick{start: Point{0, 0, 1}, end: Point{0, 1, 1}},
		2: brick{start: Point{1, 0, 1}, end: Point{1, 1, 1}},
		3: brick{start: Point{0, 0, 2}, end: Point{1, 1, 2}},
	}

	bottoms := make(Layer, 3)
	bottoms[0] = map[int]int{}
	bottoms[1] = map[int]int{1: 1, 2: 2}
	bottoms[2] = map[int]int{3: 3}

	tops := make(Layer, 3)
	tops[0] = map[int]int{}
	tops[1] = map[int]int{1: 1, 2: 2}
	tops[2] = map[int]int{3: 3}

	// Test case 1: brick can fall to the bottom
	if !canDrop(1, 0, &bricks, &bottoms, &tops) {
		t.Errorf("Expected brick 1 to be droppable, but it wasn't")
	}

	// Test case 2: brick is supported by a brick in the lower layer
	if canDrop(3, 1, &bricks, &bottoms, &tops) {
		t.Errorf("Expected brick 3 to not be droppable, but it was")
	}

}

func TestCanDropCube(t *testing.T) {
	bricks := brickMap{
		1: brick{start: Point{0, 0, 1}, end: Point{1, 1, 2}},
		2: brick{start: Point{1, 1, 4}, end: Point{2, 2, 5}},
		3: brick{start: Point{0, 0, 6}, end: Point{1, 1, 7}},
	}

	bottoms := make(Layer, 8)
	bottoms[0] = map[int]int{}
	bottoms[1] = map[int]int{1: 1}
	bottoms[2] = map[int]int{}
	bottoms[3] = map[int]int{}
	bottoms[4] = map[int]int{2: 2}
	bottoms[5] = map[int]int{}
	bottoms[6] = map[int]int{3: 3}

	tops := make(Layer, 8)
	tops[0] = map[int]int{}
	tops[1] = map[int]int{1: 1, 2: 2}
	tops[2] = map[int]int{3: 3}

	// Test case 1: brick can fall to the bottom
	if !canDrop(1, 0, &bricks, &bottoms, &tops) {
		t.Errorf("Expected brick 1 to be droppable, but it wasn't")
	}

	// Test case 2: brick is supported by a brick in the lower layer
	if canDrop(3, 1, &bricks, &bottoms, &tops) {
		t.Errorf("Expected brick 3 to not be droppable, but it was")
	}

}
