package main

import "testing"

func TestDecodeHexInstruction(t *testing.T) {
	tests := []struct {
		hex       string
		direction int
		distance  int
	}{
		{"#70c710", 1, 461937},
		{"#0dc571", 2, 56407},
		{"#1b58a2", 3, 112010},
		{"#7a21e3", 0, 500254},
	}

	for _, test := range tests {
		direction, distance := decodeHexInstruction(test.hex)
		if direction != test.direction || distance != test.distance {
			t.Errorf("decodeHexInstruction(%s) = (%d, %d), want (%d, %d)", test.hex, direction, distance, test.direction, test.distance)
		}
	}
}
