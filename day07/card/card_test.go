package card

import (
	"math"
	"testing"
)

func TestCalculateScoreLow(t *testing.T) {
	// Create test cards
	cards := [5]*Card{
		Parse('2'),
		Parse('3'),
		Parse('4'),
		Parse('5'),
		Parse('6'),
	}

	// spew.Dump(Parse('2'))

	// Calculate score
	score := CalculateScore(cards)

	// Define expected score for Part 1
	// expectedScore := 0*math.Pow(13, 4) + 1*math.Pow(13, 3) + 2*math.Pow(13, 2) + 3*math.Pow(13, 1) + 4*math.Pow(13, 0) + 0*math.Pow(13, 5)
	// Define expected score for Part 2
	expectedScore := 1*math.Pow(13, 4) + 2*math.Pow(13, 3) + 3*math.Pow(13, 2) + 4*math.Pow(13, 1) + 5*math.Pow(13, 0) + 0*math.Pow(13, 5)

	// Check if the calculated score matches the expected score
	if score != expectedScore {
		t.Errorf("Expected score: %f, but got: %f", expectedScore, score)
	}
}

func TestCalculateScoreAceHigh(t *testing.T) {
	// Create test cards
	cards := [5]*Card{
		Parse('A'),
		Parse('2'),
		Parse('3'),
		Parse('4'),
		Parse('5'),
	}

	// spew.Dump(Parse('2'))

	// Calculate score
	score := CalculateScore(cards)

	// Define expected score for part 1
	// expectedScore := 12*math.Pow(13, 4) + 0*math.Pow(13, 3) + 1*math.Pow(13, 2) + 2*math.Pow(13, 1) + 3*math.Pow(13, 0) + 0*math.Pow(13, 5)
	// Define expected score for part 2
	expectedScore := 12*math.Pow(13, 4) + 1*math.Pow(13, 3) + 2*math.Pow(13, 2) + 3*math.Pow(13, 1) + 4*math.Pow(13, 0) + 0*math.Pow(13, 5)

	// Check if the calculated score matches the expected score
	if score != expectedScore {
		t.Errorf("Expected score: %f, but got: %f", expectedScore, score)
	}
}

func TestCalculateScoreJackHigh(t *testing.T) {
	// Create test cards
	cards := [5]*Card{
		Parse('J'),
		Parse('2'),
		Parse('3'),
		Parse('4'),
		Parse('5'),
	}

	// spew.Dump(Parse('2'))

	// Calculate score
	score := CalculateScore(cards)

	// Define expected score for part 1
	// expectedScore := 9*math.Pow(13, 4) + 0*math.Pow(13, 3) + 1*math.Pow(13, 2) + 2*math.Pow(13, 1) + 3*math.Pow(13, 0) + 0*math.Pow(13, 5)
	// Define expected score for part 2
	expectedScore := 0*math.Pow(13, 4) + 1*math.Pow(13, 3) + 2*math.Pow(13, 2) + 3*math.Pow(13, 1) + 4*math.Pow(13, 0) + 0*math.Pow(13, 5)

	// Check if the calculated score matches the expected score
	if score != expectedScore {
		t.Errorf("Expected score: %f, but got: %f", expectedScore, score)
	}
}

func TestCalculateScoreTwo5Card(t *testing.T) {
	// Create test cards
	cards := [5]*Card{
		Parse('2'),
		Parse('2'),
		Parse('2'),
		Parse('2'),
		Parse('2'),
	}

	// spew.Dump(Parse('2'))

	// Calculate score
	score := CalculateScore(cards)

	// Define expected score for part 1
	// expectedScore := 0*math.Pow(13, 4) + 0*math.Pow(13, 3) + 0*math.Pow(13, 2) + 0*math.Pow(13, 1) + 0*math.Pow(13, 0) + 6*math.Pow(13, 5)
	// Define expected score for part 2
	expectedScore := 1*math.Pow(13, 4) + 1*math.Pow(13, 3) + 1*math.Pow(13, 2) + 1*math.Pow(13, 1) + 1*math.Pow(13, 0) + 6*math.Pow(13, 5)

	// Check if the calculated score matches the expected score
	if score != expectedScore {
		t.Errorf("Expected score: %f, but got: %f", expectedScore, score)
	}
}

