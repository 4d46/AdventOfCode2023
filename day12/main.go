package main

import (
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// Spring map entry
type springMapEntry struct {
	springMap   string
	springList  []int
	numUnkonwns int
	// arrangements []string
	arrangements int
}

const springExample1 = `???.### 1,1,3
.??..??...?##. 1,1,3
?#?#?#?#?#?#?#? 1,3,1,6
????.#...#... 4,1,1
????.######..#####. 1,6,5
?###???????? 3,2,1`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 12)

	test := fmt.Sprintf("%v", []int{1, 2, 3, 4}[1:])
	fmt.Println(test)

	// springMapStr := springExample1
	springMapStr := loadFileContents("damaged_springs.txt")

	fmt.Print("Building Truth table: ")
	ttable := buildSpringMapTruthTable()
	fmt.Println("DONE")
	spew.Dump(ttable[3])

	// Parse the spring map
	springMap := parseSpringMap(springMapStr)
	fmt.Println("Before:")
	printSpringMap(springMap, false)

	// Adding ufolding for part 2
	for i, _ := range springMap {
		springMap[i] = unfoldSpringMap(springMap[i])
	}
	printSpringMap(springMap, false)

	fmt.Println("\nAfter:")
	// // Calculate possible spring map arrangements for all spring maps
	// for i, _ := range springMap {

	// 	springMap[i].arrangements = calcPossibleSpringMapCombinations(springMap[i])
	// 	printSpringMapEntry(springMap[i])
	// }
	// printSpringMap(springMap, true)

	// Calculate sum of possible arrangements
	// sum := 0
	// for _, springMapEntry := range springMap {
	// 	// sum += len(springMapEntry.arrangements)
	// 	sum += springMapEntry.arrangements
	// 	panic("Not implemented")
	// }

	// Use new recursive function
	sum := 0
	for i := range springMap {
		sum += computePossibleSpringMapCombinationsRecursive2(&springMap[i], &ttable, false)
		fmt.Printf("Spring Map %2d/%d\n", i, len(springMap))
	}
	fmt.Printf("\nSum of possible arrangements: %d\n", sum)

	// wg := sync.WaitGroup{}
	// for i, sme := range springMap {
	// 	wg.Add(1)
	// 	go func(sme *springMapEntry, ttable *[15]map[string]bool, print bool, i int) {
	// 		sum := computePossibleSpringMapCombinationsRecursive(sme, ttable, false)
	// 		fmt.Printf("Spring Map %2d/%d\n", i, len(springMap))
	// 		wg.Done()
	// 	}(&sme, &ttable, false, i)
	// }
	// wg.Wait()
	// fmt.Printf("\nSum of possible arrangements: %d\n", sum)
}

// Function to unfold the spring map
func unfoldSpringMap(originalMap springMapEntry) springMapEntry {
	var returnSpringMap springMapEntry

	// Calculate the new spring map
	returnSpringMap.springMap = originalMap.springMap
	for i := 0; i < 4; i++ {
		returnSpringMap.springMap += "?" + originalMap.springMap
	}

	for i := 0; i < 5; i++ {
		returnSpringMap.springList = append(returnSpringMap.springList, originalMap.springList...)
	}

	returnSpringMap.numUnkonwns = countCharacters(returnSpringMap.springMap, '?')

	return returnSpringMap
}

// Calculate number of possible spring map combinations
func calcPossibleSpringMapCombinations(springMap springMapEntry) int {
	// Create return variable
	var possibleSpringMapCombinations int

	// Calculate number of broken springs
	var numUnkownBrokenSprings int
	spew.Dump(springMap)
	for _, springListEntry := range springMap.springList {
		numUnkownBrokenSprings += springListEntry
	}
	// Remove the number of known broken springs from the list
	numUnkownBrokenSprings -= countCharacters(springMap.springMap, '#')
	fmt.Println("numUnkownBrokenSprings:", numUnkownBrokenSprings)

	// Check if there are unknowns, if not just return 1 for the number of solutions
	if springMap.numUnkonwns < 1 {
		return 1
	}

	// Calculate the number of possible combinations
	numPossibleCombinations := uint64(1 << springMap.numUnkonwns)
	fmt.Printf("numPossibleCombinations: %d\n", numPossibleCombinations)

	// Loop over all possible combinations
	var i uint64
	for i = 0; i < numPossibleCombinations; i++ {

		// Check if the number of broken springs matches the number in the list, otherwise skip
		// fmt.Println("binary characters:", countCharacters(binaryString, '1'))
		if bits.OnesCount64(i) != numUnkownBrokenSprings {
			continue
		}
		// fmt.Printf("i: %d\n", i)
		mergedSpringMap := mergeSpringMap(springMap.springMap, i)
		// fmt.Printf("mergedSpringMap: %s\n", mergedSpringMap)
		// mergedSpringMapList := classifySpringMap(mergedSpringMap)
		mergedSpringMapList := classifySpringMapUint64(mergedSpringMap)
		if compareIntArrays(
			&mergedSpringMapList, &(springMap.springList)) {
			possibleSpringMapCombinations++
		}
	}
	return possibleSpringMapCombinations
}

