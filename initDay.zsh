#!/bin/zsh

# Check if a parameter is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <2_digit_day_number>"
    exit 1
fi

if ! [[ $1 =~ ^[[:xdigit:]]{2}$ ]]; then
    echo "Usage: $0 <2_digit_day_number>"
    exit 2
fi

DIRECTORY="day$1"

if [ -d "$DIRECTORY" ]; then
  echo "$DIRECTORY already exists."
  exit 3
fi


# Create a directory
mkdir "$DIRECTORY"

# Change directory
cd "$DIRECTORY"

# Initialize go.mod
go mod init "github.com/4d46/AdventOfCode2023/$DIRECTORY"

# Create a default Go template
echo "package main

import (
	\"fmt\"
	\"os\"
)

func main() {
    fmt.Printf(\"Advent of Code 2023 - Day %2d\\\n\", ${1})
}

// Load file contents into a string and return it
func loadFileContents(filename string) string {
        // Read contents of file into a string
        fileBytes, err := os.ReadFile(filename) // just pass the file name
        if err != nil {
                panic(err)
        }

        return string(fileBytes) // convert content to a 'string'
}" > main.go

go build

git checkout -b "$DIRECTORY"

git add main.go go.mod

git commit -m "Day ${1} Starter"

git push --set-upstream upstream "$DIRECTORY"
