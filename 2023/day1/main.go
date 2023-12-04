package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"
)

//go:embed input
var input string

var inputStrings []string

func init() {
	input = strings.TrimRight(input, "\n")
	if len(input) == 0 {
		panic("empty input file")
	}

	inputStrings = strings.Split(input, "\n")
}

func parseFlags() (part int) {
	flag.IntVar(&part, "p", 0, "part")
	flag.Parse()
	if part > 2 || part < 0 {
		fmt.Println("part specified is out of range, defaulting to both")
	}
	return
}

func main() {
	part := parseFlags()
	if part == 1 {
		partOne()
	}
	if part == 2 {
		partTwo()
	}
	if part == 0 {
		partOne()
		partTwo()
	}
}

func partOne() {

}

func partTwo() {

}
