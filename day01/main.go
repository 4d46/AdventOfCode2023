package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Define multiline string
// var multilineString = `1abc2
// pqr3stu8vwx
// a1b2c3d4e5f
// treb7uchet`

// var multilineString = `two1nine
// eightwothree
// abcone2threexyz
// xtwone3four
// 4nineeightseven2
// zoneight234
// 7pqrstsixteen`

func main() {
	//fmt.Println("hello world")
	loadedString := loadFileContents("CalibrationFilePart1.txt")
	// calVals := getCalibrationValues(multilineString)
	calVals := getCalibrationValues(loadedString)
	// spew.Dump(loadedString)
	// spew.Dump(calVals)

	calSum := sumCalibrationValues(calVals)
	// Print out integer with formatting
	fmt.Printf("Calibration Value Sum %d\n", calSum)
}

// Function that takes an interger array and returns the sum of all integers
func sumCalibrationValues(calibrationValues []int) int {
	var sum int
	for _, val := range calibrationValues {
		sum += val
	}
	return sum
}

// Function that returns array of integers from first and last number of each line
func getCalibrationValues(calibrationDocument string) []int {
	var calVals []int
	for _, line := range strings.Split(calibrationDocument, "\n") {
		calVals = append(calVals, getFirstAndLastNumber(line))
	}
	return calVals
}

func getFirstAndLastNumber(line string) int {
	var first int = -1
	var last int
	// Fix line by replacing words with digits
	fixedLine := substituteWordsWithDigits(line)
	fmt.Printf("%s => %s\n", line, fixedLine)

	for i, char := range fixedLine {
		if char >= '0' && char <= '9' {
			if first < 0 {
				first = i
			}
			last = i
		}
	}
	if first < 0 {
		return 0
	}
	// Convert character to integer
	firstDigit, err := strconv.Atoi(fixedLine[first : first+1])
	if err != nil {
		panic(err)
	}
	lastDigit, err := strconv.Atoi(fixedLine[last : last+1])
	if err != nil {
		panic(err)
	}

	// fmt.Printf("First Digit: %d, Last Digit: %d\n", firstDigit, lastDigit)

	return firstDigit*10 + lastDigit
}

// Substitute words in string with digits
func substituteWordWithDigit(line, find, replace string) string {
	return strings.ReplaceAll(line, find, replace)
}

func substituteWordsWithDigits(line string) string {
	line = substituteWordWithDigit(line, "one", "o1e")
	line = substituteWordWithDigit(line, "two", "t2o")
	line = substituteWordWithDigit(line, "three", "t3e")
	line = substituteWordWithDigit(line, "four", "4")
	line = substituteWordWithDigit(line, "five", "5e")
	line = substituteWordWithDigit(line, "six", "6")
	line = substituteWordWithDigit(line, "seven", "7n")
	line = substituteWordWithDigit(line, "eight", "e8t")
	line = substituteWordWithDigit(line, "nine", "n9e")
	// line = substituteWordWithDigit(line, "zero", "0")
	return line
}

// Load file contents into a string and return it
func loadFileContents(filename string) string {
	// Read contents of file into a string
	fileBytes, err := os.ReadFile("CalibrationFilePart1.txt") // just pass the file name
	if err != nil {
		panic(err)
	}

	return string(fileBytes) // convert content to a 'string'
}