// Start recursive SpringMap Combinations compute
func computePossibleSpringMapCombinationsRecursive2(springMap *springMapEntry, ttable *[15]map[string]bool, print bool) int {
	// Lookup
	lookup := make(map[string]int)
	// Start recursive function at begining of string and section list, with a gap because it is at the start of the string
	result := computePossibleSpringMapCombinations2(springMap, ttable, &lookup, print, 0, 0, 1)
	return result
}

// Calculate number of possible spring map combinations
func computePossibleSpringMapCombinations2(springMap *springMapEntry, ttable *[15]map[string]bool, lookup *map[string]int, print bool, pos int, listPos int, gap int) int {
	// Serialise the remaining ask
	ask := fmt.Sprintf("%s%v", springMap.springMap[pos:], springMap.springList[listPos:])
	// Check if we have already calculated this and return the result
	if val, ok := (*lookup)[ask]; ok {
		return val
	}

	if print {
		if pos >= len(springMap.springMap) {
			fmt.Printf("string: %-16s, inf: %-16s, pos: %d, listPos: %d, gap: %d\n", springMap.springMap[:pos]+"|", pos, listPos, gap)
		} else {
			fmt.Printf("string: %-16s, inf: %-16s, pos: %d, listPos: %d, gap: %d\n", springMap.springMap[:pos+1], pos, listPos, gap)
		}
	}
	// Check if we have passed the end of the string
	if pos >= len(springMap.springMap) {
		// Check if we have coded all blocks
		if listPos >= len(springMap.springList) {
			// No more segments, found a valid combination
			if print {
				fmt.Println("          ^^^^^^^^")
			}
			return 1
		} else {
			// More segments, but we've passed the end of the string, not valid
			return 0
		}
	} else if listPos >= len(springMap.springList) {
		// All segments accounted for, check if we have broken strings left in string
		if strings.ContainsRune(springMap.springMap[pos:], '#') {
			// Broken strings left in string, not valid
			return 0
		} else {
			if print {
				fmt.Println("          ^^^^^^^^")
			}
			return 1
		}
	}

	// Still valid, check next character
	if springMap.springMap[pos] == '.' {
		// If we have a . in the string next segment isn't starting, return the result of the next position, indicating gap has increased
		return computePossibleSpringMapCombinations2(springMap, ttable, lookup, print, pos+1, listPos, gap+1)
	} else if springMap.springMap[pos] == '#' {
		// If there isn't a gap, this isn't valid, return 0
		if gap < 1 {
			return 0
		}
		// If there is a gap, confirm if the next segment is valid and doesn't exceed the end of the string
		sectionLen := springMap.springList[listPos]
		// if print {
		// 	fmt.Printf("sectionLen:%d  %d/%d [%s]\n", pos, sectionLen, len(springMap.springMap), springMap.springMap[pos:pos+sectionLen])
		// }
		if pos+sectionLen <= len(springMap.springMap) {
			// If it fits in the string, check if it is a valid segment
			if _, match := ttable[sectionLen-1][springMap.springMap[pos:pos+sectionLen]]; match {
				// If the next segment is valid, return the result of the next position, increaseing sections accounted for and resetting gap
				return computePossibleSpringMapCombinations2(springMap, ttable, lookup, print, pos+sectionLen, listPos+1, 0)
			}
		}
		// Will only get here If the next segment isn't valid, return 0
		return 0
	} else if springMap.springMap[pos] == '?' {
		var sum int

		// Here we have options
		// If gap is 0, we can't start a new segment, so this must be a gap
		if gap == 0 {
			return computePossibleSpringMapCombinations2(springMap, ttable, lookup, print, pos+1, listPos, gap+1)
		}
		// If gap is > 0, we can start a new segment or add another gap, try both
		sum += computePossibleSpringMapCombinations2(springMap, ttable, lookup, print, pos+1, listPos, gap+1)
		sectionLen := springMap.springList[listPos]
		// Check this won't exceed the end of the string
		if pos+sectionLen <= len(springMap.springMap) {
			// If it fits in the string, check if it is a valid segment
			if _, match := ttable[sectionLen-1][springMap.springMap[pos:pos+sectionLen]]; match {
				sum += computePossibleSpringMapCombinations2(springMap, ttable, lookup, print, pos+sectionLen, listPos+1, 0)
			}
		}
		// Store the result in the lookup
		(*lookup)[ask] = sum
		return sum
	} else {
		panic("Invalid character in spring map")
	}
}

