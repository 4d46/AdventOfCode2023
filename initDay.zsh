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

import \"fmt\"

func main() {
    fmt.Printf(\"Advent of Code 2023 - Day %2d\\\n\", ${1})
}" > main.go

go build

git checkout -b "$DIRECTORY"