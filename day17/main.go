package main

import (
	"fmt"
	"os"
)

func main() {
    fmt.Printf("Advent of Code 2023 - Day %2d\n", 17)
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