// Start recursive SpringMap Combinations compute
func computePossibleSpringMapCombinationsRecursive(springMap *springMapEntry, ttable *[15]map[string]bool, print bool) int {
	// Start recursive function at begining of string and section list, with a gap because it is at the start of the string
	result := computePossibleSpringMapCombinations(springMap, ttable, print, "", 0, 0, 1)
	return result
}

// Calculate number of possible spring map combinations
func computePossibleSpringMapCombinations(springMap *springMapEntry, ttable *[15]map[string]bool, print bool, inferredStr string, pos int, listPos int, gap int) int {
	if print {
		if pos >= len(springMap.springMap) {
			fmt.Printf("string: %-16s, inf: %-16s, pos: %d, listPos: %d, gap: %d\n", springMap.springMap[:pos]+"|", inferredStr, pos, listPos, gap)
		} else {
			fmt.Printf("string: %-16s, inf: %-16s, pos: %d, listPos: %d, gap: %d\n", springMap.springMap[:pos+1], inferredStr, pos, listPos, gap)
		}
	}
	// Check if we have passed the end of the string
	if pos >= len(springMap.springMap) {
		// Check if we have coded all blocks
		if listPos >= len(springMap.springList) {
			// No more segments, found a valid combination
			if print {
				fmt.Println("          ^^^^^^^^")
			}
			return 1
		} else {
			// More segments, but we've passed the end of the string, not valid
			return 0
		}
	} else if listPos >= len(springMap.springList) {
		// All segments accounted for, check if we have broken strings left in string
		if strings.ContainsRune(springMap.springMap[pos:], '#') {
			// Broken strings left in string, not valid
			return 0
		} else {
			if print {
				fmt.Println("          ^^^^^^^^")
			}
			return 1
		}
	}

	// Still valid, check next character
	if springMap.springMap[pos] == '.' {
		// If we have a . in the string next segment isn't starting, return the result of the next position, indicating gap has increased
		return computePossibleSpringMapCombinations(springMap, ttable, print, inferredStr+".", pos+1, listPos, gap+1)
	} else if springMap.springMap[pos] == '#' {
		// If there isn't a gap, this isn't valid, return 0
		if gap < 1 {
			return 0
		}
		// If there is a gap, confirm if the next segment is valid and doesn't exceed the end of the string
		sectionLen := springMap.springList[listPos]
		// if print {
		// 	fmt.Printf("sectionLen:%d  %d/%d [%s]\n", pos, sectionLen, len(springMap.springMap), springMap.springMap[pos:pos+sectionLen])
		// }
		if pos+sectionLen <= len(springMap.springMap) {
			// If it fits in the string, check if it is a valid segment
			if _, match := ttable[sectionLen-1][springMap.springMap[pos:pos+sectionLen]]; match {
				// If the next segment is valid, return the result of the next position, increaseing sections accounted for and resetting gap
				return computePossibleSpringMapCombinations(springMap, ttable, print, inferredStr+strings.Repeat("#", sectionLen), pos+sectionLen, listPos+1, 0)
			}
		}
		// Will only get here If the next segment isn't valid, return 0
		return 0
	} else if springMap.springMap[pos] == '?' {
		var sum int

		// Here we have options
		// If gap is 0, we can't start a new segment, so this must be a gap
		if gap == 0 {
			return computePossibleSpringMapCombinations(springMap, ttable, print, inferredStr+".", pos+1, listPos, gap+1)
		}
		// If gap is > 0, we can start a new segment or add another gap, try both
		sum += computePossibleSpringMapCombinations(springMap, ttable, print, inferredStr+".", pos+1, listPos, gap+1)
		sectionLen := springMap.springList[listPos]
		// Check this won't exceed the end of the string
		if pos+sectionLen <= len(springMap.springMap) {
			// If it fits in the string, check if it is a valid segment
			if _, match := ttable[sectionLen-1][springMap.springMap[pos:pos+sectionLen]]; match {
				sum += computePossibleSpringMapCombinations(springMap, ttable, print, inferredStr+strings.Repeat("#", sectionLen), pos+sectionLen, listPos+1, 0)
			}
		}
		return sum
	} else {
		panic("Invalid character in spring map")
	}
}