func TestCalculateScoreAce5Card(t *testing.T) {
	// Create test cards
	cards := [5]*Card{
		Parse('A'),
		Parse('A'),
		Parse('A'),
		Parse('A'),
		Parse('A'),
	}

	// spew.Dump(Parse('2'))

	// Calculate score
	score := CalculateScore(cards)

	// Define expected score
	expectedScore := 12*math.Pow(13, 4) + 12*math.Pow(13, 3) + 12*math.Pow(13, 2) + 12*math.Pow(13, 1) + 12*math.Pow(13, 0) + 6*math.Pow(13, 5)

	// Check if the calculated score matches the expected score
	if score != expectedScore {
		t.Errorf("Expected score: %f, but got: %f", expectedScore, score)
	}
}

// 27J4K 27K4K 414945 41793 969
func TestCalculateScoreKingPairCard(t *testing.T) {
	// Create test cards
	cards := [5]*Card{
		Parse('2'),
		Parse('7'),
		Parse('K'),
		Parse('4'),
		Parse('K'),
	}

	// spew.Dump(Parse('2'))

	// Calculate score
	score := CalculateScore(cards)

	// Define expected score
	expectedScore := 1*math.Pow(13, 4) + 6*math.Pow(13, 3) + 11*math.Pow(13, 2) + 3*math.Pow(13, 1) + 11*math.Pow(13, 0) + 1*math.Pow(13, 5)

	// Check if the calculated score matches the expected score
	if score != expectedScore {
		t.Errorf("Expected score: %f, but got: %f", expectedScore, score)
	}
}

func TestCalculateScoreKingTenJackCard(t *testing.T) {
	// Create test cards
	cards := [5]*Card{
		Parse('K'),
		Parse('T'),
		Parse('J'),
		Parse('J'),
		Parse('T'),
	}

	// spew.Dump(Parse('2'))

	// Calculate score
	score := CalculateScore(cards)

	// Define expected score
	expectedScore := 11*math.Pow(13, 4) + 9*math.Pow(13, 3) + 0*math.Pow(13, 2) + 0*math.Pow(13, 1) + 9*math.Pow(13, 0) + 2*math.Pow(13, 5)

	// Check if the calculated score matches the expected score
	if score != expectedScore {
		t.Errorf("Expected score: %f, but got: %f", expectedScore, score)
	}
}

func TestCalculateBestHandOnlyJackWildcard(t *testing.T) {
	// Create test cards
	cards := [5]*Card{
		Parse('J'),
		Parse('J'),
		Parse('J'),
		Parse('J'),
		Parse('J'),
	}

	cardsExpected := [5]*Card{
		Parse('A'),
		Parse('A'),
		Parse('A'),
		Parse('A'),
		Parse('A'),
	}

	bestHand := CalculateBestHand(cards)

	matching := true
	for pos, _ := range bestHand {
		if bestHand[pos].value != cardsExpected[pos].value {
			matching = false
		}
	}
	// Check if the calculated score matches the expected score
	if !matching {
		t.Errorf("Expected: %v, but got: %v", cardsExpected, cards)
	}
}

func TestCalculateBestHandAceJackWildcard(t *testing.T) {
	// Create test cards
	cards := [5]*Card{
		Parse('A'),
		Parse('A'),
		Parse('J'),
		Parse('2'),
		Parse('3'),
	}

	cardsExpected := [5]*Card{
		Parse('A'),
		Parse('A'),
		Parse('A'),
		Parse('2'),
		Parse('3'),
	}

	bestHand := CalculateBestHand(cards)

	matching := true
	for pos, _ := range bestHand {
		if bestHand[pos].value != cardsExpected[pos].value {
			matching = false
		}
	}
	// Check if the calculated score matches the expected score
	if !matching {
		t.Errorf("Expected: %v, but got: %v", cardsExpected, cards)
	}
}

