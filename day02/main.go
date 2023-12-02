package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("hello world")
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
