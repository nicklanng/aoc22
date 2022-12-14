package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
	"unicode"
)

const stackCount = 9

//go:embed input
var input []byte

func main() {
	scanner := bufio.NewScanner(bytes.NewReader(input))

	stacks, err := parseStartingImage(scanner)
	if err != nil {
		panic(err)
	}

	for scanner.Scan() {
		s := scanner.Text()

		var numberOfMoves, src, dst int
		if _, err := fmt.Sscanf(s, "move %d from %d to %d", &numberOfMoves, &src, &dst); err != nil {
			panic(err)
		}

		// crateMover9000(&stacks, numberOfMoves, src, dst)
		crateMover9001(&stacks, numberOfMoves, src, dst)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	var topCrates string
	for i := 0; i < stackCount; i++ {
		crate, _ := stacks[i].Pop()
		topCrates += string(crate)
	}
	fmt.Printf("Top crates: %s\n", topCrates)
}

func parseStartingImage(scanner *bufio.Scanner) ([9]lib.Stack[rune], error) {
	var stacks [stackCount]lib.Stack[rune]

	for scanner.Scan() {
		s := scanner.Text()

		if len(s) == 0 {
			break
		}

		for i := 0; i < stackCount; i++ {
			index := i*4 + 1
			if index >= len(s) {
				break
			}
			if unicode.IsLetter(rune(s[index])) {
				stacks[i].Data = append([]rune{rune(s[index])}, stacks[i].Data...)
			}
		}
	}

	return stacks, scanner.Err()
}

func crateMover9000(stacks *[9]lib.Stack[rune], numberOfMoves, src, dst int) {
	for i := 0; i < numberOfMoves; i++ {
		crate, ok := stacks[src-1].Pop()
		if !ok {
			panic("not enough crates")
		}
		stacks[dst-1].Push(crate)
	}
}

func crateMover9001(stacks *[9]lib.Stack[rune], numberOfMoves, src, dst int) {
	var toBeMoved []rune

	// pop the stacks to a temp holding slice
	for i := 0; i < numberOfMoves; i++ {
		crate, ok := stacks[src-1].Pop()
		if !ok {
			panic("not enough crates")
		}
		toBeMoved = append(toBeMoved, crate)
	}

	// push in reverse
	for i := len(toBeMoved) - 1; i >= 0; i-- {
		stacks[dst-1].Push(toBeMoved[i])
	}
}