func TestCalculateBestHandFourJackWildcard(t *testing.T) {
	// Create test cards
	cards := [5]*Card{
		Parse('4'),
		Parse('4'),
		Parse('J'),
		Parse('2'),
		Parse('J'),
	}

	cardsExpected := [5]*Card{
		Parse('4'),
		Parse('4'),
		Parse('4'),
		Parse('2'),
		Parse('4'),
	}

	bestHand := CalculateBestHand(cards)

	matching := true
	for pos, _ := range bestHand {
		if bestHand[pos].value != cardsExpected[pos].value {
			matching = false
		}
	}
	// Check if the calculated score matches the expected score
	if !matching {
		t.Errorf("Expected: %v, but got: %v", cardsExpected, cards)
	}
}

func TestCalculateBestHandKingJackWildcard(t *testing.T) {
	// Create test cards
	cards := [5]*Card{
		Parse('2'),
		Parse('7'),
		Parse('J'),
		Parse('4'),
		Parse('K'),
	}

	cardsExpected := [5]*Card{
		Parse('2'),
		Parse('7'),
		Parse('K'),
		Parse('4'),
		Parse('K'),
	}

	bestHand := CalculateBestHand(cards)

	matching := true
	for pos, _ := range bestHand {
		if bestHand[pos].value != cardsExpected[pos].value {
			matching = false
		}
	}
	// Check if the calculated score matches the expected score
	if !matching {
		t.Errorf("Expected: %v, but got: %v", cardsExpected, cards)
	}
}

func TestClassifyHandTypeHighCard(t *testing.T) {
	cards := [5]*Card{
		Parse('2'),
		Parse('3'),
		Parse('4'),
		Parse('5'),
		Parse('6'),
	}

	expectedHandType := HighCard
	handType := ClassifyHandType(cards)

	if handType != expectedHandType {
		t.Errorf("Expected hand type: %v, but got: %v", expectedHandType, handType)
	}
}

func TestClassifyHandTypeOnePair(t *testing.T) {
	cards := [5]*Card{
		Parse('A'),
		Parse('2'),
		Parse('3'),
		Parse('A'),
		Parse('4'),
	}

	expectedHandType := OnePair
	handType := ClassifyHandType(cards)

	if handType != expectedHandType {
		t.Errorf("Expected hand type: %v, but got: %v", expectedHandType, handType)
	}
}

func TestClassifyHandTypeThreeOfAKind(t *testing.T) {
	cards := [5]*Card{
		Parse('T'),
		Parse('T'),
		Parse('T'),
		Parse('9'),
		Parse('8'),
	}

	expectedHandType := ThreeOfAKind
	handType := ClassifyHandType(cards)

	if handType != expectedHandType {
		t.Errorf("Expected hand type: %v, but got: %v", expectedHandType, handType)
	}
}

func TestClassifyHandTypeTwoPair(t *testing.T) {
	cards := [5]*Card{
		Parse('2'),
		Parse('3'),
		Parse('4'),
		Parse('3'),
		Parse('2'),
	}

	expectedHandType := TwoPair
	handType := ClassifyHandType(cards)

	if handType != expectedHandType {
		t.Errorf("Expected hand type: %v, but got: %v", expectedHandType, handType)
	}
}

func TestClassifyHandTypeFourOfAKind(t *testing.T) {
	cards := [5]*Card{
		Parse('A'),
		Parse('A'),
		Parse('8'),
		Parse('A'),
		Parse('A'),
	}

	expectedHandType := FourOfAKind
	handType := ClassifyHandType(cards)

	if handType != expectedHandType {
		t.Errorf("Expected hand type: %v, but got: %v", expectedHandType, handType)
	}
}

func TestClassifyHandTypeFullHouse(t *testing.T) {
	cards := [5]*Card{
		Parse('2'),
		Parse('3'),
		Parse('3'),
		Parse('3'),
		Parse('2'),
	}

	expectedHandType := FullHouse
	handType := ClassifyHandType(cards)

	if handType != expectedHandType {
		t.Errorf("Expected hand type: %v, but got: %v", expectedHandType, handType)
	}
}

func TestClassifyHandTypeFiveOfAKind(t *testing.T) {
	cards := [5]*Card{
		Parse('A'),
		Parse('A'),
		Parse('A'),
		Parse('A'),
		Parse('A'),
	}

	expectedHandType := FiveOfAKind
	handType := ClassifyHandType(cards)

	if handType != expectedHandType {
		t.Errorf("Expected hand type: %v, but got: %v", expectedHandType, handType)
	}
}
