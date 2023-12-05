package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type scratchcard struct {
	title          string
	winningNumbers []int
	numbers        []int
	copies         int
}

var part1example = `Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19
Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1
Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83
Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36
Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11`

var part2example = `Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19
Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1
Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83
Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36
Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 04)

	// inputStr := part1example
	// inputStr := part2example
	inputStr := loadFileContents("ScratchCardPart1.txt")
	// spew.Dump(inputStr)

	cards := parseScratchcards(inputStr)
	// spew.Dump(cards)

	var scores []int
	// Calculate the score of each card
	for _, card := range cards {
		score := card.score()
		scores = append(scores, score)
	}

	// add all the values in the scores list
	total := 0
	for _, score := range scores {
		total += score
	}
	fmt.Printf("Total part 1 score: %d\n", total)

	part2scoring(&cards)
	// spew.Dump(cards)
	// Count total number of card copies
	total = 0
	for _, card := range cards {
		total += card.copies
	}
	fmt.Printf("Total part 2 copies: %d\n", total)
}

// Calculate part 2 copies od scratchcards
func part2scoring(cards *[]scratchcard) {
	// fmt.Println("BEFORE Card 6 Copies: ", (*cards)[5].copies)

	// Calculate the score of each card
	for pos, card := range *cards {
		score := card.matches()
		// Add 1 to the number of copies of the next nnumber of cards matching the score
		end := pos + score + 1
		if end > len(*cards) {
			end = len(*cards)
		}
		// fmt.Println("Card: ", pos, " Score: ", score, " End: ", end)
		for i := pos + 1; i < end; i++ {
			(*cards)[i].copies += (*cards)[pos].copies
		}
		// fmt.Println("Card 6 Copies: ", (*cards)[5].copies)
	}
}

// Calculate the matches of a scratchcard
func (card *scratchcard) matches() int {
	score := 0
	for _, number := range card.numbers {
		if contains(card.winningNumbers, number) {
			score++
		}
	}
	return score
}

// Calculate the score of a scratchcard
func (card *scratchcard) score() int {
	score := 0
	for _, number := range card.numbers {
		if contains(card.winningNumbers, number) {
			score++
		}
	}
	if score > 0 {
		score = int(math.Pow(2, float64(score-1)))
	}
	return score
}

// Check if a list of numbers contains a specific number
func contains(numbers []int, number int) bool {
	for _, n := range numbers {
		if n == number {
			return true
		}
	}
	return false
}

// Parse multiline string into a set of scratchcards structures
// Each card has two lists of numbers separated by a vertical bar (|),
// a list of winning numbers and then a list of numbers in the follow format:
// Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
func parseScratchcards(input string) []scratchcard {
	var cards []scratchcard
	for _, line := range strings.Split(input, "\n") {
		var card scratchcard
		// Split the line into two parts, the card number and the numbers
		if len(line) == 0 {
			continue
		}
		titleContent := strings.Split(line, ":")
		card.title = titleContent[0]
		parts := strings.Split(titleContent[1], "|")
		card.winningNumbers = parseNumbers(parts[0])
		card.numbers = parseNumbers(parts[1])
		card.copies = 1
		cards = append(cards, card)
	}
	return cards
}

// Parse a string of numbers into a list of integers
func parseNumbers(input string) []int {
	var numbers []int
	for _, numberStr := range strings.Split(strings.Trim(input, " "), " ") {
		if len(numberStr) > 0 {
			number, err := strconv.Atoi(numberStr)
			if err != nil {
				panic(err)
			}
			numbers = append(numbers, number)
		}
	}
	return numbers
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
