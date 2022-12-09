package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"strings"
)

//go:embed input
var input []byte

func main() {
	var (
		item            rune
		lineIndex       int
		s1, s2, s3      string
		partOnePriority int
		partTwoPriority int
	)

	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		s := scanner.Text()

		// part 1
		item = findRepeatedItem(s)
		partOnePriority += priority(item)

		// part 2
		switch lineIndex % 3 {
		case 0:
			s1 = s
		case 1:
			s2 = s
		case 2:
			s3 = s
			item = findBadgeType(s1, s2, s3)
			partTwoPriority += priority(item)
		}

		lineIndex += 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Part 1 Priority: %d\n", partOnePriority)
	fmt.Printf("Part 2 Priority: %d\n", partTwoPriority)
}

func findRepeatedItem(s string) rune {
	halfway := len(s) / 2
	firstCompartment := s[:halfway]
	secondCompartment := s[halfway:]

	for _, b := range firstCompartment {
		if !strings.ContainsRune(secondCompartment, b) {
			continue
		}
		return b
	}
	panic("no repeated item")
}

func findBadgeType(s1, s2, s3 string) rune {
	for _, b := range s1 {
		if !strings.ContainsRune(s2, b) {
			continue
		}
		if !strings.ContainsRune(s3, b) {
			continue
		}
		return b
	}
	panic("no badge type")
}

func priority(b rune) int {
	if b >= 'a' {
		return int(b - 'a' + 1)
	}
	return int(b - 'A' + 27)
}