// Build truth table for spring map
// 1. Build a list of all possible spring map segments, upto 15 characters long
// 2. Calculate if the spring map segment is valid and store in the map
func buildSpringMapTruthTable() [15]map[string]bool {
	table := [15]map[string]bool{}
	for i := 0; i < 15; i++ {
		buildSpringMapTrueTableForLength(&table, i+1)
	}
	return table
}

func buildSpringMapTrueTableForLength(table *[15]map[string]bool, length int) {
	// Create a map to store the truth table
	(*table)[length-1] = make(map[string]bool)

	if length == 1 {
		for _, c := range []rune{'#', '?'} {
			// Add the character to the spring map segment
			(*table)[length-1][string(c)] = true
		}
	} else {
		// Loop ovar all previous valid spring map segments
		for springMapSegment, _ := range (*table)[length-2] {
			// Loop over all possible characters to add to the spring map segment
			// Characters are # and ? are valid, add those and ignore all others
			for _, c := range []rune{'#', '?'} {
				// Add the character to the spring map segment
				(*table)[length-1][springMapSegment+string(c)] = true
			}
		}
	}
}

// Fast function comparing if 2 int arrays are the same
func compareIntArrays(a *[]int, b *[]int) bool {
	// Check if the arrays are the same length
	if len(*a) != len(*b) {
		return false
	}

	// Loop over the arrays
	for i, _ := range *a {
		// Check if the elements are the same
		if (*a)[i] != (*b)[i] {
			return false
		}
	}

	return true
}

// // Calculate possible spring map combinations
// func calcPossibleSpringMapCombinations(springMap springMapEntry) []string {
// 	// Create return variable
// 	var possibleSpringMapCombinations []string

// 	// Calculate number of broken springs
// 	var numUnkownBrokenSprings int
// 	for _, springListEntry := range springMap.springList {
// 		numUnkownBrokenSprings += springListEntry
// 	}
// 	// Remove the number of known broken springs from the list
// 	numUnkownBrokenSprings -= countCharacters(springMap.springMap, '#')
// 	fmt.Println("numUnkownBrokenSprings:", numUnkownBrokenSprings)

// 	// CHeck if there are unknowns, if not just return add the string map and return
// 	if springMap.numUnkonwns < 1 {
// 		possibleSpringMapCombinations = append(possibleSpringMapCombinations, springMap.springMap)
// 		return possibleSpringMapCombinations
// 	}

// 	// Calculate the number of possible combinations
// 	numPossibleCombinations := twoToPowerOf(springMap.numUnkonwns)
// 	// fmt.Printf("numPossibleCombinations: %d\n", numPossibleCombinations)

// 	// Loop over all possible combinations
// 	for i := 0; i < numPossibleCombinations; i++ {
// 		// Convert the integer to a binary string
// 		binaryString := intToBinaryString(i, springMap.numUnkonwns)
// 		// fmt.Printf("binaryString: %s\n", binaryString)
// 		// Check if the number of broken springs matches the number in the list, otherwise skip
// 		// fmt.Println("binary characters:", countCharacters(binaryString, '1'))
// 		if countCharacters(binaryString, '1') != numUnkownBrokenSprings {
// 			continue
// 		}
// 		mergedSpringMap := mergeSpringMap(springMap.springMap, binaryString)
// 		// fmt.Printf("mergedSpringMap: %s\n", mergedSpringMap)
// 		mergedSpringMapList := classifySpringMap(mergedSpringMap)
// 		if reflect.DeepEqual(mergedSpringMapList, springMap.springList) {
// 			possibleSpringMapCombinations = append(possibleSpringMapCombinations, mergedSpringMap)
// 		}
// 	}
// 	return possibleSpringMapCombinations
// }

