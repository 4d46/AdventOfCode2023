package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Structure that represents a game.  Contains game ID and an array of 3 set structures
type Game struct {
	ID       int
	Sets     []Set
	MinCubes Set
}

// Structure that represents a set.  Contains integer values for red, green and blue
type Set struct {
	Red   int
	Green int
	Blue  int
}

var exampleRecord string = `Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green`

var maxCubes Set = Set{Red: 12, Green: 13, Blue: 14}

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 2)

	recordStr := loadFileContents("RecordInput.txt")

	// record := parseRecord(exampleRecord)
	record := parseRecord(recordStr)
	// spew.Dump(record)

	// printGames(record)

	fmt.Println()

	invalidGameIDsSum := sumGameIDsUptoMaxCubes(record)
	fmt.Printf("[Part 1] Sum of invalid game IDs: %d\n", invalidGameIDsSum)

	// Display the power required for all games
	power := calculatePowerForAllGames(record)
	fmt.Printf("[Part 2] Power required for all games: %d\n", power)

	// printRecord(record)
}

// Parse record string into an array of games
func parseRecord(record string) []Game {
	// Create an array of games
	var games []Game

	// Loop over each line parsing a game and adding it to the games array
	for _, line := range strings.Split(record, "\n") {
		if len(line) > 0 {
			games = append(games, parseGame(line))
		}
	}

	return games
}

// Function that parses a Game string and creates a game structure
// Example game string:
// Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
func parseGame(gameString string) Game {
	// Create a game structure
	var game Game

	// Parse the game string
	// Split the string into 2 parts to separate game id from the sets
	gameDetails := strings.Split(gameString, ":")
	// Parse the game id using a regurlar expresssion
	// Parse the game id using a regular expression
	re := regexp.MustCompile(`\d+`)
	gameIDStr := re.FindString(gameDetails[0])
	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		// handle error
		panic(err)
	}
	game.ID = gameID

	// Parse the sets
	// Split the sets string into an array of set strings
	setStrings := strings.Split(gameDetails[1], ";")
	// Initialize the sets array
	game.Sets = make([]Set, len(setStrings))

	// Loop over the set strings
	for i, setString := range setStrings {
		// Parse the set string
		// Split the set string into an array of color strings
		colorStrings := strings.Split(setString, ",")
		// Loop over the color strings
		for _, colorString := range colorStrings {
			// Parse the color string
			// Split the color string into an array of color and count strings
			colorCountStrings := strings.Split(strings.Trim(colorString, " "), " ")
			// spew.Dump(colorCountStrings[1])
			// Switch on the color string
			switch colorCountStrings[1] {
			case "red":
				// Parse the red count string
				redCount, err := strconv.Atoi(colorCountStrings[0])
				if err != nil {
					// handle error
					panic(err)
				}
				// Set the red count in the set structure
				game.Sets[i].Red = redCount
				// Remember the largest number of cubes of each color in a single game
				game.MinCubes.Red = maxInt(game.MinCubes.Red, redCount)
			case "green":
				// Parse the green count string
				greenCount, err := strconv.Atoi(colorCountStrings[0])
				if err != nil {
					// handle error
					panic(err)
				}
				// Set the green count in the set structure
				game.Sets[i].Green = greenCount
				// Remember the largest number of cubes of each color in a single game
				game.MinCubes.Green = maxInt(game.MinCubes.Green, greenCount)
			case "blue":
				// Parse the blue count string
				blueCount, err := strconv.Atoi(colorCountStrings[0])
				if err != nil {
					// handle error
					panic(err)
				}
				// Set the blue count in the set structure
				game.Sets[i].Blue = blueCount
				// Remember the largest number of cubes of each color in a single game
				game.MinCubes.Blue = maxInt(game.MinCubes.Blue, blueCount)
			}
		}
	}

	return game
}

// Function that sums the game ids that don't exceed the max cubes
func sumGameIDsUptoMaxCubes(games []Game) int {
	// Initialize the sum
	var sum int = 0

	// Loop over the games
	for _, game := range games {
		validGame := true
		for _, set := range game.Sets {
			// Check if the game exceeds the max cubes
			if set.Red > maxCubes.Red || set.Green > maxCubes.Green || set.Blue > maxCubes.Blue {
				validGame = false
			}
		}
		if validGame {
			fmt.Printf("Game %d possible\n", game.ID)
			// Add the game id to the sum
			sum += game.ID
		} else {
			fmt.Printf("- Game %d exceeds max cubes\n", game.ID)
		}
	}

	return sum
}

// Function that calculates the power required for all games
func calculatePowerForAllGames(games []Game) int {
	// Initialize the power
	var power int = 0

	// Loop over the games
	for _, game := range games {
		// Add the power for the game to the total power
		power += calculatePowerForGame(game)
	}

	return power
}

// Function that calculates the power required for a game
func calculatePowerForGame(game Game) int {
	// Initialize the power
	var power int = 0

	// Power is calculated as the product of the minimum number of each colour required for the game
	power = game.MinCubes.Red * game.MinCubes.Green * game.MinCubes.Blue

	return power
}

// Function that pretty prints an array of games
func printGames(games []Game) {
	for _, game := range games {
		printGame(game)
	}
}

// Function that pretty prints a game
func printGame(game Game) {
	fmt.Printf("Game %d\n", game.ID)
	for _, set := range game.Sets {
		printSet(set)
	}
	printTotals(game)
}

// Function that pretty prints the totals for a game
func printTotals(game Game) {
	fmt.Printf("  Min Cube Required: ")
	printSet(game.MinCubes)
}

// Function that pretty prints a set
func printSet(set Set) {
	fmt.Printf("  Red: %2d, Green: %2d, Blue: %2d\n", set.Red, set.Green, set.Blue)
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

// Write a function that outputs all record details in the original format
func printRecord(record []Game) {
	for _, game := range record {
		fmt.Printf("Game %d:", game.ID)
		firstSet := true
		for _, set := range game.Sets {
			if !firstSet {
				fmt.Printf(";")
			}
			firstColour := true

			if set.Red > 0 {
				if !firstColour {
					fmt.Printf(",")
				}
				firstColour = false
				fmt.Printf(" %d red", set.Red)
			}
			if set.Green > 0 {
				if !firstColour {
					fmt.Printf(",")
				}
				firstColour = false
				fmt.Printf(" %d green", set.Green)
			}
			if set.Blue > 0 {
				if !firstColour {
					fmt.Printf(",")
				}
				firstColour = false
				fmt.Printf(" %d blue", set.Blue)
			}
			firstSet = false
		}
		fmt.Printf("\n")
	}
}

// Function that returns the biggest integer of passed parameters
func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
