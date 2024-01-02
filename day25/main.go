package main

import (
	"fmt"
	"os"
	"strings"
)

type void struct{}

// const EMPTY void = void{}

const componentLinks string = `jqt: rhn xhk nvd
rsh: frs pzl lsr
xhk: hfx
cmg: qnr nvd lhk bvb
rhn: xhk bvb hfx
bvb: xhk hfx
pzl: lsr hfx nvd
qnr: nvd
ntq: jqt hfx bvb xhk
nvd: lhk
lsr: lhk
rzs: qnr cmg lsr rsh
frs: qnr lhk lsr`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 25)

	// Load the input data
	// input := componentLinks
	input := loadFileContents("component_links.txt")

	// Process the input data
	links := parseLinks(input)

	// Print the results
	for k, v := range links {
		if len(v) > 7 {
			fmt.Printf("%s: %d\n", k, len(v))
		}
	}

	var largestAnswer int = 0
	var keyList []string
	for k, _ := range links {
		keyList = append(keyList, k)
	}
	var set [2][]string
	checkSet(&links, &keyList, set, 0, &largestAnswer)
	fmt.Printf("Largest answer: %d\n", largestAnswer)
}

// func everyCombination(allLinks *map[string]map[string]void,componentList *[]string, set [2][]string) {
//         // fmt.Printf("%s: %d\n", currentComponent, currentLength)
//         if currentLength > 7 {
//                 fmt.Printf("%s: %d\n", currentComponent, currentLength)
//         }
//         for k, v := range (*allLinks)[currentComponent] {
//                 if _, ok := (*currentLinks)[k]; !ok {
//                         (*currentLinks)[k] = v
//                         everyCombination(allLinks, currentLinks, k, currentLength+1)
//                         delete(*currentLinks, k)
//                 }
//         }
// }

func checkSet(allLinks *map[string]map[string]void, keyList *[]string, set [2][]string, depth int, largestAnswer *int) {
	// If we haven't reached the end and assigned all the components, continue on and assign next component
	if depth < len(*keyList) {
		// Try assigning to both sets, but alternate which one to try first depending on the depth
		if depth%2 == 0 {
			// Try assigning to set 1 first
			set[0] = append(set[0], (*keyList)[depth])
			checkSet(allLinks, keyList, set, depth+1, largestAnswer)
			set[0] = set[0][:len(set[0])-1]
			set[1] = append(set[1], (*keyList)[depth])
			checkSet(allLinks, keyList, set, depth+1, largestAnswer)
			set[1] = set[1][:len(set[1])-1]
		} else {
			// Try assigning to set 2 first
			set[1] = append(set[1], (*keyList)[depth])
			checkSet(allLinks, keyList, set, depth+1, largestAnswer)
			set[1] = set[1][:len(set[1])-1]
			set[0] = append(set[0], (*keyList)[depth])
			checkSet(allLinks, keyList, set, depth+1, largestAnswer)
			set[0] = set[0][:len(set[0])-1]
		}
	} else {
		// We are at the end of the list, so check the set
		// Check if there is already a better answer than this, in which case we can stop
		if len(set[0])*len(set[1]) < *largestAnswer {
			// We can stop checking this option
			return
		}
		// Check if the sets only have 3 links between them
		// Loop over all elements in the first set and count if lkinked to element in second set
		linkCount := 0
		for _, v1 := range set[0] {
			for _, v2 := range set[1] {
				if _, ok := (*allLinks)[v1][v2]; ok {
					linkCount++
					if linkCount > 3 {
						// Too many links, we can stop checking this option
						return
					}
				}
			}
		}
		// Check if the sets have 3 links between them
		if linkCount == 3 {
			// We have a new best answer
			*largestAnswer = len(set[0]) * len(set[1])
			fmt.Printf("New best answer: %d\n", *largestAnswer)
		}
	}
}

// Parse the input data into a map of components
func parseLinks(input string) map[string]map[string]void {
	// Create a map to hold the components
	links := make(map[string]map[string]void)

	// Split the input into lines
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		// Split the line into components
		components := strings.Split(line, ": ")
		subcomponents := strings.Split(components[1], " ")
		for _, subcomponent := range subcomponents {
			// Add the subcomponent to the map
			if _, ok := links[components[0]]; !ok {
				links[components[0]] = make(map[string]void)
			}
			links[components[0]][subcomponent] = void{}
			if _, ok := links[subcomponent]; !ok {
				links[subcomponent] = make(map[string]void)
			}
			links[subcomponent][components[0]] = void{}
		}
	}

	return links
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
