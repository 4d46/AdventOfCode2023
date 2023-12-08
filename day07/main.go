package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/4d46/AdventOfCode2023/day07/card"
)

// Hand definition
type Hand struct {
	cards        [5]*card.Card
	bestHand     [5]*card.Card
	bestHandType card.HandType
	bid          int
	score        float64
	rawScore     float64
}

// Example hand part 1
var exampleHandString = `32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 07)

	// Load input
	// hands := parseHands(exampleHandString)
	hands := parseHands(loadFileContents("list_of_hands.txt"))

	// spew.Dump(hands)

	// Sort hands by score
	SortHands(hands)

	fmt.Printf("Hands: %d\n", len(hands))

	// spew.Dump(hands)
	printHands(hands)

	// Calculate Winnings
	winnings := calculateWinnings(hands)

	// Print Winnings
	fmt.Printf("Winnings: %d\n", winnings)
}

// Function that prints our hands
func printHands(hands []Hand) {
	// Print hands
	for _, h := range hands {
		fmt.Printf("%s\n", card.FormatHand(h.cards))
		// fmt.Printf("%s %s %d %.0f %.0f %d\n", card.FormatHand(h.cards), card.FormatHand(h.bestHand), h.bestHandType, h.score, h.rawScore, h.bid)
	}
}

// Sort hands by score
func SortHands(hands []Hand) {
	// Sort hands by score
	sort.Slice(hands, func(i, j int) bool {
		if hands[i].bestHandType == hands[j].bestHandType {
			return hands[i].rawScore < hands[j].rawScore
		}
		return hands[i].bestHandType < hands[j].bestHandType
		// if hands[i].score == hands[j].score {
		// 	return hands[i].rawScore < hands[j].rawScore
		// }
		// return hands[i].score < hands[j].score
	})
}

// Calculate Winnings
func calculateWinnings(hands []Hand) int {
	// Calculate winnings
	var winnings int
	for i, h := range hands {
		winnings += h.bid * (i + 1)
	}

	return winnings
}

// Parse input string into hands
func parseHands(input string) []Hand {
	// Split input into lines
	lines := strings.Split(input, "\n")

	// Create a slice of hands
	hands := make([]Hand, len(lines))

	// Parse each line into a hand
	for i, line := range lines {
		// Split line into cards and bid
		parts := strings.Split(line, " ")
		cardsStr := parts[0]
		bid, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}
		cards := parseCards(cardsStr)
		bestHand := card.CalculateBestHand(cards)
		fmt.Printf("%s Best Hand: %s\n", card.FormatHand(cards), card.FormatHand(bestHand))
		// spew.Dump(bestHand)
		// Create a hand
		hands[i] = Hand{
			cards:        cards,
			bestHand:     bestHand,
			bid:          bid,
			bestHandType: card.ClassifyHandType(bestHand),
			score:        card.CalculateScore(bestHand),
			rawScore:     card.CalculateScore(cards),
		}
	}

	return hands
}

// Parse cards from string
func parseCards(cards string) [5]*card.Card {
	// Create a slice of cards
	var cardsSlice [5]*card.Card

	// Parse each card
	for i, c := range cards {
		cardsSlice[i] = card.Parse(c)
	}

	return cardsSlice
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
