package main

import (
	"fmt"
	"os"
)

// Race record structure
type RaceRecord struct {
	time     int
	distance int
}

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 06)

	// records := []RaceRecord{{7, 9}, {15, 40}, {30, 200}}
	// records := []RaceRecord{{48, 255}, {87, 1288}, {69, 1117}, {81, 1623}}
	// records := []RaceRecord{{71530, 940200}}
	records := []RaceRecord{{48876981, 255128811171623}}
	var marginOfError int = 0

	for _, record := range records {
		chargeTimes := findFurthestDistance(record.time, record.distance)
		if marginOfError > 0 {
			marginOfError *= len(chargeTimes)
		} else {
			marginOfError = len(chargeTimes)
		}
		// Print count of winning results
		fmt.Printf("Race time: %d, Distance: %d, Winning results: %d\n", record.time, record.distance, len(chargeTimes))
		// spew.Dump(chargeTimes)
	}
	fmt.Printf("Total Margin of error: %d\n", marginOfError)

}

// Find fastest distance for given time of race
func findFurthestDistance(raceTime int, existingRecord int) []int {
	var chargeTimes []int
	var maxDistance int = 0
	var lastDistance int = 0

	// Loop from 1 until maxDistance starts to increase again
	for i := 1; i < raceTime; i++ {
		lastDistance = (raceTime - i) * i
		// fmt.Printf("  Boost time: %d, Distance: %d, Existing Record: %d\n", i, lastDistance, existingRecord)
		if lastDistance > maxDistance {
			maxDistance = lastDistance
		}
		if lastDistance > existingRecord {
			chargeTimes = append(chargeTimes, i)
		}
	}

	return chargeTimes
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