// Merge the spring map with the binary numer
func mergeSpringMap(springMap string, bits uint64) uint64 {
	var result uint64
	// Current position in the binary string
	var bitsPos int = 0

	// Loop over the spring map
	strMapLen := len(springMap)
	for i := 0; i < strMapLen; i++ {
		// Check if the character is a question mark
		if springMap[i] == '?' {
			// Add the next character of the binary string to the string builder
			// Translate 0 to . and 1 to #
			if bits&1<<bitsPos == 0 {
				result |= 1 << uint64(strMapLen-i-1)
			} else {
				result &= ^(1 << uint64(strMapLen-i-1))
			}
			// Increment the binary string position
			bitsPos++
		} else if springMap[i] == '#' {
			// Add the character to the string builder
			result |= 1 << uint64(strMapLen-i-1)
		} else {
			result &= ^(1 << uint64(strMapLen-i-1))
		}
	}

	return result
}

// // Merge the spring map with the binary string
// func mergeSpringMap(springMap string, binaryString string) string {
// 	// Create a string builder
// 	var sb strings.Builder

// 	// Current position in the binary string
// 	var binaryStringPos int = 0

// 	// Loop over the spring map
// 	for _, c := range springMap {
// 		// Check if the character is a question mark
// 		if c == '?' {
// 			// Check if there are any characters in the binary string
// 			if binaryStringPos < len(binaryString) {
// 				// Add the next character of the binary string to the string builder
// 				// Translate 0 to . and 1 to #
// 				if rune(binaryString[binaryStringPos]) == '0' {
// 					sb.WriteRune('.')
// 				} else {
// 					sb.WriteRune('#')
// 				}
// 				// Increment the binary string position
// 				binaryStringPos++
// 			} else {
// 				// We've run out of binary string which shouldn't happen, panic
// 				panic("Ran out of binary string")
// 			}
// 		} else {
// 			// Add the character to the string builder
// 			sb.WriteRune(c)
// 		}
// 	}

// 	return sb.String()
// }

// Convert the integer to a binary string
func intToBinaryString(n int, length int) string {
	// Convert the integer to a binary string
	binaryString := strconv.FormatInt(int64(n), 2)

	// Check the length of the binary string
	if len(binaryString) < length {
		// Pad the binary string with zeros
		binaryString = strings.Repeat("0", length-len(binaryString)) + binaryString
	}

	return binaryString
}

// Calculate 2 to the power of n
func twoToPowerOf(n int) int {
	// Calculate 2 to the power of n
	twoToPowerOf := 1 << n
	return twoToPowerOf
}

// Function to classify a spring map entry encoded in a unit64 as a spring list
func classifySpringMapUint64(springMap uint64) []int {
	// Define the return list
	var springList []int

	// Remeber state of whether inside a spring set and the spring set count
	var insideSpringSet bool
	var springSetCount int = 0

	// Loop over the spring map
	for i := 0; i < 64; i++ {
		// Check if we are inside a spring set
		if insideSpringSet {
			// Check if we are at the end of the spring set
			if springMap&(1<<uint64(i)) == 0 {
				// We are at the end of the spring set
				insideSpringSet = false
				// If the count is > 0 then add it to the spring list and reset the count
				if springSetCount > 0 {
					springList = append(springList, springSetCount)
					springSetCount = 0
				}
			} else {
				// We are not at the end of the spring set
				springSetCount++
			}
		} else {
			// Check if we are at the start of a spring set
			if springMap&(1<<uint64(i)) != 0 {
				// We are at the start of a spring set
				insideSpringSet = true
				springSetCount = 1
			}
		}

		// Check if we are at the end of a spring set
		if !insideSpringSet && springSetCount > 0 {
			// We are at the end of a spring set
			springList = append(springList, springSetCount)
		}
	}
	if springSetCount > 0 {
		// We are at the end of a spring set
		springList = append(springList, springSetCount)
	}
	return springList
}

