package card

import (
	"math"
	"strings"
)

type Card struct {
	value int
}

type HandType int

const NumCards = 13

// Part 1 Card List
// const (
// 	c2 = iota
// 	c3
// 	c4
// 	c5
// 	c6
// 	c7
// 	c8
// 	c9
// 	cT
// 	cJ
// 	cQ
// 	cK
// 	cA
// )

// Part 2 Card List
const (
	cJ = iota
	c2
	c3
	c4
	c5
	c6
	c7
	c8
	c9
	cT
	cQ
	cK
	cA
)

const (
	HighCard HandType = iota
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

// Parse a string into a card
func Parse(s rune) *Card {
	var c Card
	switch s {
	case '2':
		c.value = c2
	case '3':
		c.value = c3
	case '4':
		c.value = c4
	case '5':
		c.value = c5
	case '6':
		c.value = c6
	case '7':
		c.value = c7
	case '8':
		c.value = c8
	case '9':
		c.value = c9
	case 'T':
		c.value = cT
	case 'J':
		c.value = cJ
	case 'Q':
		c.value = cQ
	case 'K':
		c.value = cK
	case 'A':
		c.value = cA
	default:
		panic("Invalid card")
	}
	return &c
}

// Get the string of a card
func (c *Card) String() string {
	var s string
	switch c.value {
	case c2:
		s = "2"
	case c3:
		s = "3"
	case c4:
		s = "4"
	case c5:
		s = "5"
	case c6:
		s = "6"
	case c7:
		s = "7"
	case c8:
		s = "8"
	case c9:
		s = "9"
	case cT:
		s = "T"
	case cJ:
		s = "J"
	case cQ:
		s = "Q"
	case cK:
		s = "K"
	case cA:
		s = "A"
	default:
		panic("Invalid card")
	}
	return s
}

func (c *Card) Value() int {
	return c.value
}

// Calculate the best hand from a set of cards
// func CalculateBestHand(cards [5]*Card) [5]*Card {
// 	var bestHand [5]*Card
// 	var bestScore float64

// 	// Loop through all possible hand variations
// 	for _, hand := range allHandVariations(cards) {
// 		// fmt.Printf("  %s All Hands: %s\n", FormatHand(cards), FormatHand(hand))

// 		// Calculate score for hand
// 		score := CalculateScore(hand)
// 		// Check if score is better than best hand
// 		if score > bestScore {
// 			// If so, set best hand to current hand
// 			bestHand = hand
// 			bestScore = score
// 		}
// 	}

// 	return bestHand
// }

func FormatHand(cards [5]*Card) string {
	// Create a slice of cards
	var cardStr []string

	// Parse each card
	for _, c := range cards {
		cardStr = append(cardStr, c.String())
	}

	return strings.Join(cardStr, "")
}

func CalculateScore(cards [5]*Card) float64 {
	var score float64

	// spew.Dump(cards)

	// Calculate hand type

	// handType := ClassifyHandType(cards)
	// spew.Dump(cardTally)
	// fmt.Printf("Classified Hand Type: %d\n", handType)

	// Calculate score
	// Add score for each card
	handsize := len(cards)
	for pos, c := range cards {
		newscore := float64(c.Value()) * math.Pow(NumCards, float64(handsize-pos-1))
		score += newscore
		// fmt.Printf("Card: %d, Position: %d, Power: %d, Score: %f\n", c.Value(), pos, handsize-pos-1, newscore)
	}

	// **NOTE** This was used for part 1 but splitting for part 2
	// Add score for hand type
	// handscore := float64(handType) * math.Pow(float64(NumCards), float64(handsize))
	// // fmt.Printf("Hand Type: %d, Score: %f\n", handType, handscore)
	// score += handscore

	return score
}

// Calculate all possible hand types if the Joker is wild
// func allHandVariations(cards [5]*Card) [][5]*Card {
// 	var hands [][5]*Card
// 	var nonJokerCards map[*Card]bool = make(map[*Card]bool)
// 	var hasJoker bool

// 	// Add passed hand to hands
// 	hands = append(hands, cards)

// 	// Loop through cards adding unique cards to map
// 	for _, c := range cards {
// 		if c.Value() != cJ {
// 			nonJokerCards[c] = true
// 		} else {
// 			hasJoker = true
// 		}
// 	}

// 	// If there is no joker, return the original hand
// 	if !hasJoker {
// 		return hands
// 	}

// 	//
// 	if len(nonJokerCards) == 0 {
// 		nonJokerCards[Parse('A')] = true
// 	}

// 	// Loop over non-joker cards in map and create a new hand for each, replacing any jokers
// 	for c := range nonJokerCards {
// 		// Create a new hand
// 		var newHand [5]*Card
// 		// Loop through cards
// 		for pos, c2 := range cards {
// 			// If the card is a joker, replace it with the non-joker card
// 			if c2.Value() == cJ {
// 				newHand[pos] = c
// 			} else {
// 				newHand[pos] = c2
// 			}
// 		}
// 		// Add new hand to hands
// 		hands = append(hands, newHand)
// 	}

// 	return hands

// }

func ClassifyHandType(cards [5]*Card) HandType { // Classify hand
	var handType HandType

	// Add cards to map
	var cardTally = make(map[int]int)
	var jokers int
	// Loop through cards
	for _, c := range cards {
		if c.Value() == cJ {
			jokers++
		} else {
			// Add card tally to map
			cardTally[c.Value()]++
		}
	}

	// Find card with the highest count
	var mostCommonCard int
	for k, v := range cardTally {
		if mostCommonCard == 0 || v > cardTally[mostCommonCard] {
			mostCommonCard = k
		}
	}

	// Add jokers to most common card
	cardTally[mostCommonCard] += jokers

	// spew.Dump(cardTally)
	switch len(cardTally) {
	case 5:
		handType = HighCard
	case 4:
		handType = OnePair
	case 3:
		// Could be multiple options
		// Check for 3 of a kind
		for _, v := range cardTally {
			if v == 3 {
				handType = ThreeOfAKind
				break
			}
		}
		// Check for 2 pairs
		if handType == 0 {
			handType = TwoPair
		}
	case 2:
		// Could be multiple options
		// Check for 4 of a kind
		for _, v := range cardTally {
			if v == 4 {
				handType = FourOfAKind
				break
			}
		}
		// Check for full house
		if handType == 0 {
			handType = FullHouse
		}
	case 1:
		handType = FiveOfAKind
	default:
		panic("Invalid hand")
	}

	return handType
}

// func ClassifyHandType(cards [5]*Card) HandType { // Classify hand
// 	var handType HandType

// 	// Add cards to map
// 	var cardTally = make(map[int]int)
// 	// Loop through cards
// 	for _, c := range cards {
// 		// Add card tally to map
// 		cardTally[c.Value()]++
// 	}

// 	// spew.Dump(cardTally)
// 	switch len(cardTally) {
// 	case 5:
// 		handType = HighCard
// 	case 4:
// 		handType = OnePair
// 	case 3:
// 		// Could be multiple options
// 		// Check for 3 of a kind
// 		for _, v := range cardTally {
// 			if v == 3 {
// 				handType = ThreeOfAKind
// 				break
// 			}
// 		}
// 		// Check for 2 pairs
// 		if handType == 0 {
// 			handType = TwoPair
// 		}
// 	case 2:
// 		// Could be multiple options
// 		// Check for 4 of a kind
// 		for _, v := range cardTally {
// 			if v == 4 {
// 				handType = FourOfAKind
// 				break
// 			}
// 		}
// 		// Check for full house
// 		if handType == 0 {
// 			handType = FullHouse
// 		}
// 	case 1:
// 		handType = FiveOfAKind
// 	default:
// 		panic("Invalid hand")
// 	}

// 	return handType
// }
