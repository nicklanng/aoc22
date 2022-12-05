package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
	"log"
	"strings"
)

//go:embed input
var input []byte

func main() {
	var (
		partOneCount int
		partTwoCount int
	)

	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		s := scanner.Text()

		assignment1Str, assignment2Str, ok := strings.Cut(s, ",")
		if !ok {
			panic("no comma in assignment pairs")
		}

		assignment1, ok := parseAssignment(assignment1Str)
		if !ok {
			panic("failed to parse assignment 1 string")
		}

		assignment2, ok := parseAssignment(assignment2Str)
		if !ok {
			panic("failed to parse assignment 2 string")
		}

		if assignment1.contains(assignment2) || assignment2.contains(assignment1) {
			partOneCount++
		}

		if assignment1.intersects(assignment2) {
			partTwoCount++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Part 1 Count: %d\n", partOneCount)
	fmt.Printf("Part 2 Count: %d\n", partTwoCount)
}

func parseAssignment(s string) (assignment, bool) {
	minStr, maxStr, ok := strings.Cut(s, "-")
	if !ok {
		return assignment{}, false
	}

	min := lib.MustParseInt(minStr)
	max := lib.MustParseInt(maxStr)

	return assignment{min: min, max: max}, true
}

type assignment struct {
	min, max int
}

func (a assignment) contains(a2 assignment) bool {
	return a2.min >= a.min && a2.max <= a.max
}

func (a assignment) intersects(a2 assignment) bool {
	return !(a2.max < a.min || a2.min > a.max)
}
