package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Structure for almanacMap
type almanacMap struct {
	source      string
	destination string
	entries     []almanacMapEntry
}

type almanacMapEntry struct {
	destinationStart int
	sourceStart      int
	length           int
}

var part1example = `seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 05)

	// Load input file
	inputContents := loadFileContents("almanac.txt")
	// inputContents := part1example

	seeds := parseSeeds(inputContents)
	// spew.Dump(seeds)

	// Parse input file
	almanac := parseAlmanac(inputContents)
	// spew.Dump(almanac)

	// mapSeedsToSoilPart1(seeds, almanac)
	mapSeedsToSoilPart2(seeds, almanac)

}

// Map seeds to soil with part 1 interpretation
func mapSeedsToSoilPart1(seeds []int, almanac map[string]almanacMap) {
	lowestLocation := math.MaxInt32
	// Loop over all seeds
	for _, seed := range seeds {
		location := processMap(seed, "seed", almanac)
		if location < lowestLocation {
			lowestLocation = location
		}
		fmt.Printf("Seed %d maps to Location %d\n\n", seed, location)
	}
	fmt.Printf("Lowest location is %d\n", lowestLocation)
}

// Map seeds to soil with part 2 interpretation
func mapSeedsToSoilPart2(seeds []int, almanac map[string]almanacMap) {
	lowestLocation := math.MaxInt32
	// Loop over all seeds
	for pos, _ := range seeds {
		if pos%2 == 0 {
			continue
		}
		fmt.Printf("Processing seeds %d to %d\n", seeds[pos-1], seeds[pos-1]+seeds[pos])
		for seed := seeds[pos-1]; seed < seeds[pos-1]+seeds[pos]; seed++ {
			location := processMap(seed, "seed", almanac)
			if location < lowestLocation {
				lowestLocation = location
			}
			// fmt.Printf("Seed %d maps to Location %d\n\n", seed, location)
		}
	}
	fmt.Printf("Lowest part 2 location is %d\n", lowestLocation)
}

// Process a map
func processMap(input int, source string, almanac map[string]almanacMap) int {
	// Get almanacMap from almanac
	almanacMap := almanac[source]
	mappedValue := -1

	// Loop over all entries in almanacMap
	for _, entry := range almanacMap.entries {
		// Check if input is in range of sourceStart and sourceStart + length
		if input >= entry.sourceStart && input < entry.sourceStart+entry.length {
			// Calculate destination
			mappedValue = entry.destinationStart + (input - entry.sourceStart)
		}
	}
	if mappedValue < 0 {
		mappedValue = input
	}
	// fmt.Printf("  mapped %s %d to %s %d\n", source, input, almanacMap.destination, mappedValue)

	if almanacMap.destination != "location" {
		mappedValue = processMap(mappedValue, almanacMap.destination, almanac)
	}

	return mappedValue
}

// Parse seeds from input string
func parseSeeds(input string) []int {
	seeds := []int{}

	// Loop over all lines in input, splitting on newline
	for _, line := range strings.Split(input, "\n") {
		// If line starts with the word seeds, split on space and take the second element
		if strings.HasPrefix(line, "seeds") {
			// Split on space
			splitLine := strings.Split(line, " ")
			for pos, seed := range splitLine {
				if pos > 0 {
					// Convert seed to int
					seedInt, err := strconv.Atoi(seed)
					if err != nil {
						panic(err)
					}
					// Add seed to seeds
					seeds = append(seeds, seedInt)
				}
			}
		}
	}

	return seeds
}

// Parse alamac input string
func parseAlmanac(input string) map[string]almanacMap {
	var currentMap almanacMap
	almanac := make(map[string]almanacMap)

	// Create a regex which checks if a line starts with a letter
	nameRegex := regexp.MustCompile(`^([a-zA-Z]+)-to-([a-zA-Z]+) map`)
	charRegex := regexp.MustCompile(`^[a-zA-Z]`)
	numRegex := regexp.MustCompile(`^[0-9]`)

	// Loop over all lines in input, splitting on newline
	for _, line := range strings.Split(input, "\n") {
		// If current line doesn't start with a number and currentEntry is not empty add entry to almanac
		if !numRegex.MatchString(line) && len(currentMap.entries) > 0 {
			// Add currentEntry to almanac
			almanac[currentMap.source] = currentMap
			// Reset currentEntry
			currentMap = almanacMap{}
		}

		// If current line starts with a letter and not the word seeds, start a new entry
		if charRegex.MatchString(line) && !strings.HasPrefix(line, "seeds") {
			// Set name of currentEntry
			nameMatch := nameRegex.FindStringSubmatch(line)
			currentMap.source = nameMatch[1]
			currentMap.destination = nameMatch[2]
		}

		// If current line starts with a number, split it by space and add it to currentEntry
		if numRegex.MatchString(line) {
			var err error
			entry := almanacMapEntry{}
			// Split line by space
			splitLine := strings.Split(line, " ")
			entry.destinationStart, err = strconv.Atoi(splitLine[0])
			if err != nil {
				panic(err)
			}
			entry.sourceStart, err = strconv.Atoi(splitLine[1])
			if err != nil {
				panic(err)
			}
			entry.length, err = strconv.Atoi(splitLine[2])
			if err != nil {
				panic(err)
			}
			// Add entry to currentMap
			currentMap.entries = append(currentMap.entries, entry)
		}
	}
	// If there is an active maps add it to the almanac
	if len(currentMap.entries) > 0 {
		// Add currentEntry to almanac
		almanac[currentMap.source] = currentMap
		// Reset currentEntry
		currentMap = almanacMap{}
	}

	return almanac
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