// Function to classify a spring map as a spring list
func classifySpringMap(springMap string) []int {
	// Define the return list
	var springList []int

	// Remeber state of whether inside a spring set and the spring set count
	var insideSpringSet bool
	var springSetCount int = 0

	// Loop over the spring map
	for _, c := range springMap {
		// Check if we are inside a spring set
		if insideSpringSet {
			// Check if we are at the end of the spring set
			if c == '.' {
				// We are at the end of the spring set
				insideSpringSet = false
				// If the count is > 0 then add it to the spring list and reset the count
				if springSetCount > 0 {
					springList = append(springList, springSetCount)
					springSetCount = 0
				}
			} else {
				// We are not at the end of the spring set
				springSetCount++
			}
		} else {
			// Check if we are at the start of a spring set
			if c == '#' {
				// We are at the start of a spring set
				insideSpringSet = true
				springSetCount = 1
			}
		}

		// Check if we are at the end of a spring set
		if !insideSpringSet && springSetCount > 0 {
			// We are at the end of a spring set
			springList = append(springList, springSetCount)
		}
	}
	if springSetCount > 0 {
		// We are at the end of a spring set
		springList = append(springList, springSetCount)
	}
	return springList
}

// Print the spring map
func printSpringMap(springMap []springMapEntry, includeArrangements bool) {
	// Loop over the spring map
	for _, springMapEntry := range springMap {
		// Print the spring map entry
		fmt.Printf("%s %v %d", springMapEntry.springMap, springMapEntry.springList, springMapEntry.numUnkonwns)
		if includeArrangements {
			// fmt.Printf(" - arrangements %d", len(springMapEntry.arrangements))
			fmt.Printf(" - arrangements %d", springMapEntry.arrangements)
		}
		fmt.Println()
	}
}

// Print the spring map

// Print spring map entry
func printSpringMapEntry(springMapEntry springMapEntry) {
	// Print the spring map entry
	// fmt.Printf("%s %v %d - arrangements %d\n", springMapEntry.springMap, springMapEntry.springList, springMapEntry.numUnkonwns, len(springMapEntry.arrangements))
	fmt.Printf("%s %v %d - arrangements %d\n", springMapEntry.springMap, springMapEntry.springList, springMapEntry.numUnkonwns, springMapEntry.arrangements)
}

// Parse the spring map
func parseSpringMap(springMapStr string) []springMapEntry {
	// Split the string into lines
	springMapLines := splitLines(springMapStr)

	// Create a slice of spring map entries
	springMap := make([]springMapEntry, len(springMapLines))

	// Loop over the lines
	for i, line := range springMapLines {
		// Split the line into the spring map and the spring list
		springMapEntry := splitSpringMapEntry(line)

		// Add the spring map entry to the slice
		springMap[i] = springMapEntry
	}

	return springMap
}

// Split the spring map entry into the spring map and the spring list
func splitSpringMapEntry(springMapEntryStr string) springMapEntry {
	// Split the string into the spring map and the spring list
	springMapEntryParts := splitString(springMapEntryStr, " ")

	// Create a spring map entry
	springMapEntry := springMapEntry{
		springMap:  springMapEntryParts[0],
		springList: splitStringToInt(springMapEntryParts[1], ","),
	}

	// Calculate the number of unknowns
	springMapEntry.numUnkonwns = countCharacters(springMapEntry.springMap, '?')

	return springMapEntry
}

// Split a string into a slice of strings
func splitString(str string, sep string) []string {
	// Split the string into parts
	parts := strings.Split(str, sep)

	return parts
}

// Split a string into a slice of ints
func splitStringToInt(str string, sep string) []int {
	// Split the string into parts
	parts := splitString(str, sep)

	// Create a slice of ints
	ints := make([]int, len(parts))

	// Loop over the parts
	for i, part := range parts {
		// Convert the part to an int
		ints[i] = strToInt(part)
	}

	return ints
}

// Split a string into a slice of strings
func splitLines(str string) []string {
	// Split the string into lines
	lines := strings.Split(str, "\n")

	return lines
}

// Convert a string to an int
func strToInt(str string) int {
	// Convert the string to an int
	num, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}

	return num
}

// Count the number of characters in a string
func countCharacters(str string, char rune) int {
	// Count the number of characters
	count := 0
	for _, c := range str {
		if c == char {
			count++
		}
	}

	return count
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
